package backend

import (
	"context"
	"encoding/json"
	"math"
	"sync"
	"time"

	"github.com/canhlinh/cointracker/config"
	"github.com/canhlinh/cointracker/types"
	socket "github.com/canhlinh/tradingview-scraper/v2"
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

type Price struct {
	LastTimePrice float64
	CurrentPrice  float64
	Started       bool
}

type App struct {
	DB           *DB
	wsClosedChan chan bool
	workerPool   *work.WorkerPool
	enqueuer     *work.Enqueuer
	cronJob      *cron.Cron
	symbols      map[string]bool
	mutex        *sync.Mutex
	tvSocket     socket.SocketInterface
	wait         *sync.WaitGroup
}

func NewApp() *App {
	app := &App{
		DB:      NewDB(),
		mutex:   &sync.Mutex{},
		wait:    &sync.WaitGroup{},
		symbols: map[string]bool{},
	}
	return app.initBackgroundJobs()
}

type PayloadCheckAlert struct {
	AlertConfigID int64     `json:"alert_config_id"`
	Timeframe     int       `json:"timeframe"`
	ClosedAt      time.Time `json:"closed_at"`
}

func (payload PayloadCheckAlert) String() string {
	b, err := json.Marshal(payload)
	if err != nil {
		log.Panic(err)
	}
	return string(b)
}

func (app *App) cronJobScanPriceChanges() {
	now := time.Now().Truncate(time.Minute)
	log.Infof("start scanning price changes at:%s", now.Format(time.RFC3339))

	alertConfigs, err := app.DB.GetAllAlertConfigs()
	if err != nil {
		log.Panic(err)
	}

	h := now.Hour()
	m := now.Minute()
	for _, alertConfig := range alertConfigs {
		timeframes, err := app.DB.GetTimeframes(alertConfig.Timeframes, 0)
		if err != nil {
			log.Error(err)
			return
		}
		disableTimeframes := alertConfig.MapDisabledTimeframes()

		for _, timeframe := range timeframes {
			if _, ok := disableTimeframes[timeframe.Timeframe]; ok {
				continue
			}

			delta := ((h*60 + m) % timeframe.Timeframe)
			if delta == 0 {
				app.enqueuer.Enqueue(config.RedisJobCheckPriceChange, map[string]interface{}{
					"payload": (&PayloadCheckAlert{
						AlertConfigID: alertConfig.ID,
						Timeframe:     timeframe.Timeframe,
						ClosedAt:      now,
					}).String(),
				})
			}
		}

	}
}

func (app *App) jobCheckPriceChange(job *work.Job) error {
	var jobPayload PayloadCheckAlert
	if err := json.Unmarshal([]byte(job.ArgString("payload")), &jobPayload); err != nil {
		log.Error(err)
		return nil
	}
	alertConfig, err := app.DB.GetAlertConfig(jobPayload.AlertConfigID, false)
	if err != nil {
		log.Error(err)
		return nil
	}
	timeframe, err := app.DB.GetTimeframe(int64(jobPayload.Timeframe))
	if err != nil {
		log.Error(err)
		return nil
	}

	opennedAt := jobPayload.ClosedAt.Add(-time.Minute * time.Duration(timeframe.Timeframe))
	closedAt := jobPayload.ClosedAt

	openPrice, closePrice := app.DB.GetFirstAndLastPrice(alertConfig.Symbol, opennedAt, closedAt)
	if openPrice != 0 || closePrice != 0 {
		percentChanged := calculatePercentChanged(openPrice, closePrice)
		log.Infof("symbol:%s price_changed:%0.2f timeframe:%v open:%v close:%v openned_at:%v closed_at:%v",
			alertConfig.Symbol, percentChanged, timeframe, openPrice, closePrice, opennedAt, closedAt)

		kline := &types.Kline{
			Symbol:         alertConfig.Symbol,
			Timeframe:      jobPayload.Timeframe,
			Open:           openPrice,
			Close:          closePrice,
			PercentChanged: percentChanged,
			OpennedAt:      opennedAt,
			ClosedAt:       closedAt,
		}

		isPriceUp := closePrice > openPrice && (alertConfig.Direction == types.DirectionBoth || alertConfig.Direction == types.DirectionUp)
		isPriceDown := closePrice < openPrice && (alertConfig.Direction == types.DirectionBoth || alertConfig.Direction == types.DirectionDown)
		if (math.Abs(percentChanged) >= timeframe.PercentAlert) && (isPriceUp || isPriceDown) {
			sendTeleMsg(buildNotificationMessage(alertConfig.Symbol, percentChanged, timeframe.Timeframe, openPrice, closePrice))

			alert := kline.ToAlert()
			alert.AlertConfigID = jobPayload.AlertConfigID
			if err := app.DB.WriteAlert(alert); err != nil {
				log.Error(err)
			}

			tx := app.DB.Begin()
			defer tx.Commit()
			alertConfig, err := tx.GetAlertConfig(jobPayload.AlertConfigID, true)
			if err != nil {
				log.Error(err)
				return nil
			}
			if err := tx.UpdateTriggered(alertConfig, jobPayload.Timeframe, time.Now()); err != nil {
				log.Error(err)
			}
			tx.Commit()
		}

		if config.Config().StoreKline {
			if err := app.DB.WriteKline(kline); err != nil {
				log.Error(err)
				return err
			}
		}
	}
	return nil
}

func (app *App) cronDeleteOldData() {
	log.Info("cronDeleteOldData")
	now := time.Now()
	start := time.Now().AddDate(-1, 0, 0)
	dataQuoteDataStopAt := now.Add(-time.Hour * time.Duration(config.Config().DeleteQuoteDataOlderThan))
	dataKlineDataStopAt := now.Add(-time.Hour * time.Duration(config.Config().DeleteKlineDataOlderThan))
	dataAlertDataStopAt := now.Add(-time.Hour * 48)

	if err := app.DB.influxdbClient.DeleteAPI().DeleteWithName(context.Background(), config.Config().InfluxDBOrg, config.Config().InfluxDBBucket, start, dataQuoteDataStopAt, `_measurement="tickers"`); err != nil {
		log.Error(err)
	} else {
		log.Info("Deleted old tickers")
	}
	if err := app.DB.sql.Where("created_at <= ?", dataKlineDataStopAt).Delete(&types.Kline{}).Error; err != nil {
		log.Error(err)
	} else {
		log.Info("Deleted old klines")
	}
	if err := app.DB.sql.Where("created_at <= ?", dataAlertDataStopAt).Delete(&types.Alert{}).Error; err != nil {
		log.Error(err)
	} else {
		log.Info("Deleted old alerts")
	}
}

func (app *App) cronCleanDisabledConfigs() {
	log.Info("cronCleanDisabledConfigs")
	if err := app.DB.CleanDisabledTimeframes(); err != nil {
		log.Error(err)
	}
}

func (app *App) initBackgroundJobs() *App {
	var redisPool = &redis.Pool{
		MaxActive: 10,
		MaxIdle:   10,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(config.Config().RedisDsn)
		},
	}
	workerPool := work.NewWorkerPool(*app, 10, config.RedisJobNamespace, redisPool)
	workerPool.Job(config.RedisJobCheckPriceChange, app.jobCheckPriceChange)
	workerPool.Start()

	cronJob := cron.New()
	cronJob.AddFunc("* * * * *", app.cronJobScanPriceChanges)
	cronJob.AddFunc("0 * * * *", app.cronDeleteOldData)
	cronJob.AddFunc("0 0 * * *", app.cronCleanDisabledConfigs)
	cronJob.Start()

	app.workerPool = workerPool
	app.enqueuer = work.NewEnqueuer(config.RedisJobNamespace, redisPool)
	app.cronJob = cronJob
	return app
}

func (app *App) Run() {
	app.wsClosedChan = make(chan bool)

CONNECT_WS:
	app.wait.Add(1)
	log.Info("connecting to TradingView")

	tvSocket, err := socket.Connect(
		func(symbol string, data *socket.QuoteData) {
			if data.Price != nil {
				if err := app.DB.WriteQuoteData(&types.Ticker{
					Symbol: symbol,
					Price:  *data.Price,
				}); err != nil {
					log.Error(err)
				}
			}
		},
		func(err error, context string) {
			log.Warnf("error -> %#v context-> %v", err.Error(), context)
			sendTeleMsg(err.Error())
			app.wsClosedChan <- true
		},
	)
	if err != nil {
		log.Fatal("failed to connect to TradingView, error: " + err.Error())
	}

	app.wait.Done()
	app.tvSocket = tvSocket
	log.Info("connected to TradingView")

	if err := app.updateListenSymbols(); err != nil {
		log.Error(err)
	}

	<-app.wsClosedChan
	<-time.After(time.Second * 10)
	app.tvSocket = nil
	sendTeleMsg("disconnected from TradingView")
	log.Warn("disconnected from TradingView")
	goto CONNECT_WS
}

func (app *App) updateListenSymbols() error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	app.wait.Wait()

	if app.tvSocket == nil {
		log.Error("tvSocket empty")
	}

	for symbol, _ := range app.symbols {
		if err := app.tvSocket.RemoveSymbol(symbol); err != nil {
			return err
		}
	}
	app.symbols = make(map[string]bool)

	symbols, err := app.DB.GetConfigSymbols()
	if err != nil {
		return err
	}

	for _, symbol := range symbols {
		if _, ok := app.symbols[symbol]; !ok {
			if err := app.tvSocket.AddSymbol(symbol); err != nil {
				return err
			}
			app.symbols[symbol] = true
		}
	}
	return nil
}

func (app *App) AddSymbol(symbol string) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	app.wait.Wait()

	if _, ok := app.symbols[symbol]; ok {
		return nil
	}

	if err := app.tvSocket.AddSymbol(symbol); err != nil {
		return err
	}
	app.symbols[symbol] = true
	return nil
}

func (app *App) RemoveSymbol(symbol string) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	app.wait.Wait()

	if _, ok := app.symbols[symbol]; !ok {
		return nil
	}

	if err := app.tvSocket.RemoveSymbol(symbol); err != nil {
		return err
	}
	delete(app.symbols, symbol)
	return nil
}

func (app *App) CreateAlertConfig(alertConfig *types.AlertConfig) error {
	if len(alertConfig.Timeframes) == 0 {
		timeframes, err := app.DB.GetDefaultTimeframes()
		if err != nil {
			return err
		}
		alertConfig.Timeframes = timeframes
	}

	if err := app.DB.SaveAlertConfig(alertConfig); err != nil {
		return err
	}

	if err := app.AddSymbol(alertConfig.Symbol); err != nil {
		return err
	}
	return nil
}

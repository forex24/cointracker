package backend

import (
	"context"
	"fmt"

	"time"

	"github.com/canhlinh/cointracker/config"
	"github.com/canhlinh/cointracker/types"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	log "github.com/sirupsen/logrus"
	"gopkg.in/guregu/null.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type DB struct {
	influxdbClient influxdb2.Client
	Org            string
	Bucket         string
	sql            *gorm.DB
}

func NewDB() *DB {
	gormConfig := &gorm.Config{}
	if config.Config().Debug {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}
	mysqlDSN := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=UTC",
		config.Config().MySQLUser,
		config.Config().MySQLPassword,
		config.Config().MySQLHost,
		config.Config().MySQLDatabase,
	)
	log.Info(mysqlDSN)

	schema.RegisterSerializer("json", JSONSerializer{})
	sql, err := gorm.Open(mysql.Open(mysqlDSN), gormConfig)
	if err != nil {
		panic("failed to connect database")
	}

	influxdbClient := influxdb2.NewClient(config.Config().InfluxDBUrl, config.Config().InfluxDBToken)
	db := &DB{
		influxdbClient: influxdbClient,
		Org:            config.Config().InfluxDBOrg,
		Bucket:         config.Config().InfluxDBBucket,
		sql:            sql,
	}
	db.migrate()
	// db.loadSettings()
	return db
}

func (db *DB) Begin() *DB {
	return &DB{
		sql:            db.sql.Begin(),
		influxdbClient: db.influxdbClient,
		Org:            db.Org,
		Bucket:         db.Bucket,
	}
}

func (db *DB) Commit() error {
	return db.sql.Commit().Error
}

func (db *DB) migrate() {
	if err := db.sql.AutoMigrate(&types.AlertConfig{}, &types.Kline{}, &types.Timeframe{}, &types.Alert{}); err != nil {
		log.Error(err)
	}

	for _, timeframe := range types.Timeframes {
		if err := db.sql.FirstOrCreate(timeframe).Error; err != nil {
			log.Error(err)
		}
	}
}

func (db *DB) DeleteConfigs() error {
	return db.sql.Migrator().DropTable("alert_configs")
}

func (db *DB) SaveAlertConfig(alertConfig *types.AlertConfig) error {
	return db.sql.Where("symbol = ?", alertConfig.Symbol).Attrs(alertConfig).FirstOrCreate(alertConfig).Error
}

func (db *DB) UpdateAlertConfig(alertConfig *types.AlertConfig) error {
	if err := db.sql.Where("id = ?", alertConfig.ID).Updates(alertConfig).Error; err != nil {
		return nil
	}
	return nil
}

func (db *DB) GetAlertConfigs(search *string, offset, limit int) ([]*types.AlertConfig, int64, error) {
	var alertConfigs []*types.AlertConfig
	var total int64

	q := db.sql
	if search != nil {
		q = q.Where("symbol LIKE ?", formatLikeQuery(*search))
	}
	q.Model(alertConfigs).Count(&total)

	if err := q.Offset(offset).Limit(limit).Find(&alertConfigs).Error; err != nil {
		return nil, total, err
	}
	return alertConfigs, total, nil
}

func (db *DB) GetConfigSymbols() ([]string, error) {
	var symbols []string

	if err := db.sql.Select("symbol").Model(&types.AlertConfig{}).Pluck("symbol", &symbols).Error; err != nil {
		return nil, err
	}

	return symbols, nil
}

func (db *DB) GetAlertConfig(id int64, forUpdate bool) (*types.AlertConfig, error) {
	var alertConfig types.AlertConfig
	query := db.sql.Where("id = ?", id)
	if forUpdate {
		query = query.Clauses(clause.Locking{Strength: "UPDATE"})
	}

	if err := query.First(&alertConfig).Error; err != nil {
		return nil, err
	}

	return &alertConfig, nil
}

func (db *DB) DeleteAlertConfig(id int64) error {
	if err := db.sql.Where("id = ?", id).Delete(&types.AlertConfig{}).Error; err != nil {
		return err
	}

	return nil
}

func (db *DB) WriteQuoteData(ticker *types.Ticker) error {
	writeAPI := db.influxdbClient.WriteAPI(db.Org, db.Bucket)
	tags := map[string]string{
		"symbol": ticker.Symbol,
	}
	fields := map[string]interface{}{
		"price": ticker.Price,
	}
	writeAPI.WritePoint(write.NewPoint("tickers", tags, fields, time.Now()))
	return nil
}

func (db *DB) WriteKline(kline *types.Kline) error {
	return db.sql.Create(kline).Error
}

func (db *DB) WriteAlert(alert *types.Alert) error {
	return db.sql.Create(alert).Error
}

func (db *DB) GetFirstAndLastPrice(symbol string, startAt time.Time, stopAt time.Time) (float64, float64) {
	query := fmt.Sprintf(`
	data = from(bucket: "%s")
		|> range(start: %s, stop: %s)
		|> filter(fn: (r) => r["_measurement"] == "tickers")
		|> filter(fn: (r) => r._field == "price")
		|> filter(fn: (r) => r["symbol"] == "%s")
	first_record = data |> first() |> set(key: "_field", value: "delta")
	last_record = data |> last() |> set(key: "_field", value: "delta")
	union(tables: [first_record, last_record])`, db.Bucket, startAt.Format(time.RFC3339), stopAt.Format(time.RFC3339), symbol)
	result, err := db.influxdbClient.QueryAPI(db.Org).Query(context.Background(), query)
	if err != nil {
		panic(err)
	}
	var data []float64
	for result.Next() {
		data = append(data, result.Record().Value().(float64))
	}
	if len(data) != 2 {
		log.Errorf("invalid open and close price, symbol:%v", symbol)
		return 0, 0
	}
	return data[0], data[1]
}

func (db *DB) GetKlines(search *string, offset int, limit int) ([]*types.Kline, int64, error) {
	var klines []*types.Kline
	var total int64

	query := db.sql
	if search != nil {
		query = query.Where("symbol LIKE ?", formatLikeQuery(*search))
	}

	query.Model(klines).Count(&total)

	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&klines).Error; err != nil {
		return nil, 0, err
	}
	return klines, total, nil
}

func (db *DB) GetTimeframes(ids []int, alertConfigID int64) ([]*types.Timeframe, error) {
	query := db.sql.Where("enable = 1")
	if len(ids) > 0 {
		query = query.Where("timeframe IN (?)", ids)
	}
	if alertConfigID > 0 {
		query = query.Joins("JOIN alert_config_timeframes ON alert_config_timeframes.alert_config_id = ?", alertConfigID)
	}

	var timeframes []*types.Timeframe
	if err := query.Find(&timeframes).Error; err != nil {
		return nil, err
	}
	return timeframes, nil
}

func (db *DB) GetAlerts(search *string, offset int, limit int) ([]*types.Alert, int64, error) {
	var alerts []*types.Alert
	var total int64

	query := db.sql
	if search != nil {
		query = query.Where("symbol LIKE ?", formatLikeQuery(*search))
	}

	query.Model(alerts).Count(&total)

	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&alerts).Error; err != nil {
		return nil, 0, err
	}
	return alerts, total, nil
}

func (db *DB) UpdateTriggered(alertConfig *types.AlertConfig, timeframe int, t time.Time) error {
	alertConfig.LastTriggeredAt = null.TimeFrom(t)
	if alertConfig.AutoDisableAfterTrigger.Bool {
		if _, ok := alertConfig.MapDisabledTimeframes()[timeframe]; !ok {
			alertConfig.DisabledTimeframes = append(alertConfig.DisabledTimeframes, timeframe)
		}
	}
	return db.sql.Model(alertConfig).Where("id = ?", alertConfig.ID).Updates(alertConfig).Error
}

func (db *DB) GetTimeframe(id int64) (*types.Timeframe, error) {
	var timeframe types.Timeframe
	if err := db.sql.Where("timeframe = ?", id).First(&timeframe).Error; err != nil {
		return nil, err
	}
	return &timeframe, nil
}

func (db *DB) UpdateTimeframe(timeframe *types.Timeframe) error {
	if err := db.sql.Where("timeframe = ?", timeframe.ID).Updates(timeframe).Error; err != nil {
		return nil
	}
	return nil
}

func (db *DB) GetDefaultTimeframes() ([]int, error) {
	var timeframes []int
	if err := db.sql.Model(&types.Timeframe{}).Where("`default` = 1 AND enable = 1").Pluck("timeframe", &timeframes).Error; err != nil {
		return nil, err
	}
	return timeframes, nil
}

func (db *DB) CleanDisabledTimeframes() error {
	updates := map[string]interface{}{
		"disabled_timeframes": gorm.Expr("NULL"),
	}
	return db.sql.Session(&gorm.Session{AllowGlobalUpdate: true}).Model(&types.AlertConfig{}).Updates(updates).Error
}

func (db *DB) GetAllAlertConfigs() ([]*types.AlertConfig, error) {
	var tickers []*types.AlertConfig
	if err := db.sql.Find(&tickers).Error; err != nil {
		return nil, err
	}
	return tickers, nil
}

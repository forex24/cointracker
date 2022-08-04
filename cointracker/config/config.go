package config

import (
	"fmt"

	"github.com/getsentry/sentry-go"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

const DefaultTimeFrame = 30 // 30 minutes
const OneDay = 24 * 60
const SqlLitePath = "/data/config.sql"

const RedisJobNamespace = "cointracker"
const RedisJobCheckPriceChange = "job_check_price_change"
const RedisJobScanPriceChanges = "job_scan_price_changes"

type config struct {
	Symbols                  []string `required:"true" envconfig:"SYMBOLS"`
	BotChatID                int      `required:"true" envconfig:"BOT_CHAT_ID"`
	BotToken                 string   `required:"true" envconfig:"BOT_TOKEN"`
	PercentPriceChangedAlert float64  `required:"true" envconfig:"PERCENT_PRICE_CHANGED_ALERT"`
	InfluxDBToken            string   `required:"true" envconfig:"INFLUXDB_TOKEN"`
	InfluxDBBucket           string   `required:"true" envconfig:"INFLUXDB_BUCKET"`
	InfluxDBOrg              string   `required:"true" envconfig:"INFLUXDB_ORG"`
	InfluxDBUrl              string   `required:"true" envconfig:"INFLUXDB_URL"`
	TimeFrame                int      `envconfig:"TIMEFRAME"` // minutes
	RedisDsn                 string   `required:"true" envconfig:"REDIS_DSN"`
	SentryDsn                string   `envconfig:"SENTRY_DSN"`
	Debug                    bool     `envconfig:"DEBUG"`
	MySQLDatabase            string   `required:"true" envconfig:"MYSQL_DATABASE"`
	MySQLHost                string   `required:"true" envconfig:"MYSQL_HOST"`
	MySQLUser                string   `required:"true" envconfig:"MYSQL_USER"`
	MySQLPassword            string   `required:"true" envconfig:"MYSQL_PASSWORD"`
	DeleteQuoteDataOlderThan int      `envconfig:"DELETE_QUOTE_DATA_OLDER_THAN" default:"48"`
	StoreKline               bool     `envconfig:"STORE_KLINE" default:"false"`
	DeleteKlineDataOlderThan int      `envconfig:"DELETE_KLINE_DATA_OLDER_THAN" default:"48"`
	DisplayTimezone          string   `envconfig:"DISPLAY_TIMEZONE" default:"Asia/Ho_Chi_Minh"`
}

var instance = &config{}

func init() {
	if err := envconfig.Process("", instance); err != nil {
		panic(err)
	}
	if instance.TimeFrame == 0 {
		instance.TimeFrame = DefaultTimeFrame
	}

	if instance.TimeFrame >= OneDay {
		panic("TimeFrame not supported")
	}
	if len(instance.SentryDsn) > 0 {
		instance.initSentry()
	}
}

func Config() config {
	return *instance
}

func (c *config) initSentry() {
	if err := sentry.Init(sentry.ClientOptions{Dsn: c.SentryDsn}); err != nil {
		logrus.WithError(err).Fatal("failed to setup sentry logging")
	}
	logrus.AddHook(NewSentryHook([]logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel, logrus.WarnLevel}))
}

func TelegramBotAPI() string {
	return fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", instance.BotToken)
}

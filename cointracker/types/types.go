package types

import (
	"encoding/json"
	"time"

	"gopkg.in/guregu/null.v3"
)

type Ticker struct {
	Symbol string
	Price  float64
}

const (
	DirectionUp   = "up"
	DirectionDown = "down"
	DirectionBoth = "both"
)

type ArrayInt []int

type AlertConfig struct {
	ID                      int64     `json:"id" gorm:"PRIMARY_KEY"`
	CreatedAt               time.Time `json:"created_at" gorm:"type:DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP"`
	UpdatedAt               time.Time `json:"updated_at" gorm:"type:DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP"`
	Symbol                  string    `json:"symbol" gorm:"unique;type:VARCHAR(64)"`
	Direction               string    `json:"direction" gorm:"type:VARCHAR(4);default:up"`
	LastTriggeredAt         null.Time `json:"last_triggered_at"`
	Timeframes              ArrayInt  `json:"timeframes" gorm:"type:JSON;serializer:json"`
	AutoDisableAfterTrigger null.Bool `json:"auto_disable_after_trigger" gorm:"default:0"`
	DisabledTimeframes      ArrayInt  `json:"disabled_timeframes" gorm:"type:JSON;serializer:json"`
}

func (alertConfig AlertConfig) MapDisabledTimeframes() map[int]bool {
	m := map[int]bool{}
	for _, t := range alertConfig.DisabledTimeframes {
		m[t] = true
	}
	return m
}

func (alertConfig *AlertConfig) String() string {
	b, _ := json.Marshal(alertConfig)
	return string(b)
}

type Kline struct {
	ID             int64     `json:"id" gorm:"PRIMARY_KEY" `
	CreatedAt      time.Time `json:"created_at" gorm:"type:DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP"`
	Symbol         string    `json:"symbol" gorm:"type:VARCHAR(64)"`
	Timeframe      int       `json:"timeframe"`
	PercentChanged float64   `json:"percent_changed"`
	Open           float64   `json:"open"`
	Close          float64   `json:"close"`
	OpennedAt      time.Time `json:"openned_at"`
	ClosedAt       time.Time `json:"closed_at"`
}

type Alert struct {
	ID             int64     `json:"id" gorm:"PRIMARY_KEY" `
	CreatedAt      time.Time `json:"created_at" gorm:"type:DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP"`
	Symbol         string    `json:"symbol" gorm:"type:VARCHAR(64)"`
	Timeframe      int       `json:"timeframe"`
	PercentChanged float64   `json:"percent_changed"`
	Open           float64   `json:"open"`
	Close          float64   `json:"close"`
	OpennedAt      time.Time `json:"openned_at"`
	ClosedAt       time.Time `json:"closed_at"`
	AlertConfigID  int64     `json:"alert_config_id"`
}

func (kline Kline) ToAlert() *Alert {
	return &Alert{
		Symbol:         kline.Symbol,
		Timeframe:      kline.Timeframe,
		PercentChanged: kline.PercentChanged,
		Open:           kline.Open,
		Close:          kline.Close,
		OpennedAt:      kline.OpennedAt,
		ClosedAt:       kline.ClosedAt,
	}
}

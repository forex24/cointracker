package types

import (
	"gopkg.in/guregu/null.v3"
	"gorm.io/gorm"
)

type Timeframe struct {
	ID           int       `json:"id" gorm:"-"`
	Timeframe    int       `json:"timeframe" gorm:"PRIMARY_KEY"`
	Format       string    `json:"format"`
	PercentAlert float64   `json:"percent_alert" gorm:"type:DECIMAL(5,2);default:10"`
	Default      null.Bool `json:"default" gorm:"default:0"`
	Enable       bool      `json:"enable" gorm:"default:0"`
}

func (t *Timeframe) AfterFind(tx *gorm.DB) error {
	t.ID = t.Timeframe
	return nil
}

var Timeframes = []*Timeframe{
	{
		Timeframe: 1,
		Format:    "1m",
		Default:   null.BoolFrom(true),
		Enable:    true,
	},
	{
		Timeframe: 3,
		Format:    "3m",
		Default:   null.BoolFrom(true),
		Enable:    true,
	},
	{
		Timeframe: 5,
		Format:    "5m",
		Default:   null.BoolFrom(true),
		Enable:    true,
	},
	{
		Timeframe: 10,
		Format:    "10m",
	},
	{
		Timeframe: 15,
		Format:    "15m",
		Default:   null.BoolFrom(true),
		Enable:    true,
	},
	{
		Timeframe: 30,
		Format:    "30m",
	},
	{
		Timeframe: 40,
		Format:    "40m",
	},
	{
		Timeframe: 60,
		Format:    "1h",
		Enable:    true,
		Default:   null.BoolFrom(true),
	},
	{
		Timeframe: 120,
		Format:    "2h",
	},
	{
		Timeframe: 180,
		Format:    "3h",
	},
	{
		Timeframe: 240,
		Format:    "4h",
		Enable:    true,
	},
	{
		Timeframe: 360,
		Format:    "6h",
	},
	{
		Timeframe: 720,
		Format:    "12h",
	},
	{
		Timeframe: 1440,
		Format:    "24h",
		Enable:    true,
	},
}

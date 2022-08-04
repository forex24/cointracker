package backend

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/canhlinh/cointracker/config"
	"gorm.io/gorm/schema"
)

func timeLocation() *time.Location {
	timezoneVN, _ := time.LoadLocation(config.Config().DisplayTimezone)
	return timezoneVN
}

func formatPriceChanged(priceChanged float64) string {
	prefix := "up +"
	if priceChanged < 0 {
		prefix = "down "
	}
	return fmt.Sprintf("%s<strong>%v</strong>", prefix, fmt.Sprintf("%.2f%%", priceChanged))
}

type TelegramMessage struct {
	ChatId                int    `json:"chat_id"`
	Text                  string `json:"text"`
	ParseMode             string `json:"parse_mode"`
	DisableWebPagePreview bool   `json:"disable_web_page_preview"`
}

func sendTeleMsg(msg string) {
	buf := &bytes.Buffer{}
	json.NewEncoder(buf).Encode(&TelegramMessage{
		ChatId:                config.Config().BotChatID,
		Text:                  msg,
		ParseMode:             "html",
		DisableWebPagePreview: true,
	})

	res, err := http.Post(config.TelegramBotAPI(), "application/json", buf)
	if err != nil {
		fmt.Println(err)
	}
	if res.StatusCode != 200 {
		fmt.Println(res.Status)
	}
}

// Sample: RNDRUSDT is up +1.36% to $0.587 in 15m
func buildNotificationMessage(symbol string, percentChanged float64, timeframe int, openPrice float64, closedPrice float64) string {
	return fmt.Sprintf("%s is %s from <strong>$%v</strong> to <strong>$%v</strong> in %s",
		formatSymbol(symbol),
		formatPriceChanged(percentChanged),
		openPrice,
		closedPrice,
		formatTimeframe(timeframe),
	)
}

func formatTimeframe(timeframe int) string {
	if timeframe > 60 {
		d := time.Minute * time.Duration(timeframe)
		return fmt.Sprintf("%.2fh", d.Hours())
	}
	return fmt.Sprintf("%dm", timeframe)
}

var currencies = []string{"USDT", "USD", "BTC", "ETH", "BUSD", "USDC", "DAI", "TUSD", "USDN", "LUSD", "USDD", "USDP"}

func separatePair(pair string) (string, string) {
	for _, currency := range currencies {
		if strings.HasSuffix(pair, currency) {
			return strings.TrimSuffix(pair, currency), currency
		}
	}
	return pair, ""
}

func getExchangeTradeURL(exchange, pair string) string {
	coin, currency := separatePair(pair)

	switch exchange {
	case "BINANCE":
		return "https://www.binance.com/en/trade/" + coin + "_" + currency
	case "GATEIO":
		return "https://www.gate.io/trade/" + coin + "_" + currency
	case "COINBASE":
		return "https://pro.coinbase.com/trade/" + coin + "_" + currency
	case "KUCOIN":
		return "https://www.kucoin.com/vi/trade/" + coin + "-" + currency
	case "BITFINEX":
		return "https://trading.bitfinex.com/t/" + coin + "_" + currency
	case "HUOBI":
		return "https://www.huobi.com/en-us/exchange/" + coin + "_" + currency
	case "OKX", "OKEX":
		return "https://www.okx.com/vi/trade-spot/" + coin + "-" + currency
	case "FTX":
		return "https://ftx.com/trade/" + coin + "/" + currency
	case "MEXC":
		return "https://www.mexc.com/vi-VN/exchange/" + coin + "_" + currency
	default:
		return "https://www.tradingview.com/chart?symbol=" + exchange + ":" + pair
	}
}

func formatSymbol(symbol string) string {
	s := strings.Split(symbol, ":")
	exchange := s[0]
	pair := s[1]
	return fmt.Sprintf("<a href='%s'>%s</a>", getExchangeTradeURL(exchange, pair), pair)
}

// JSONSerializer json serializer
type JSONSerializer struct {
}

// Scan implements serializer interface
func (JSONSerializer) Scan(ctx context.Context, field *schema.Field, dst reflect.Value, dbValue interface{}) (err error) {
	fieldValue := reflect.New(field.FieldType)

	if dbValue != nil {
		var bytes []byte
		switch v := dbValue.(type) {
		case []byte:
			bytes = v
		case string:
			bytes = []byte(v)
		default:
			return fmt.Errorf("failed to unmarshal JSONB value: %#v", dbValue)
		}

		err = json.Unmarshal(bytes, fieldValue.Interface())
	}

	field.ReflectValueOf(ctx, dst).Set(fieldValue.Elem())
	return
}

// Value implements serializer interface
func (JSONSerializer) Value(ctx context.Context, field *schema.Field, dst reflect.Value, fieldValue interface{}) (interface{}, error) {
	return json.Marshal(fieldValue)
}

func formatLikeQuery(s string) string {
	return fmt.Sprintf("%%%s%%", s)
}

func calculatePercentChanged(openPrice, closePrice float64) (percentChanged float64) {
	if closePrice > openPrice {
		percentChanged = (closePrice - openPrice) / openPrice * 100
	} else {
		percentChanged = (closePrice - openPrice) / closePrice * 100
	}
	percentChanged = math.Round(percentChanged*100) / 100
	return
}

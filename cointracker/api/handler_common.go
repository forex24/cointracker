package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/canhlinh/cointracker/types"
)

func GetKlines(c *Context) (interface{}, error) {
	var form FormSearch
	if err := c.BindQuery(&form); err != nil {
		return nil, err
	}

	klines, total, err := c.App.DB.GetKlines(form.Text, form.Offset, form.Limit)
	if err != nil {
		return nil, err
	}
	return NewResponse(klines).SetTotal(total), nil
}

func GetSymbols(c *Context) (interface{}, error) {
	var form FormSearchSymbol
	if err := c.BindQuery(&form); err != nil {
		return nil, err
	}

	if len(form.IDs) == 1 {
		ts := strings.Split(form.IDs[0], ":")
		if len(ts) == 2 {
			form.MatchSymbol = form.IDs[0]
			form.Exchange = ts[0]
			form.Text = ts[1]
		}
	} else if len(form.Text) > 0 {
		ts := strings.Split(form.Text, ":")
		if len(ts) == 2 {
			form.MatchSymbol = form.Text
			form.Exchange = ts[0]
			form.Text = ts[1]
		}
	}

	query := url.Values{}
	query.Add("text", form.Text)
	query.Add("exchange", form.Exchange)
	query.Add("lang", "en")
	query.Add("type", "bitcoin,crypto")
	query.Add("domain", "production")

	r, err := http.Get("https://symbol-search.tradingview.com/symbol_search?" + query.Encode())
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.New(r.Status)
	}

	var symbols types.TVSymbols
	if err := json.NewDecoder(r.Body).Decode(&symbols); err != nil {
		return nil, err
	}
	return NewResponse(symbols.Format(form.MatchSymbol)).SetTotal(int64(len(symbols))), nil
}

func GetAlerts(c *Context) (interface{}, error) {
	var form FormSearch
	if err := c.BindQuery(&form); err != nil {
		return nil, err
	}

	alerts, total, err := c.App.DB.GetAlerts(form.Text, form.Offset, form.Limit)
	if err != nil {
		return nil, err
	}
	return NewResponse(alerts).SetTotal(total), nil
}

package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"strings"

	"github.com/canhlinh/cointracker/types"
	"github.com/gabriel-vasile/mimetype"
	log "github.com/sirupsen/logrus"
)

func DeleteAlertConfig(c *Context) (interface{}, error) {
	alertConfig, err := c.App.DB.GetAlertConfig(c.ParamInt("id"), false)
	if err != nil {
		return nil, err
	}

	if err := c.App.DB.DeleteAlertConfig(c.ParamInt("id")); err != nil {
		return nil, err
	}

	if err := c.App.RemoveSymbol(alertConfig.Symbol); err != nil {
		log.Error(err)
	}
	return NewResponse("DELETED"), nil
}

func DeleteAlertConfigs(c *Context) (interface{}, error) {
	var ids []int64
	json.NewDecoder(c.Request.Body).Decode(&ids)
	deletedAlertConfigs := []*types.AlertConfig{}

	for _, id := range ids {
		alertConfig, err := c.App.DB.GetAlertConfig(id, false)
		if err != nil {
			return nil, err
		}

		if err := c.App.DB.DeleteAlertConfig(id); err != nil {
			return nil, err
		}

		if err := c.App.RemoveSymbol(alertConfig.Symbol); err != nil {
			return nil, err
		}
		deletedAlertConfigs = append(deletedAlertConfigs, alertConfig)
	}
	return NewResponse(deletedAlertConfigs), nil
}

func CreateAlertConfigs(c *Context) (interface{}, error) {
	var alertConfig types.AlertConfig
	if err := c.BindJSON(&alertConfig); err != nil {
		return nil, err
	}
	if !isValidSymbol(alertConfig.Symbol) {
		return nil, errors.New("invalid symbol")
	}

	if err := c.App.CreateAlertConfig(&alertConfig); err != nil {
		return nil, err
	}
	return NewResponse(alertConfig), nil
}

func GetAlertConfigs(c *Context) (interface{}, error) {
	var form FormSearch
	if err := c.BindQuery(&form); err != nil {
		return nil, err
	}

	alertConfigs, total, err := c.App.DB.GetAlertConfigs(form.Text, form.Offset, form.Limit)
	if err != nil {
		return nil, err
	}
	return NewResponse(alertConfigs).SetTotal(total), nil
}

func GetAlertConfig(c *Context) (interface{}, error) {
	alertConfig, err := c.App.DB.GetAlertConfig(c.ParamInt("id"), false)
	if err != nil {
		return nil, err
	}
	return NewResponse(alertConfig), nil
}

func UpdateAlertConfig(c *Context) (interface{}, error) {
	tx := c.App.DB.Begin()
	defer tx.Commit()

	old, err := tx.GetAlertConfig(c.ParamInt("id"), true)
	if err != nil {
		return nil, err
	}

	var new types.AlertConfig
	if err := c.BindJSON(&new); err != nil {
		return nil, err
	}
	new.ID = old.ID

	if err := tx.UpdateAlertConfig(&new); err != nil {
		return nil, err
	}

	if old.Symbol != new.Symbol {
		if err := c.App.RemoveSymbol(old.Symbol); err != nil {
			log.Error(err)
		}
		if err := c.App.AddSymbol(new.Symbol); err != nil {
			log.Error(err)
		}
	}

	tx.Commit()
	return NewResponse(new), nil
}

func ImportSymbols(c *Context) (interface{}, error) {
	err := c.Request.ParseMultipartForm(32 << 20)
	if err != nil {
		return nil, err
	}
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buf := &bytes.Buffer{}
	if _, err := io.Copy(buf, file); err != nil {
		return nil, err
	}
	mime := mimetype.Detect(buf.Bytes())
	if !mime.Is("text/plain") {
		return nil, errors.New("file type not supported")
	}
	symbols := strings.Split(buf.String(), ",")
	alertConfigs := []*types.AlertConfig{}

	for _, symbol := range symbols {
		if !isValidSymbol(symbol) {
			continue
		}

		alertConfig := &types.AlertConfig{Symbol: symbol}
		if err := c.App.CreateAlertConfig(alertConfig); err != nil {
			return nil, err
		}
		alertConfigs = append(alertConfigs, alertConfig)
	}

	return NewResponse(alertConfigs), nil
}

func isValidSymbol(symbol string) bool {
	if len(symbol) == 0 {
		return false
	}
	s := strings.Split(symbol, ":")
	return len(s) == 2
}

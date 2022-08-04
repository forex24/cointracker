package api

import "github.com/canhlinh/cointracker/types"

func GetTimeframes(c *Context) (interface{}, error) {
	var form FormGetTimeframes
	if err := c.BindQuery(&form); err != nil {
		return nil, err
	}

	timeframes, err := c.App.DB.GetTimeframes(form.IDs, form.AlertConfigID)
	if err != nil {
		return nil, err
	}
	return NewResponse(timeframes).SetTotal(int64(len(timeframes))), nil
}

func GetTimeframe(c *Context) (interface{}, error) {
	timeframe, err := c.App.DB.GetTimeframe(c.ParamInt("id"))
	if err != nil {
		return nil, err
	}
	return NewResponse(timeframe), nil
}

func UpdateTimeframe(c *Context) (interface{}, error) {
	old, err := c.App.DB.GetTimeframe(c.ParamInt("id"))
	if err != nil {
		return nil, err
	}

	var new types.Timeframe
	if err := c.BindJSON(&new); err != nil {
		return nil, err
	}
	new.ID = old.ID

	if err := c.App.DB.UpdateTimeframe(&new); err != nil {
		return nil, err
	}
	return NewResponse(new), nil
}

package api

type FormSearchSymbol struct {
	IDs         []string `form:"ids"`
	Text        string   `form:"text"`
	Exchange    string   `form:"exchange"`
	MatchSymbol string   `form:"-"`
}

type FormGetTimeframes struct {
	IDs           []int `form:"ids"`
	AlertConfigID int64 `form:"alert_config_id"`
}

type FormSearch struct {
	Text   *string `form:"text"`
	Offset int     `form:"offset"`
	Limit  int     `form:"limit"`
}

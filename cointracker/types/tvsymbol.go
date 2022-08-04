package types

type TVSymbol struct {
	ID           string `json:"id"`
	Symbol       string `json:"symbol"`
	Exchange     string `json:"exchange"`
	CurrencyCode string `json:"currency_code"`
	Description  string `json:"description"`
	Format       string `json:"format"`
}

type TVSymbols []*TVSymbol

func (tvs TVSymbols) Format(matchSymbol string) TVSymbols {
	if len(matchSymbol) > 0 && len(tvs) > 0 {
		if tvs[0].Exchange+":"+tvs[0].Symbol == matchSymbol {
			tvs = tvs[:1]
		}
	}

	for _, symbol := range tvs {
		symbol.Format = symbol.Exchange + ":" + symbol.Symbol
		symbol.ID = symbol.Format
	}
	return tvs
}

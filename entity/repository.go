package entity

type GetCBRatesRequest struct {
	Country string
}

type GetCBRatesResponse struct {
	Rates *ExchangeRates
}

type GetExchangeRateRequest struct {
	Country          string
	BaseCurrencyID   string
	TargetCurrencyID string
	Amount           int
}

type GetExchangeRateResponse struct {
	Rate Rate
}

package entity

import (
	"time"
)

// GetExchangeRatesRequest is a request to get exchange rates from a central bank of provided country
type GetExchangeRatesRequest struct {
	Country string `json:"country,omitempty"`
}

// GetExchangeRatesResponse is a container with exchange rates for the external API request
type GetExchangeRatesResponse struct {
	Rates map[string]Rate `json:"rates"`
}

// ExchangeRates is a container to store internal exchange rates data for a single central bank
type ExchangeRates struct {
	Country    string
	DateLoaded string
	TimeZone   *time.Location
	Rates      map[string]Rate
}

// Rate is a container for a single exchange rate
// E.g. 25.159600 THAI bt = 100 JPY means
// Nominal = 100,
// BaseCurrency = THB,
// TargetCurrency = JPY,
// RateTargetToBase = 25159600
type Rate struct {
	Nominal          int     `json:"nominal,omitempty"` // amount of target currency to be used for ratio
	BaseCurrency     string  `json:"base_currency,omitempty"`
	TargetCurrency   string  `json:"target_currency,omitempty"`
	RateTargetToBase float64 `json:"rate_target_to_base,omitempty"` // contains 6 digits decimal ratio in bigint.
}

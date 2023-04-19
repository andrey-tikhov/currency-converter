package entity

// ConvertCurrencyRequest is a container to store the currency conversion request
// Country represents the central bank that is expected as source of exchange rates (default will be used if omitted)
// The request represents the following question: How much in TargetCurrency will be the Amount of SourceCurrency
type ConvertCurrencyRequest struct {
	Country        *string `json:"country,omitempty"`
	SourceCurrency string  `json:"source_currency,omitempty"`
	TargetCurrency string  `json:"target_currency,omitempty"`
	Amount         int     `json:"amount,omitempty"`
}

// ConvertCurrencyResponse represents the resulted Amount of SourceCurrency
type ConvertCurrencyResponse struct {
	Amount float64 `json:"amount"`
}

package mapper

import (
	"errors"
	"my_go/entity"
)

// GetExchangeRatesRequestToGetCBRatesRequest converts incoming controller request to
// request consumed by repository
func GetExchangeRatesRequestToGetCBRatesRequest(
	r *entity.GetExchangeRatesRequest,
) (*entity.GetCBRatesRequest, error) {
	if r == nil {
		return nil, errors.New("nil GetExchangeRatesRequest")
	}
	return &entity.GetCBRatesRequest{
		Country: r.Country,
	}, nil
}

// GetCBRRatesResponseToGetExchangeRatesResponse converts repository response to GetExchangeRatesResponse
// that will be marshalled and used for external API
func GetCBRRatesResponseToGetExchangeRatesResponse(
	r *entity.GetCBRatesResponse,
) (*entity.GetExchangeRatesResponse, error) {
	if r == nil {
		return nil, errors.New("nil GetCBRatesResponse")
	}
	return &entity.GetExchangeRatesResponse{
		Rates: r.Rates.Rates,
	}, nil
}

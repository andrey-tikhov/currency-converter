package mapper

import (
	"encoding/json"
	"errors"
	"fmt"
	"my_go/entity"
	"strconv"
)

// BodyToConvertCurrencyRequest converts the http response body to internal entity.ConvertCurrencyRequest
func BodyToConvertCurrencyRequest(body []byte) (*entity.ConvertCurrencyRequest, error) {
	var r entity.ConvertCurrencyRequest
	err := json.Unmarshal(body, &r)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %s", err)
	}
	return &r, nil
}

// ConvertCurrencyResponseToBytes converts internal entity.ConvertCurrencyResponse to http response body
func ConvertCurrencyResponseToBytes(response *entity.ConvertCurrencyResponse) ([]byte, error) {
	b, err := json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal: %s", err) // unreachable in tests
	}
	return b, nil
}

// CBRRatesAndGetExchangeRateRequestToGetExchangeRateResponse converts rates amd get exchange rate request to
// entity.GetExchangeRateResponse
func CBRRatesAndGetExchangeRateRequestToGetExchangeRateResponse(
	r *entity.ExchangeRates,
	req *entity.GetExchangeRateRequest,
) (*entity.GetExchangeRateResponse, error) {
	if r == nil {
		return nil, errors.New("nil ExchangeRates")
	}
	if req == nil {
		return nil, errors.New("nil GetExchangeRateRequest")
	}
	targetRate, ok := r.Rates[req.TargetCurrencyID]
	if !ok {
		return nil, fmt.Errorf(
			"targetRate, currency %s not supported by CB %s", req.TargetCurrencyID, req.Country,
		)
	}
	baseRate, ok := r.Rates[req.BaseCurrencyID]
	if !ok {
		return nil, fmt.Errorf("baseRate, currency %s not supported by CB %s", req.TargetCurrencyID, req.Country)
	}

	newRawRate := float64(targetRate.Nominal) / targetRate.RateTargetToBase *
		(baseRate.RateTargetToBase / float64(baseRate.Nominal)) * float64(req.Amount)
	strRate := fmt.Sprintf("%.4f", newRawRate)
	rate, _ := strconv.ParseFloat(strRate, 64)

	return &entity.GetExchangeRateResponse{
		Rate: entity.Rate{
			Nominal:          req.Amount,
			BaseCurrency:     req.BaseCurrencyID,
			TargetCurrency:   req.TargetCurrencyID,
			RateTargetToBase: rate, // int64 overflow
		},
	}, nil
}

// ConvertCurrencyRequestToGetExchangeRateRequest converts entity.ConvertCurrencyRequest
// to entity.GetExchangeRateRequest. defaultCB is used as a fallback if country if not provided.
func ConvertCurrencyRequestToGetExchangeRateRequest(
	req *entity.ConvertCurrencyRequest,
	defaultCB string,
) (*entity.GetExchangeRateRequest, error) {
	if req == nil {
		return nil, fmt.Errorf("nil ConvertCurrencyRequest")
	}
	country := defaultCB
	if req.Country != nil {
		country = *req.Country
	}
	return &entity.GetExchangeRateRequest{
		Country:          country,
		BaseCurrencyID:   req.SourceCurrency,
		TargetCurrencyID: req.TargetCurrency,
		Amount:           req.Amount,
	}, nil
}

// GetExchangeRateResponseToConvertCurrencyResponse converts entity.GetExchangeRateResponse
// to entity.ConvertCurrencyResponse
func GetExchangeRateResponseToConvertCurrencyResponse(
	r *entity.GetExchangeRateResponse,
) (*entity.ConvertCurrencyResponse, error) {
	if r == nil {
		return nil, fmt.Errorf("nil Rate")
	}
	return &entity.ConvertCurrencyResponse{
		Amount: r.Rate.RateTargetToBase,
	}, nil
}

package mapper

import (
	"github.com/stretchr/testify/assert"
	"my_go/entity"
	"testing"
	"time"
)

func TestGetExchangeRatesRequestToGetCBRatesRequest(t *testing.T) {
	getExchangeRatesRequest := &entity.GetExchangeRatesRequest{
		Country: "russia",
	}
	type args struct {
		r *entity.GetExchangeRatesRequest
	}
	tests := []struct {
		name      string
		args      args
		want      *entity.GetCBRatesRequest
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Happy path",
			args: args{
				r: getExchangeRatesRequest,
			},
			want: &entity.GetCBRatesRequest{
				Country: "russia",
			},
			assertion: assert.NoError,
		},
		{
			name: "nil request",
			args: args{
				r: nil,
			},
			want:      nil,
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetExchangeRatesRequestToGetCBRatesRequest(tt.args.r)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetCBRRatesResponseToGetExchangeRatesResponse(t *testing.T) {
	ruLoc, _ := time.LoadLocation("Europe/Moscow")
	response := &entity.GetCBRatesResponse{
		Rates: &entity.ExchangeRates{
			Country:    "russia",
			DateLoaded: "2023-01-01",
			TimeZone:   ruLoc,
			Rates: map[string]entity.Rate{
				"USD": {
					Nominal:          100,
					BaseCurrency:     "RUR",
					TargetCurrency:   "USD",
					RateTargetToBase: 123.56,
				},
				"JPY": {
					Nominal:          10,
					BaseCurrency:     "RUR",
					TargetCurrency:   "JPY",
					RateTargetToBase: 987.65,
				},
			},
		},
	}
	type args struct {
		r *entity.GetCBRatesResponse
	}
	tests := []struct {
		name      string
		args      args
		want      *entity.GetExchangeRatesResponse
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Happy path",
			args: args{
				r: response,
			},
			want: &entity.GetExchangeRatesResponse{
				Rates: map[string]entity.Rate{
					"USD": {
						Nominal:          100,
						BaseCurrency:     "RUR",
						TargetCurrency:   "USD",
						RateTargetToBase: 123.56,
					},
					"JPY": {
						Nominal:          10,
						BaseCurrency:     "RUR",
						TargetCurrency:   "JPY",
						RateTargetToBase: 987.65,
					},
				},
			},
			assertion: assert.NoError,
		},
		{
			name: "nil input",
			args: args{
				r: nil,
			},
			want:      nil,
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCBRRatesResponseToGetExchangeRatesResponse(tt.args.r)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

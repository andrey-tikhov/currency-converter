package mapper

import (
	"github.com/stretchr/testify/assert"
	"my_go/entity"
	"my_go/utils"
	"testing"
	"time"
)

func TestCBRRatesAndGetExchangeRateRequestToGetExchangeRateResponse(t *testing.T) {
	thaiTZ, _ := time.LoadLocation("Asia/Bangkok")
	type args struct {
		r   *entity.ExchangeRates
		req *entity.GetExchangeRateRequest
	}
	tests := []struct {
		name      string
		args      args
		want      *entity.GetExchangeRateResponse
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Happy Path",
			args: args{
				r: &entity.ExchangeRates{
					Country:    "thailand",
					DateLoaded: "2022-01-01",
					TimeZone:   thaiTZ,
					Rates: map[string]entity.Rate{
						"USD": {
							Nominal:          1,
							BaseCurrency:     "THB",
							TargetCurrency:   "USD",
							RateTargetToBase: 34.079000,
						},
						"JPY": {
							Nominal:          100,
							BaseCurrency:     "THB",
							TargetCurrency:   "JPY",
							RateTargetToBase: 25.159600,
						},
					},
				},
				req: &entity.GetExchangeRateRequest{
					Country:          "thailand",
					BaseCurrencyID:   "USD",
					TargetCurrencyID: "JPY",
					Amount:           25,
				},
			},
			want: &entity.GetExchangeRateResponse{
				Rate: entity.Rate{
					Nominal:          25,
					BaseCurrency:     "USD",
					TargetCurrency:   "JPY",
					RateTargetToBase: 3386.282,
				},
			},
			assertion: assert.NoError,
		},
		{
			name: "nil request",
			args: args{
				r: nil,
				req: &entity.GetExchangeRateRequest{
					Country:          "thailand",
					BaseCurrencyID:   "USD",
					TargetCurrencyID: "JPY",
					Amount:           25,
				},
			},
			want:      nil,
			assertion: assert.Error,
		},
		{
			name: "nil rates",
			args: args{
				r: &entity.ExchangeRates{
					Rates: nil,
				},
				req: nil,
			},
			assertion: assert.Error,
		},
		{
			name: "Base Currency not in rates",
			args: args{
				r: &entity.ExchangeRates{
					Country:    "thailand",
					DateLoaded: "2022-01-01",
					TimeZone:   thaiTZ,
					Rates: map[string]entity.Rate{
						"JPY": {
							Nominal:          100,
							BaseCurrency:     "THB",
							TargetCurrency:   "JPY",
							RateTargetToBase: 25.159600,
						},
					},
				},
				req: &entity.GetExchangeRateRequest{
					Country:          "thailand",
					BaseCurrencyID:   "USD",
					TargetCurrencyID: "JPY",
					Amount:           25,
				},
			},
			want:      nil,
			assertion: assert.Error,
		},
		{
			name: "Target Currency not in rates",
			args: args{
				r: &entity.ExchangeRates{
					Country:    "thailand",
					DateLoaded: "2022-01-01",
					TimeZone:   thaiTZ,
					Rates: map[string]entity.Rate{
						"USD": {
							Nominal:          1,
							BaseCurrency:     "THB",
							TargetCurrency:   "USD",
							RateTargetToBase: 34.079000,
						},
					},
				},
				req: &entity.GetExchangeRateRequest{
					Country:          "thailand",
					BaseCurrencyID:   "USD",
					TargetCurrencyID: "JPY",
					Amount:           25,
				},
			},
			want:      nil,
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CBRRatesAndGetExchangeRateRequestToGetExchangeRateResponse(tt.args.r, tt.args.req)
			assert.Equal(t, tt.want, got)
			tt.assertion(t, err)
		})
	}
}

func TestBodyToConvertCurrencyRequest(t *testing.T) {
	correctJSON := []byte(`{"country":"russia","source_currency":"RUR","target_currency":"USD","amount":100}`)
	type args struct {
		body []byte
	}
	tests := []struct {
		name      string
		args      args
		want      *entity.ConvertCurrencyRequest
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Happy path",
			args: args{
				body: correctJSON,
			},
			want: &entity.ConvertCurrencyRequest{
				Country:        utils.ToPointer("russia"),
				SourceCurrency: "RUR",
				TargetCurrency: "USD",
				Amount:         100,
			},
			assertion: assert.NoError,
		},
		{
			name: "bad json",
			args: args{
				body: []byte(`{"swrgwrgr:`),
			},
			want:      nil,
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BodyToConvertCurrencyRequest(tt.args.body)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestConvertCurrencyResponseToBytes(t *testing.T) {
	type args struct {
		response *entity.ConvertCurrencyResponse
	}
	tests := []struct {
		name      string
		args      args
		want      []byte
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Happy path",
			args: args{
				response: &entity.ConvertCurrencyResponse{
					Amount: 123.34,
				},
			},
			want:      []byte(`{"amount":123.34}`),
			assertion: assert.NoError,
		},
		{
			name: "nil response",
			args: args{
				response: nil,
			},
			want:      []byte(`null`),
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertCurrencyResponseToBytes(tt.args.response)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestConvertCurrencyRequestToGetExchangeRateRequest(t *testing.T) {
	requestWithCountry := &entity.ConvertCurrencyRequest{
		Country:        utils.ToPointer("thailand"),
		SourceCurrency: "JPY",
		TargetCurrency: "USD",
		Amount:         12,
	}
	requestWithoutCountry := &entity.ConvertCurrencyRequest{
		Country:        nil,
		SourceCurrency: "JPY",
		TargetCurrency: "USD",
		Amount:         12,
	}
	type args struct {
		req       *entity.ConvertCurrencyRequest
		defaultCB string
	}
	tests := []struct {
		name      string
		args      args
		want      *entity.GetExchangeRateRequest
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Happy path",
			args: args{
				req:       requestWithCountry,
				defaultCB: "russia",
			},
			want: &entity.GetExchangeRateRequest{
				Country:          "thailand",
				BaseCurrencyID:   "JPY",
				TargetCurrencyID: "USD",
				Amount:           12,
			},
			assertion: assert.NoError,
		},
		{
			name: "Happy path, no country, fallback",
			args: args{
				req:       requestWithoutCountry,
				defaultCB: "russia",
			},
			want: &entity.GetExchangeRateRequest{
				Country:          "russia",
				BaseCurrencyID:   "JPY",
				TargetCurrencyID: "USD",
				Amount:           12,
			},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertCurrencyRequestToGetExchangeRateRequest(tt.args.req, tt.args.defaultCB)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetCBRRateToConvertCurrencyResponse(t *testing.T) {
	rate := entity.Rate{
		Nominal:          100,
		BaseCurrency:     "RUR",
		TargetCurrency:   "USD",
		RateTargetToBase: 123.56,
	}
	type args struct {
		r *entity.GetExchangeRateResponse
	}
	tests := []struct {
		name      string
		args      args
		want      *entity.ConvertCurrencyResponse
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Happy path",
			args: args{
				r: &entity.GetExchangeRateResponse{
					Rate: rate,
				},
			},
			want: &entity.ConvertCurrencyResponse{
				Amount: 123.56,
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
			got, err := GetExchangeRateResponseToConvertCurrencyResponse(tt.args.r)
			tt.assertion(t, err)
			assert.Equalf(t, tt.want, got, "GetExchangeRateResponseToConvertCurrencyResponse(%v)", tt.args.r)
		})
	}
}

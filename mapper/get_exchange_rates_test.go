package mapper

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"my_go/entity"
	"testing"
)

func TestBodyToGetExchangeRateRequest(t *testing.T) {
	correctBody := []byte(`{"country":"thailand"}`)
	type args struct {
		body []byte
	}
	tests := []struct {
		name      string
		args      args
		want      *entity.GetExchangeRatesRequest
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Happy path",
			args: args{
				body: correctBody,
			},
			want: &entity.GetExchangeRatesRequest{
				Country: "thailand",
			},
			assertion: assert.NoError,
		},
		{
			name: "bad json",
			args: args{
				body: []byte(`{s"r`),
			},
			want:      nil,
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BodyToGetExchangeRatesRequest(tt.args.body)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetCBRateResponseToBytes(t *testing.T) {
	response := &entity.GetExchangeRatesResponse{
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
	}
	type args struct {
		r *entity.GetExchangeRatesResponse
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
				r: response,
			},
			want:      []byte(`{"rates":{"USD":{"nominal":100,"base_currency":"RUR","target_currency":"USD","rate_target_to_base":123.56},"JPY":{"nominal":10,"base_currency":"RUR","target_currency":"JPY","rate_target_to_base":987.65}}}`),
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetExchangeRatesResponseToBytes(tt.args.r)
			tt.assertion(t, err)
			var gotStruct entity.GetExchangeRateResponse
			err = json.Unmarshal(got, &gotStruct)
			assert.NoError(t, err)
			var expectedStruct entity.GetExchangeRateResponse
			err = json.Unmarshal(tt.want, &expectedStruct)
			assert.Equal(t, expectedStruct, gotStruct)
		})
	}
}

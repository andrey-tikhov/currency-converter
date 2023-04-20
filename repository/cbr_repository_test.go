package repository

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"my_go/entity"
	"my_go/gateway"
	russiagatewaymock "my_go/mocks/gateway/russia"
	thailandgatewaymock "my_go/mocks/gateway/thailand"
	"sync"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRussiaCB := russiagatewaymock.NewMockGateway(ctrl)
	mockThailandCB := thailandgatewaymock.NewMockGateway(ctrl)
	c, err := New(Params{
		ThaiGateway:   mockRussiaCB,
		RussiaGateway: mockThailandCB,
	})
	assert.NoError(t, err)
	assert.NotNil(t, c)
}

func Test_cbr_GetCBRates(t *testing.T) {
	timeNow := func() time.Time {
		return time.Unix(100, 0)
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ruTZ, _ := time.LoadLocation("Europe/Moscow")
	ruGatewayResponse := &entity.ExchangeRates{
		Country:    "russia",
		DateLoaded: "2023-01-01",
		TimeZone:   ruTZ,
		Rates: map[string]entity.Rate{
			"USD": {
				Nominal:          100,
				BaseCurrency:     "RUB",
				TargetCurrency:   "USD",
				RateTargetToBase: 123.56,
			},
			"JPY": {
				Nominal:          10,
				BaseCurrency:     "RUB",
				TargetCurrency:   "JPY",
				RateTargetToBase: 987.65,
			},
		},
	}
	thTZ, _ := time.LoadLocation("Asia/Bangkok")
	thGatewayResponse := &entity.ExchangeRates{
		Country:    "thailand",
		DateLoaded: "2023-01-01",
		TimeZone:   thTZ,
		Rates: map[string]entity.Rate{
			"THB": {
				Nominal:          1,
				RateTargetToBase: 1,
				BaseCurrency:     "THB",
				TargetCurrency:   "THB",
			},
			"GBP": {
				Nominal:          1,
				RateTargetToBase: 42.5171,
				BaseCurrency:     "THB",
				TargetCurrency:   "GBP",
			},
			"USD": {
				Nominal:          1,
				RateTargetToBase: 34.28235,
				BaseCurrency:     "THB",
				TargetCurrency:   "USD",
			},
		},
	}

	type fields struct {
		TimeNow    func() time.Time
		RatesCache map[string]entity.ExchangeRates
	}
	type args struct {
		req *entity.GetCBRatesRequest
	}
	type mockCBGateway struct {
		res *entity.ExchangeRates
		err error
	}
	type expectedNumberOfGatewayCalls struct {
		RussiaGetCBRates   int
		ThailandGetCBRates int
	}
	tests := []struct {
		name                      string
		fields                    fields
		args                      args
		mockRussiaCBGateway       *mockCBGateway
		mockThailandCBGateway     *mockCBGateway
		expectedNumOfGatewayCalls expectedNumberOfGatewayCalls
		want                      *entity.GetCBRatesResponse
		assertion                 assert.ErrorAssertionFunc
	}{
		{
			name: "Happy path. Russia",
			fields: fields{
				TimeNow: timeNow,
				RatesCache: map[string]entity.ExchangeRates{
					entity.Russia:   {},
					entity.Thailand: {},
				},
			},
			args: args{
				req: &entity.GetCBRatesRequest{
					Country: "russia",
				},
			},
			mockRussiaCBGateway: &mockCBGateway{
				res: ruGatewayResponse,
				err: nil,
			},
			expectedNumOfGatewayCalls: expectedNumberOfGatewayCalls{
				RussiaGetCBRates: 1,
			},
			want: &entity.GetCBRatesResponse{
				Rates: ruGatewayResponse,
			},
			assertion: assert.NoError,
		},
		{
			name: "Happy path. Thailand",
			fields: fields{
				RatesCache: map[string]entity.ExchangeRates{
					entity.Russia:   {},
					entity.Thailand: {},
				},
			},
			args: args{
				req: &entity.GetCBRatesRequest{
					Country: "thailand",
				},
			},
			mockThailandCBGateway: &mockCBGateway{
				res: thGatewayResponse,
				err: nil,
			},
			expectedNumOfGatewayCalls: expectedNumberOfGatewayCalls{
				ThailandGetCBRates: 1,
			},
			want: &entity.GetCBRatesResponse{
				Rates: thGatewayResponse,
			},
			assertion: assert.NoError,
		},
		{
			name: "country key is not supported",
			fields: fields{
				RatesCache: map[string]entity.ExchangeRates{},
			},
			args: args{
				req: &entity.GetCBRatesRequest{
					Country: "some unsupported country",
				},
			},
			want:      nil,
			assertion: assert.Error,
		},
		{
			name: "nil request",
			fields: fields{
				TimeNow: timeNow,
				RatesCache: map[string]entity.ExchangeRates{
					entity.Russia:   {},
					entity.Thailand: {},
				},
			},
			args: args{
				req: nil,
			},
			want:      nil,
			assertion: assert.Error,
		},
		{
			name: "cache reloading fails",
			fields: fields{
				TimeNow: timeNow,
				RatesCache: map[string]entity.ExchangeRates{
					entity.Russia:   {},
					entity.Thailand: {},
				},
			},
			args: args{
				req: &entity.GetCBRatesRequest{
					Country: "russia",
				},
			},
			mockRussiaCBGateway: &mockCBGateway{
				res: nil,
				err: errors.New("some error"),
			},
			expectedNumOfGatewayCalls: expectedNumberOfGatewayCalls{
				RussiaGetCBRates: 1,
			},
			want:      nil,
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockRussiaCB := russiagatewaymock.NewMockGateway(ctrl)
			if tt.mockRussiaCBGateway != nil {
				mockRussiaCB.
					EXPECT().
					GetCBRRates(ctx).
					MaxTimes(tt.expectedNumOfGatewayCalls.RussiaGetCBRates).
					MinTimes(tt.expectedNumOfGatewayCalls.RussiaGetCBRates).
					Return(tt.mockRussiaCBGateway.res, tt.mockRussiaCBGateway.err)
			}
			mockThailandCB := thailandgatewaymock.NewMockGateway(ctrl)
			if tt.mockThailandCBGateway != nil {
				mockThailandCB.
					EXPECT().
					GetCBRRates(ctx).
					MaxTimes(tt.expectedNumOfGatewayCalls.ThailandGetCBRates).
					MinTimes(tt.expectedNumOfGatewayCalls.ThailandGetCBRates).
					Return(tt.mockThailandCBGateway.res, tt.mockThailandCBGateway.err)
			}

			c := &cbr{
				RWMutex: sync.RWMutex{},
				TimeNow: tt.fields.TimeNow,
				Gateways: map[string]gateway.CBGateway{
					entity.Russia:   mockRussiaCB,
					entity.Thailand: mockThailandCB,
				},
				RatesCache: tt.fields.RatesCache,
			}
			got, err := c.GetCBRates(ctx, tt.args.req)

			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_cbr_needsRefresh(t *testing.T) {
	thTZ, _ := time.LoadLocation("Asia/Bangkok")
	rates := &entity.ExchangeRates{
		Country:    "thailand",
		DateLoaded: "1970-01-01",
		TimeZone:   thTZ,
		Rates: map[string]entity.Rate{
			"THB": {
				Nominal:          1,
				RateTargetToBase: 1,
				BaseCurrency:     "THB",
				TargetCurrency:   "THB",
			},
			"GBP": {
				Nominal:          1,
				RateTargetToBase: 42.5171,
				BaseCurrency:     "THB",
				TargetCurrency:   "GBP",
			},
			"USD": {
				Nominal:          1,
				RateTargetToBase: 34.28235,
				BaseCurrency:     "THB",
				TargetCurrency:   "USD",
			},
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	type fields struct {
		TimeNow    func() time.Time
		RatesCache map[string]entity.ExchangeRates
	}
	type args struct {
		r *entity.ExchangeRates
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Happy path. Rates need refresh",
			fields: fields{
				TimeNow: func() time.Time {
					return time.Unix(100000000, 0)
				},
			},
			args: args{
				r: rates,
			},
			want: true,
		},
		{
			name: "Happy path. Rates do not need refresh",
			fields: fields{
				TimeNow: func() time.Time {
					return time.Unix(100, 0)
				},
			},
			args: args{
				r: rates,
			},
			want: false,
		},
		{
			name: "nil rates",
			fields: fields{
				TimeNow: func() time.Time {
					return time.Unix(100, 0)
				},
			},
			args: args{
				r: nil,
			},
			want: true,
		},
		{
			name: "Rates with empty date loaded",
			fields: fields{
				TimeNow: func() time.Time {
					return time.Unix(100, 0)
				},
			},
			args: args{
				r: &entity.ExchangeRates{},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRussiaCB := russiagatewaymock.NewMockGateway(ctrl)
			mockThailandCB := thailandgatewaymock.NewMockGateway(ctrl)
			c := &cbr{
				RWMutex: sync.RWMutex{},
				TimeNow: tt.fields.TimeNow,
				Gateways: map[string]gateway.CBGateway{
					entity.Russia:   mockRussiaCB,
					entity.Thailand: mockThailandCB,
				},
				RatesCache: tt.fields.RatesCache,
			}
			assert.Equal(t, tt.want, c.needsRefresh(tt.args.r))
		})
	}
}

func Test_cbr_reloadCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	thTZ, _ := time.LoadLocation("Asia/Bangkok")
	rates := entity.ExchangeRates{
		Country:    "thailand",
		DateLoaded: "1970-01-01",
		TimeZone:   thTZ,
		Rates: map[string]entity.Rate{
			"RUB": {
				Nominal:          1,
				RateTargetToBase: 1,
				BaseCurrency:     "RUB",
				TargetCurrency:   "RUB",
			},
			"GBP": {
				Nominal:          1,
				RateTargetToBase: 42.5171,
				BaseCurrency:     "THB",
				TargetCurrency:   "GBP",
			},
			"USD": {
				Nominal:          1,
				RateTargetToBase: 34.28235,
				BaseCurrency:     "THB",
				TargetCurrency:   "USD",
			},
		},
	}

	otherRates := entity.ExchangeRates{
		Country:    "thailand",
		DateLoaded: "1973-03-03",
		TimeZone:   thTZ,
		Rates: map[string]entity.Rate{
			"RUB": {
				Nominal:          1,
				RateTargetToBase: 1,
				BaseCurrency:     "RUB",
				TargetCurrency:   "RUB",
			},
			"GBP": {
				Nominal:          1,
				RateTargetToBase: 42.5171,
				BaseCurrency:     "THB",
				TargetCurrency:   "GBP",
			},
			"USD": {
				Nominal:          1,
				RateTargetToBase: 34.28235,
				BaseCurrency:     "THB",
				TargetCurrency:   "USD",
			},
		},
	}

	type fields struct {
		TimeNow    func() time.Time
		RatesCache map[string]entity.ExchangeRates
	}
	type args struct {
		country string
	}
	type mockCBGateway struct {
		res *entity.ExchangeRates
		err error
	}
	tests := []struct {
		name                  string
		fields                fields
		args                  args
		mockRussiaCBGateway   *mockCBGateway
		mockThailandCBGateway *mockCBGateway
		expectedCache         map[string]entity.ExchangeRates
		assertion             assert.ErrorAssertionFunc
	}{
		{
			name: "Happy path. Rates already in cache and do not need refresh",
			fields: fields{
				TimeNow: func() time.Time {
					return time.Unix(100, 0)
				},
				RatesCache: map[string]entity.ExchangeRates{
					"russia": rates,
				},
			},
			args: args{
				country: "russia",
			},
			expectedCache: map[string]entity.ExchangeRates{
				"russia": rates,
			},
			assertion: assert.NoError,
		},
		{
			name: "Happy path. Rates already in cache and need refresh. Russia.",
			fields: fields{
				TimeNow: func() time.Time {
					return time.Unix(100000000, 0)
				},
				RatesCache: map[string]entity.ExchangeRates{
					"russia": rates,
				},
			},
			mockRussiaCBGateway: &mockCBGateway{
				res: &otherRates,
				err: nil,
			},
			args: args{
				country: "russia",
			},
			expectedCache: map[string]entity.ExchangeRates{
				"russia": otherRates,
			},
			assertion: assert.NoError,
		},
		{
			name: "Happy path. Rates already in cache and need refresh. Thailand",
			fields: fields{
				TimeNow: func() time.Time {
					return time.Unix(100000000, 0)
				},
				RatesCache: map[string]entity.ExchangeRates{
					"thailand": rates,
				},
			},
			mockThailandCBGateway: &mockCBGateway{
				res: &otherRates,
				err: nil,
			},
			args: args{
				country: "thailand",
			},
			expectedCache: map[string]entity.ExchangeRates{
				"thailand": otherRates,
			},
			assertion: assert.NoError,
		},
		{
			name: "Gateway for country doesn't exist",
			fields: fields{
				TimeNow: func() time.Time {
					return time.Unix(100, 0)
				},
				RatesCache: map[string]entity.ExchangeRates{
					"russia": rates,
				},
			},
			args: args{
				country: "some country",
			},
			expectedCache: map[string]entity.ExchangeRates{
				"russia": rates,
			},
			assertion: assert.Error,
		},
		{
			name: "Rates already in cache and need refresh. Gateway fails.",
			fields: fields{
				TimeNow: func() time.Time {
					return time.Unix(100000000, 0)
				},
				RatesCache: map[string]entity.ExchangeRates{
					"russia": rates,
				},
			},
			mockRussiaCBGateway: &mockCBGateway{
				res: nil,
				err: errors.New("some error"),
			},
			args: args{
				country: "russia",
			},
			expectedCache: map[string]entity.ExchangeRates{
				"russia": rates,
			},
			assertion: assert.Error,
		},
		{
			name: "Rates already in cache and need refresh. Gateway returns nil and no error.",
			fields: fields{
				TimeNow: func() time.Time {
					return time.Unix(100000000, 0)
				},
				RatesCache: map[string]entity.ExchangeRates{
					"russia": rates,
				},
			},
			mockRussiaCBGateway: &mockCBGateway{
				res: nil,
				err: nil,
			},
			args: args{
				country: "russia",
			},
			expectedCache: map[string]entity.ExchangeRates{
				"russia": rates,
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockRussiaCB := russiagatewaymock.NewMockGateway(ctrl)
			if tt.mockRussiaCBGateway != nil {
				mockRussiaCB.
					EXPECT().
					GetCBRRates(ctx).
					Return(tt.mockRussiaCBGateway.res, tt.mockRussiaCBGateway.err)
			}
			mockThailandCB := thailandgatewaymock.NewMockGateway(ctrl)
			if tt.mockThailandCBGateway != nil {
				mockThailandCB.
					EXPECT().
					GetCBRRates(ctx).
					Return(tt.mockThailandCBGateway.res, tt.mockThailandCBGateway.err)
			}
			c := &cbr{
				RWMutex: sync.RWMutex{},
				TimeNow: tt.fields.TimeNow,
				Gateways: map[string]gateway.CBGateway{
					entity.Russia:   mockRussiaCB,
					entity.Thailand: mockThailandCB,
				},
				RatesCache: tt.fields.RatesCache,
			}
			tt.assertion(t, c.reloadCache(ctx, tt.args.country))
			assert.Equal(t, tt.expectedCache, c.RatesCache)
		})
	}
}

func Test_cbr_GetExchangeRate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ruTZ, _ := time.LoadLocation("Europe/Moscow")
	rates := entity.ExchangeRates{
		Country:    "thailand",
		DateLoaded: "1970-01-01",
		TimeZone:   ruTZ,
		Rates: map[string]entity.Rate{
			"RUB": {
				Nominal:          1,
				RateTargetToBase: 1,
				BaseCurrency:     "RUB",
				TargetCurrency:   "RUB",
			},
			"GBP": {
				Nominal:          1,
				RateTargetToBase: 42.5171,
				BaseCurrency:     "RUB",
				TargetCurrency:   "GBP",
			},
			"USD": {
				Nominal:          1,
				RateTargetToBase: 34.28235,
				BaseCurrency:     "RUB",
				TargetCurrency:   "USD",
			},
		},
	}

	type fields struct {
		TimeNow    func() time.Time
		RatesCache map[string]entity.ExchangeRates
	}
	type args struct {
		req *entity.GetExchangeRateRequest
	}
	type mockCBGateway struct {
		res *entity.ExchangeRates
		err error
	}
	tests := []struct {
		name                  string
		fields                fields
		args                  args
		mockRussiaCBGateway   *mockCBGateway
		mockThailandCBGateway *mockCBGateway
		want                  *entity.GetExchangeRateResponse
		assertion             assert.ErrorAssertionFunc
	}{
		{
			name: "Happy path, rates in cache",
			fields: fields{
				TimeNow: func() time.Time {
					return time.Unix(100, 0)
				},
				RatesCache: map[string]entity.ExchangeRates{
					"russia": rates,
				},
			},
			args: args{
				req: &entity.GetExchangeRateRequest{
					Country:          "russia",
					BaseCurrencyID:   "GBP",
					TargetCurrencyID: "USD",
					Amount:           10,
				},
			},
			want: &entity.GetExchangeRateResponse{
				Rate: entity.Rate{
					Nominal:          10,
					BaseCurrency:     "GBP",
					TargetCurrency:   "USD",
					RateTargetToBase: 12.402,
				},
			},
			assertion: assert.NoError,
		},
		{
			name: "nil request",
			fields: fields{
				TimeNow: func() time.Time {
					return time.Unix(100, 0)
				},
				RatesCache: map[string]entity.ExchangeRates{
					"russia": rates,
				},
			},
			args: args{
				req: nil,
			},
			assertion: assert.Error,
		},
		{
			name: "get cbr rate fails, country not supported",
			fields: fields{
				TimeNow: func() time.Time {
					return time.Unix(100, 0)
				},
				RatesCache: map[string]entity.ExchangeRates{
					"russia": rates,
				},
			},
			args: args{
				req: &entity.GetExchangeRateRequest{
					Country:          "some random country",
					BaseCurrencyID:   "GBP",
					TargetCurrencyID: "USD",
					Amount:           10,
				},
			},
			assertion: assert.Error,
		},
		{
			name: "cache refreshed with nil rates from gateway, no error",
			fields: fields{
				TimeNow: func() time.Time {
					return time.Unix(100, 0)
				},
				RatesCache: map[string]entity.ExchangeRates{},
			},
			args: args{
				req: &entity.GetExchangeRateRequest{
					Country:          "russia",
					BaseCurrencyID:   "GBP",
					TargetCurrencyID: "USD",
					Amount:           10,
				},
			},
			mockRussiaCBGateway: &mockCBGateway{
				res: nil,
				err: nil,
			},
			assertion: assert.Error,
		},
		{
			name: "cache refreshed with nil rates map from gateway, no error",
			fields: fields{
				TimeNow: func() time.Time {
					return time.Unix(100, 0)
				},
				RatesCache: map[string]entity.ExchangeRates{},
			},
			args: args{
				req: &entity.GetExchangeRateRequest{
					Country:          "russia",
					BaseCurrencyID:   "GBP",
					TargetCurrencyID: "USD",
					Amount:           10,
				},
			},
			mockRussiaCBGateway: &mockCBGateway{
				res: &entity.ExchangeRates{
					Rates: nil,
				},
				err: nil,
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockRussiaCB := russiagatewaymock.NewMockGateway(ctrl)
			if tt.mockRussiaCBGateway != nil {
				mockRussiaCB.
					EXPECT().
					GetCBRRates(ctx).
					Return(tt.mockRussiaCBGateway.res, tt.mockRussiaCBGateway.err)
			}
			mockThailandCB := thailandgatewaymock.NewMockGateway(ctrl)
			if tt.mockThailandCBGateway != nil {
				mockThailandCB.
					EXPECT().
					GetCBRRates(ctx).
					Return(tt.mockThailandCBGateway.res, tt.mockThailandCBGateway.err)
			}
			c := &cbr{
				RWMutex: sync.RWMutex{},
				TimeNow: tt.fields.TimeNow,
				Gateways: map[string]gateway.CBGateway{
					entity.Russia:   mockRussiaCB,
					entity.Thailand: mockThailandCB,
				},
				RatesCache: tt.fields.RatesCache,
			}
			got, err := c.GetExchangeRate(ctx, tt.args.req)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

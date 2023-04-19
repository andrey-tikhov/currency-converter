package cb_repository

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"my_go/entity"
	"my_go/mapper"
	repositorymock "my_go/mocks/repository"
	"testing"
	"time"
)

func Test_controller_GetCBRates(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	anyRequest := &entity.GetExchangeRatesRequest{
		Country: "russia",
	}
	type args struct {
		r *entity.GetExchangeRatesRequest
	}
	type mockRepository struct {
		res *entity.GetCBRatesResponse
		err error
	}
	ruTZ, _ := time.LoadLocation("Europe/Moscow")
	tests := []struct {
		name           string
		args           args
		mockRepository *mockRepository
		want           *entity.GetExchangeRatesResponse
		assertion      assert.ErrorAssertionFunc
	}{
		{
			name: "Happy path",
			args: args{
				r: anyRequest,
			},
			mockRepository: &mockRepository{
				res: &entity.GetCBRatesResponse{
					Rates: &entity.ExchangeRates{
						Country:    "russia",
						DateLoaded: "2023-01-01",
						TimeZone:   ruTZ,
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
				},
				err: nil,
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
			name: "nil request",
			args: args{
				r: nil,
			},
			want:      nil,
			assertion: assert.Error,
		},
		{
			name: "controller fails",
			args: args{
				r: anyRequest,
			},
			mockRepository: &mockRepository{
				res: nil,
				err: errors.New("some error"),
			},
			want:      nil,
			assertion: assert.Error,
		},
		{
			name: "controller returns nil and no error",
			args: args{
				r: anyRequest,
			},
			mockRepository: &mockRepository{
				res: nil,
				err: nil,
			},
			want:      nil,
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockrepository := repositorymock.NewMockCBR(ctrl)
			if tt.mockRepository != nil {
				if req, err := mapper.GetExchangeRatesRequestToGetCBRatesRequest(tt.args.r); err == nil {
					mockrepository.
						EXPECT().
						GetCBRates(ctx, req).
						Return(tt.mockRepository.res, tt.mockRepository.err)
				}
			}
			c := &controller{
				repository: mockrepository,
			}
			got, err := c.GetCBRates(ctx, tt.args.r)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockrepository := repositorymock.NewMockCBR(ctrl)
	got, err := New(Params{
		Repository: mockrepository,
	})
	assert.NoError(t, err)
	assert.NotNil(t, got)
}

func Test_controller_GetExchangeRate(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		req *entity.GetExchangeRateRequest
	}
	type mockRepository struct {
		res *entity.GetExchangeRateResponse
		err error
	}
	tests := []struct {
		name           string
		args           args
		mockRepository *mockRepository
		want           *entity.GetExchangeRateResponse
		assertion      assert.ErrorAssertionFunc
	}{
		{
			name: "Happy path",
			args: args{
				req: &entity.GetExchangeRateRequest{
					Country:          "russia",
					BaseCurrencyID:   "GBP",
					TargetCurrencyID: "USD",
					Amount:           10,
				},
			},
			mockRepository: &mockRepository{
				res: &entity.GetExchangeRateResponse{
					Rate: entity.Rate{
						Nominal:          10,
						BaseCurrency:     "GBP",
						TargetCurrency:   "USD",
						RateTargetToBase: 12.402,
					},
				},
				err: nil,
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
			name: "repository fails",
			args: args{
				req: &entity.GetExchangeRateRequest{
					Country:          "russia",
					BaseCurrencyID:   "GBP",
					TargetCurrencyID: "USD",
					Amount:           10,
				},
			},
			mockRepository: &mockRepository{
				res: nil,
				err: errors.New("some error"),
			},
			want:      nil,
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockrepository := repositorymock.NewMockCBR(ctrl)
			if tt.mockRepository != nil {
				mockrepository.
					EXPECT().
					GetExchangeRate(ctx, tt.args.req).
					Return(tt.mockRepository.res, tt.mockRepository.err)
			}
			c := &controller{
				repository: mockrepository,
			}
			got, err := c.GetExchangeRate(ctx, tt.args.req)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

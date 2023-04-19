package conversion

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/config"
	internalconfig "my_go/config"
	"my_go/entity"
	"my_go/mapper"
	controllermock "my_go/mocks/controller/cb_repository"
	"my_go/utils"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	src := config.Source(
		strings.NewReader(`{}`),
	)
	providerGood, _ := config.NewYAML(src)
	mock := controllermock.NewMockController(ctrl)
	got, err := New(Params{
		Config:               providerGood,
		RepositoryController: mock,
	})
	assert.NoError(t, err)
	assert.NotNil(t, got)
}

func Test_controller_Convert(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

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

	type fields struct {
		config internalconfig.Defaults
	}
	type args struct {
		req *entity.ConvertCurrencyRequest
	}
	type mockRepositoryController struct {
		res *entity.GetExchangeRateResponse
		err error
	}
	tests := []struct {
		name                     string
		fields                   fields
		args                     args
		mockRepositoryController *mockRepositoryController
		want                     *entity.ConvertCurrencyResponse
		assertion                assert.ErrorAssertionFunc
	}{
		{
			name: "Happy path, country provided",
			fields: fields{
				config: internalconfig.Defaults{
					DefaultCB: "russia",
				},
			},
			args: args{
				req: requestWithCountry,
			},
			mockRepositoryController: &mockRepositoryController{
				res: &entity.GetExchangeRateResponse{
					Rate: entity.Rate{
						Nominal:          12,
						BaseCurrency:     "JPY",
						TargetCurrency:   "USD",
						RateTargetToBase: 1.2345,
					},
				},
				err: nil,
			},
			want: &entity.ConvertCurrencyResponse{
				Amount: 1.2345,
			},
			assertion: assert.NoError,
		},
		{
			name: "Happy path, country not provided",
			fields: fields{
				config: internalconfig.Defaults{
					DefaultCB: "russia",
				},
			},
			args: args{
				req: requestWithoutCountry,
			},
			mockRepositoryController: &mockRepositoryController{
				res: &entity.GetExchangeRateResponse{
					Rate: entity.Rate{
						Nominal:          12,
						BaseCurrency:     "JPY",
						TargetCurrency:   "USD",
						RateTargetToBase: 1.2345,
					},
				},
				err: nil,
			},
			want: &entity.ConvertCurrencyResponse{
				Amount: 1.2345,
			},
			assertion: assert.NoError,
		},
		{
			name: "nil request",
			args: args{
				req: nil,
			},
			want:      nil,
			assertion: assert.Error,
		},
		{
			name: "controller fails",
			fields: fields{
				config: internalconfig.Defaults{
					DefaultCB: "russia",
				},
			},
			args: args{
				req: requestWithCountry,
			},
			mockRepositoryController: &mockRepositoryController{
				res: nil,
				err: errors.New("some error"),
			},
			want:      nil,
			assertion: assert.Error,
		},
		{
			name: "controller returns nil and no error",
			fields: fields{
				config: internalconfig.Defaults{
					DefaultCB: "russia",
				},
			},
			args: args{
				req: requestWithCountry,
			},
			mockRepositoryController: &mockRepositoryController{
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
			mockRepositoryCtrl := controllermock.NewMockController(ctrl)
			if req, err := mapper.ConvertCurrencyRequestToGetExchangeRateRequest(
				tt.args.req, tt.fields.config.DefaultCB,
			); err == nil {
				mockRepositoryCtrl.
					EXPECT().
					GetExchangeRate(ctx, req).
					Return(tt.mockRepositoryController.res, tt.mockRepositoryController.err)
			}
			c := &controller{
				config:               tt.fields.config,
				repositoryController: mockRepositoryCtrl,
			}
			got, err := c.Convert(ctx, tt.args.req)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

package handler

import (
	"bytes"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"my_go/entity"
	"my_go/mapper"
	cb_repositorymock "my_go/mocks/controller/cb_repository"
	conversionmock "my_go/mocks/controller/conversion"
	"net/http"
	"net/http/httptest"
	"testing"
)

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

func TestNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repositoryCtrlMock := cb_repositorymock.NewMockController(ctrl)
	conversionCtrlMock := conversionmock.NewMockController(ctrl)
	got, err := New(Params{
		CBRepositoryController: repositoryCtrlMock,
		ConversionController:   conversionCtrlMock,
	})
	assert.NoError(t, err)
	assert.NotNil(t, got)
}

func Test_handler_GetCBRates(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type mockCBRepositoryController struct {
		res *entity.GetExchangeRatesResponse
		err error
	}
	type args struct {
		method string
		body   []byte
		url    string
	}
	tests := []struct {
		name                       string
		args                       args
		failureBody                bool
		mockCBRepositoryController *mockCBRepositoryController
		expectedStatusCode         int
		expectedResponse           string
	}{
		{
			name: "Happy path",
			args: args{
				method: "POST",
				body:   []byte(`{"country":"thailand"}`),
				url:    "/get_exchange_rates",
			},
			mockCBRepositoryController: &mockCBRepositoryController{
				res: &entity.GetExchangeRatesResponse{
					Rates: map[string]entity.Rate{
						"USD": {
							Nominal:          100,
							BaseCurrency:     "RUR",
							TargetCurrency:   "USD",
							RateTargetToBase: 123.56,
						},
					},
				},
				err: nil,
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"rates":{"USD":{"nominal":100,"base_currency":"RUR","target_currency":"USD","rate_target_to_base":123.56}}}`,
		},
		{
			name: "wrong request method",
			args: args{
				method: "GET",
				body:   []byte(`{"country":"thailand"}`),
				url:    "/get_exchange_rates",
			},
			expectedStatusCode: http.StatusMethodNotAllowed,
			expectedResponse:   "method not allowed\n",
		},
		{
			name: "failed to read body",
			args: args{
				method: "POST",
				url:    "/get_exchange_rates",
			},
			failureBody:        true,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   "unable to read the body\n",
		},
		{
			name: "failed to convert body to the internal entity",
			args: args{
				method: "POST",
				body:   []byte(`{"srgsrgs_`),
				url:    "/get_exchange_rates",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   "bad request, err failed to unmarshal: unexpected end of JSON input\n",
		},
		{
			name: "controller fails",
			args: args{
				method: "POST",
				body:   []byte(`{"country":"thailand"}`),
				url:    "/get_exchange_rates",
			},
			mockCBRepositoryController: &mockCBRepositoryController{
				res: nil,
				err: errors.New("some error"),
			},
			expectedStatusCode: http.StatusBadGateway,
			expectedResponse:   "failed to process the request, err some error\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpreq, _ := http.NewRequest(tt.args.method, tt.args.url, bytes.NewReader(tt.args.body))
			if tt.failureBody {
				httpreq, _ = http.NewRequest(tt.args.method, tt.args.url, errReader(0))
			}
			repositoryCtrlMock := cb_repositorymock.NewMockController(ctrl)
			if req, err := mapper.BodyToGetExchangeRatesRequest(tt.args.body); err == nil &&
				tt.mockCBRepositoryController != nil {
				repositoryCtrlMock.
					EXPECT().
					GetCBRates(httpreq.Context(), req).
					Return(tt.mockCBRepositoryController.res, tt.mockCBRepositoryController.err)
			}
			conversionCtrlMock := conversionmock.NewMockController(ctrl)
			h := &handler{
				repositoryCtrl: repositoryCtrlMock,
				conversionCtrl: conversionCtrlMock,
			}
			rr := httptest.NewRecorder()
			testhandler := http.HandlerFunc(h.GetCBRates)
			testhandler.ServeHTTP(rr, httpreq)
			assert.Equal(t, tt.expectedStatusCode, rr.Code)
			assert.Equal(t, tt.expectedResponse, rr.Body.String())
		})
	}
}

func Test_handler_ConvertCurrency(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type mockConversionController struct {
		res *entity.ConvertCurrencyResponse
		err error
	}
	type args struct {
		method string
		body   []byte
		url    string
	}
	tests := []struct {
		name                     string
		args                     args
		failureBody              bool
		mockConversionController *mockConversionController
		expectedStatusCode       int
		expectedResponse         string
	}{
		{
			name: "Happy path",
			args: args{
				method: "POST",
				body:   []byte(`{"country":"russia","source_currency":"RUR","target_currency":"USD","amount":100}`),
				url:    "/convert",
			},
			mockConversionController: &mockConversionController{
				res: &entity.ConvertCurrencyResponse{
					Amount: 12.2345,
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"amount":12.2345}`,
		},
		{
			name: "wrong request method",
			args: args{
				method: "GET",
				body:   []byte(`{"country":"russia","source_currency":"RUR","target_currency":"USD","amount":100}`),
				url:    "/convert",
			},
			expectedStatusCode: http.StatusMethodNotAllowed,
			expectedResponse:   "method not allowed\n",
		},
		{
			name: "failed to read body",
			args: args{
				method: "POST",
				url:    "/convert",
			},
			failureBody:        true,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   "unable to read the body\n",
		},
		{
			name: "failed to convert body to the internal entity",
			args: args{
				method: "POST",
				body:   []byte(`{"srgsrgs_`),
				url:    "/convert",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   "bad request, err failed to unmarshal: unexpected end of JSON input\n",
		},
		{
			name: "controller fails",
			args: args{
				method: "POST",
				body:   []byte(`{"country":"russia","source_currency":"RUR","target_currency":"USD","amount":100}`),
				url:    "/convert",
			},
			mockConversionController: &mockConversionController{
				res: nil,
				err: errors.New("some error"),
			},
			expectedStatusCode: http.StatusBadGateway,
			expectedResponse:   "failed to process the request, err some error\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpreq, _ := http.NewRequest(tt.args.method, tt.args.url, bytes.NewReader(tt.args.body))
			if tt.failureBody {
				httpreq, _ = http.NewRequest(tt.args.method, tt.args.url, errReader(0))
			}

			repositoryCtrlMock := cb_repositorymock.NewMockController(ctrl)
			conversionCtrlMock := conversionmock.NewMockController(ctrl)
			if req, err := mapper.BodyToConvertCurrencyRequest(tt.args.body); err == nil &&
				tt.mockConversionController != nil {
				conversionCtrlMock.
					EXPECT().
					Convert(httpreq.Context(), req).
					Return(tt.mockConversionController.res, tt.mockConversionController.err)
			}
			h := &handler{
				repositoryCtrl: repositoryCtrlMock,
				conversionCtrl: conversionCtrlMock,
			}
			rr := httptest.NewRecorder()
			testhandler := http.HandlerFunc(h.ConvertCurrency)
			testhandler.ServeHTTP(rr, httpreq)
			assert.Equal(t, tt.expectedStatusCode, rr.Code)
			assert.Equal(t, tt.expectedResponse, rr.Body.String())
		})
	}
}

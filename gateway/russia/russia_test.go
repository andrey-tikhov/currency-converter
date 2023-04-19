package russia

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/config"
	internalconfig "my_go/config"
	"my_go/entity"
	"my_go/utils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	src := config.Source(
		strings.NewReader(`{}`),
	)
	providerGood, _ := config.NewYAML(src)
	type args struct {
		c config.Provider
	}
	tests := []struct {
		name      string
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Happy path",
			args: args{
				c: providerGood,
			},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.c)
			tt.assertion(t, err)
			assert.NotNil(t, got)
		})
	}
}

func Test_russiaCRBGateway_GetCBRRates(t *testing.T) {
	var nilContext context.Context
	timeNow := func() time.Time {
		return time.Unix(100, 0)
	}
	ruTZ, _ := time.LoadLocation("Europe/Moscow")
	correctXML := []byte(`<ValCurs Date="18.04.2023" name="Foreign Currency Market">
  <Valute ID="R01010">
    <NumCode>036</NumCode>
    <CharCode>AUD</CharCode>
    <Nominal>1</Nominal>
    <Name>РђРІСЃС‚СЂР°Р»РёР№СЃРєРёР№ РґРѕР»Р»Р°СЂ</Name>
    <Value>54,8131</Value>
  </Valute>
  <Valute ID="R01020A">
    <NumCode>944</NumCode>
    <CharCode>AZN</CharCode>
    <Nominal>10</Nominal>
    <Name>РђР·РµСЂР±Р°Р№РґР¶Р°РЅСЃРєРёР№ РјР°РЅР°С‚</Name>
    <Value>48,0164</Value>
  </Valute></ValCurs>
`)
	type fields struct {
		TimeNow  func() time.Time
		TimeZone string
	}
	tests := []struct {
		name                     string
		ContentLenghth           *string
		httpRequestCreationFails bool
		isTimedOut               bool
		httpRespStatusCode       int
		httpRespBody             []byte
		fields                   fields
		want                     *entity.ExchangeRates
		assertion                assert.ErrorAssertionFunc
	}{
		{
			name:               "Happy path",
			httpRespStatusCode: 200,
			httpRespBody:       correctXML,
			fields: fields{
				TimeNow:  timeNow,
				TimeZone: "Europe/Moscow",
			},
			want: &entity.ExchangeRates{
				Country:    "russia",
				DateLoaded: timeNow().Format("2006-01-02"),
				TimeZone:   ruTZ,
				Rates: map[string]entity.Rate{
					"AUD": {
						Nominal:          1,
						BaseCurrency:     "RUR",
						TargetCurrency:   "AUD",
						RateTargetToBase: 54.8131,
					},
					"AZN": {
						Nominal:          10,
						BaseCurrency:     "RUR",
						TargetCurrency:   "AZN",
						RateTargetToBase: 48.0164,
					},
					"RUR": {
						Nominal:          1,
						BaseCurrency:     "RUR",
						TargetCurrency:   "RUR",
						RateTargetToBase: 1,
					},
				},
			},
			assertion: assert.NoError,
		},
		{
			name:       "network timeout",
			isTimedOut: true,
			want:       nil,
			assertion:  assert.Error,
		},
		{
			name:               "unexpected status code from server",
			httpRespStatusCode: 404,
			want:               nil,
			assertion:          assert.Error,
		},
		{
			name:               "incorrect xml arrived from server",
			httpRespStatusCode: 200,
			httpRespBody:       []byte(`<this is for sure not xml`),
			want:               nil,
			assertion:          assert.Error,
		},
		{
			name:               "failed to read body",
			ContentLenghth:     utils.ToPointer("1"),
			httpRespStatusCode: 200,
			httpRespBody:       correctXML,
			want:               nil,
			assertion:          assert.Error,
		},
		{
			name:                     "http request creation fails",
			httpRequestCreationFails: true,
			want:                     nil,
			assertion:                assert.Error,
		},
		{
			name:               "bad timezone in gateway config",
			httpRespStatusCode: 200,
			httpRespBody:       correctXML,
			fields: fields{
				TimeNow:  timeNow,
				TimeZone: "1234",
			},
			want:      nil,
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			release := make(chan interface{})
			timedOutMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				<-release
			}))
			defer timedOutMock.Close()
			defer close(release)

			TestServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.ContentLenghth != nil {
					w.Header().Set("Content-Length", *tt.ContentLenghth)
				}
				w.WriteHeader(tt.httpRespStatusCode)
				w.Write(tt.httpRespBody)
			}))
			defer TestServer.Close()

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
			defer cancel()
			if tt.httpRequestCreationFails {
				ctx = nilContext
			}

			url := TestServer.URL
			if tt.isTimedOut {
				url = timedOutMock.URL
			}
			g := &russiaCRBGateway{
				TimeNow: tt.fields.TimeNow,
				Config: internalconfig.RussiaCBConfig{
					APIURL:   url,
					Timezone: tt.fields.TimeZone,
				},
			}
			got, err := g.GetCBRRates(ctx)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

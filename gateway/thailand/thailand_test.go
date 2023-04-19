package thailand

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
		want      Gateway
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

func Test_thailandCRBGateway_GetCBRRates(t *testing.T) {
	var nilContext context.Context
	timeNow := func() time.Time {
		return time.Unix(100, 0)
	}
	thTZ, _ := time.LoadLocation("Asia/Bangkok")
	correctXML := []byte(`
<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#" xmlns="http://purl.org/rss/1.0/" xmlns:cb="http://centralbanks.org/cb/1.0/" xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:dcterms="http://purl.org/dc/terms/" xmlns:xsi="http://www.w3c.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.w3c.org/1999/02/22-rdf-syntax-ns#rdf.xsd">
<channel rdf:about="http://www.bot.or.th/english/statistics/financialmarkets/exchangerate/_layouts/Application/ExchangeRate/ExchangeRate.aspx">
<title>Bank of Thailand: Daily Foreign Exchange Rates</title>
<link>http://www.bot.or.th/english/statistics/financialmarkets/exchangerate/_layouts/Application/ExchangeRate/ExchangeRate.aspx</link>
<description>Daily foreign exchange rates rates from Bank of Thailand.</description>
<items>
<rdf:Seq>
<rdf:li rdf:resource="https://www.bot.or.th/App/RSS/fxrate-USD.xml"/>
<rdf:li rdf:resource="https://www.bot.or.th/App/RSS/fxrate-GBP.xml"/>
</rdf:Seq>
</items>
<dc:language>en</dc:language>
<dc:date>2023-04-17</dc:date>
</channel>
<item rdf:about="https://www.bot.or.th/App/RSS/fxrate-USD.xml">
<title>TH: 34.0625 THB = 1 USD 2023-04-17 Bank of Thailand Average Buying Sight Bill</title>
<link>http://www.bot.or.th/english/statistics/financialmarkets/exchangerate/_layouts/Application/ExchangeRate/ExchangeRate.aspx</link>
<description>34.0625 Thai Baht = 1 USD</description>
<dc:language>en</dc:language>
<dc:date>2023-04-17</dc:date>
<dc:format>text/html</dc:format>
<cb:country>TH</cb:country>
<cb:baseCurrency>THB</cb:baseCurrency>
<cb:targetCurrency>USD</cb:targetCurrency>
<cb:value frequency="business" decimals="4">34.0625</cb:value>
<cb:rateType>Daily</cb:rateType>
<cb:application>statistics</cb:application>
</item>
<item rdf:about="https://www.bot.or.th/App/RSS/fxrate-USD.xml">
<title>TH: 34.1656 THB = 1 USD 2023-04-17 Bank of Thailand Average Buying Transfer</title>
<link>http://www.bot.or.th/english/statistics/financialmarkets/exchangerate/_layouts/Application/ExchangeRate/ExchangeRate.aspx</link>
<description>34.1656 Thai Baht = 1 USD</description>
<dc:language>en</dc:language>
<dc:date>2023-04-17</dc:date>
<dc:format>text/html</dc:format>
<cb:country>TH</cb:country>
<cb:baseCurrency>THB</cb:baseCurrency>
<cb:targetCurrency>USD</cb:targetCurrency>
<cb:value frequency="business" decimals="4">34.1656</cb:value>
<cb:rateType>Daily</cb:rateType>
<cb:application>statistics</cb:application>
</item>
<item rdf:about="https://www.bot.or.th/App/RSS/fxrate-USD.xml">
<title>TH: 34.5022 THB = 1 USD 2023-04-17 Bank of Thailand Average Selling Rate</title>
<link>http://www.bot.or.th/english/statistics/financialmarkets/exchangerate/_layouts/Application/ExchangeRate/ExchangeRate.aspx</link>
<description>34.5022 Thai Baht = 1 USD</description>
<dc:language>en</dc:language>
<dc:date>2023-04-17</dc:date>
<dc:format>text/html</dc:format>
<cb:country>TH</cb:country>
<cb:baseCurrency>THB</cb:baseCurrency>
<cb:targetCurrency>USD</cb:targetCurrency>
<cb:value frequency="business" decimals="4">34.5022</cb:value>
<cb:rateType>Daily</cb:rateType>
<cb:application>statistics</cb:application>
</item>
<item rdf:about="https://www.bot.or.th/App/RSS/fxrate-GBP.xml">
<title>TH: 42.0040 THB = 1 GBP 2023-04-17 Bank of Thailand Average Buying Sight Bill</title>
<link>http://www.bot.or.th/english/statistics/financialmarkets/exchangerate/_layouts/Application/ExchangeRate/ExchangeRate.aspx</link>
<description>42.0040 Thai Baht = 1 GBP</description>
<dc:language>en</dc:language>
<dc:date>2023-04-17</dc:date>
<dc:format>text/html</dc:format>
<cb:country>TH</cb:country>
<cb:baseCurrency>THB</cb:baseCurrency>
<cb:targetCurrency>GBP</cb:targetCurrency>
<cb:value frequency="business" decimals="4">42.0040</cb:value>
<cb:rateType>Daily</cb:rateType>
<cb:application>statistics</cb:application>
</item>
<item rdf:about="https://www.bot.or.th/App/RSS/fxrate-GBP.xml">
<title>TH: 42.1612 THB = 1 GBP 2023-04-17 Bank of Thailand Average Buying Transfer</title>
<link>http://www.bot.or.th/english/statistics/financialmarkets/exchangerate/_layouts/Application/ExchangeRate/ExchangeRate.aspx</link>
<description>42.1612 Thai Baht = 1 GBP</description>
<dc:language>en</dc:language>
<dc:date>2023-04-17</dc:date>
<dc:format>text/html</dc:format>
<cb:country>TH</cb:country>
<cb:baseCurrency>THB</cb:baseCurrency>
<cb:targetCurrency>GBP</cb:targetCurrency>
<cb:value frequency="business" decimals="4">42.1612</cb:value>
<cb:rateType>Daily</cb:rateType>
<cb:application>statistics</cb:application>
</item>
<item rdf:about="https://www.bot.or.th/App/RSS/fxrate-GBP.xml">
<title>TH: 43.0302 THB = 1 GBP 2023-04-17 Bank of Thailand Average Selling Rate</title>
<link>http://www.bot.or.th/english/statistics/financialmarkets/exchangerate/_layouts/Application/ExchangeRate/ExchangeRate.aspx</link>
<description>43.0302 Thai Baht = 1 GBP</description>
<dc:language>en</dc:language>
<dc:date>2023-04-17</dc:date>
<dc:format>text/html</dc:format>
<cb:country>TH</cb:country>
<cb:baseCurrency>THB</cb:baseCurrency>
<cb:targetCurrency>GBP</cb:targetCurrency>
<cb:value frequency="business" decimals="4">43.0302</cb:value>
<cb:rateType>Daily</cb:rateType>
<cb:application>statistics</cb:application>
</item>
</rdf:RDF>
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
				TimeZone: "Asia/Bangkok",
			},
			want: &entity.ExchangeRates{
				Country:    "thailand",
				DateLoaded: timeNow().Format("2006-01-02"),
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
			g := &thailandCRBGateway{
				TimeNow: tt.fields.TimeNow,
				Config: internalconfig.ThailandCBConfig{
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

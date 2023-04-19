package cb

import (
	"github.com/stretchr/testify/assert"
	"my_go/entity"
	"testing"
)

func TestThailandCBRResponseToRates(t *testing.T) {
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
	correctXMLWithBadDescription := []byte(`
<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#" xmlns="http://purl.org/rss/1.0/" xmlns:cb="http://centralbanks.org/cb/1.0/" xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:dcterms="http://purl.org/dc/terms/" xmlns:xsi="http://www.w3c.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.w3c.org/1999/02/22-rdf-syntax-ns#rdf.xsd">
<item rdf:about="https://www.bot.or.th/App/RSS/fxrate-USD.xml">
<title>TH: 34.0625 THB = 1 USD 2023-04-17 Bank of Thailand Average Buying Sight Bill</title>
<link>http://www.bot.or.th/english/statistics/financialmarkets/exchangerate/_layouts/Application/ExchangeRate/ExchangeRate.aspx</link>
<description>some unexpected description</description>
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
</rdf:RDF>
`)
	type args struct {
		body []byte
	}
	tests := []struct {
		name      string
		args      args
		want      map[string]entity.Rate
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Happy path",
			args: args{
				body: correctXML,
			},
			want: map[string]entity.Rate{
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
			assertion: assert.NoError,
		},
		{
			name: "bad xml received",
			args: args{
				body: []byte(`<asef3_355`),
			},
			want:      nil,
			assertion: assert.Error,
		},
		{
			name: "description contains bad string",
			args: args{
				body: correctXMLWithBadDescription,
			},
			want: map[string]entity.Rate{
				"THB": {
					Nominal:          1,
					RateTargetToBase: 1,
					BaseCurrency:     "THB",
					TargetCurrency:   "THB",
				},
				"USD": {
					Nominal:          1,
					RateTargetToBase: 34.5022,
					BaseCurrency:     "THB",
					TargetCurrency:   "USD",
				},
			},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ThailandCBRResponseToRates(tt.args.body)
			assert.Equal(t, tt.want, got)
			tt.assertion(t, err)
		})
	}
}

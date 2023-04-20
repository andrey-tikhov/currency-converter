package cb

import (
	"github.com/stretchr/testify/assert"
	"my_go/entity"
	"testing"
)

func TestRussiaCBRResponseToRates(t *testing.T) {
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
	correctXMLwithBadValue := []byte(`<ValCurs Date="18.04.2023" name="Foreign Currency Market">
  <Valute ID="R01010">
    <NumCode>036</NumCode>
    <CharCode>AUD</CharCode>
    <Nominal>1</Nominal>
    <Name>РђРІСЃС‚СЂР°Р»РёР№СЃРєРёР№ РґРѕР»Р»Р°СЂ</Name>
    <Value>some weird data</Value>
  </Valute>
  <Valute ID="R01020A">
    <NumCode>944</NumCode>
    <CharCode>AZN</CharCode>
    <Nominal>10</Nominal>
    <Name>РђР·РµСЂР±Р°Р№РґР¶Р°РЅСЃРєРёР№ РјР°РЅР°С‚</Name>
    <Value>48,0164</Value>
  </Valute></ValCurs>
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
				"AUD": {
					Nominal:          1,
					BaseCurrency:     "RUB",
					TargetCurrency:   "AUD",
					RateTargetToBase: 54.8131,
				},
				"AZN": {
					Nominal:          10,
					BaseCurrency:     "RUB",
					TargetCurrency:   "AZN",
					RateTargetToBase: 48.0164,
				},
				"RUB": {
					Nominal:          1,
					BaseCurrency:     "RUB",
					TargetCurrency:   "RUB",
					RateTargetToBase: 1,
				},
			},
			assertion: assert.NoError,
		},
		{
			name: "bad xml",
			args: args{
				body: []byte(`<_q352462**)$5`),
			},
			want:      nil,
			assertion: assert.Error,
		},
		{
			name: "response from central bank contains bad exchange rate value",
			args: args{
				body: correctXMLwithBadValue,
			},
			want: map[string]entity.Rate{
				"AZN": {
					Nominal:          10,
					BaseCurrency:     "RUB",
					TargetCurrency:   "AZN",
					RateTargetToBase: 48.0164,
				},
				"RUB": {
					Nominal:          1,
					BaseCurrency:     "RUB",
					TargetCurrency:   "RUB",
					RateTargetToBase: 1,
				},
			},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RussiaCBRResponseToRates(tt.args.body)
			assert.Equal(t, tt.want, got)
			tt.assertion(t, err)
		})
	}
}

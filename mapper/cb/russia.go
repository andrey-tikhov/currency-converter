package cb

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"golang.org/x/net/html/charset"
	"my_go/entity"
	cb_entity "my_go/entity/cb"
	"strconv"
	"strings"
)

// RussiaCBRResponseToRates converts the response from Central Bank of Russia to
// map where keys are currency ID and values are entity.Rate.
// For the convenience of the conversion calculation the rate of RUR to RUR conversion is added.
func RussiaCBRResponseToRates(body []byte) (map[string]entity.Rate, error) {
	reader := bytes.NewReader(body)
	parser := xml.NewDecoder(reader)
	parser.CharsetReader = charset.NewReaderLabel
	resp := cb_entity.RussiaCBRData{}
	err := parser.Decode(&resp)
	if err != nil {
		return nil, fmt.Errorf("xml unmarshal failed %s", err)
	}
	m := map[string]entity.Rate{}
	for _, r := range resp.Rates {
		ratio, err := strconv.ParseFloat(strings.Replace(r.Value, ",", ".", 1), 64)
		if err != nil {
			// TODO add error log here
			continue
		}
		m[r.CurrencyID] = entity.Rate{
			Nominal:          r.Nominal,
			BaseCurrency:     "RUR",
			TargetCurrency:   r.CurrencyID,
			RateTargetToBase: ratio,
		}
	}
	m["RUR"] = entity.Rate{
		Nominal:          1,
		BaseCurrency:     "RUR",
		TargetCurrency:   "RUR",
		RateTargetToBase: 1,
	}
	return m, nil
}

package cb

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"golang.org/x/net/html/charset"
	"my_go/entity"
	cb_entity "my_go/entity/cb"
	"regexp"
	"strconv"
	"strings"
)

const (
	buyRateString  = "Bank of Thailand Average Buying Sight Bill"
	sellRateString = "Bank of Thailand Average Selling Rate"
)

var descriptionRegEx = regexp.MustCompile("[0-9\\.]+ Thai Baht = ([0-9]+) [A-Z]{3}")

// ThailandCBRResponseToRates converts xml response from Thai central bank to the map of
// currency ids to map where keys are currency ID and values are entity.Rate.
// The foreign exchange rate for Thai bank is taken as average between Buying Sight bill and Selling Rate as no clear
// definition (as for example bank of Russia has) for a daily rate is provided
// For the convenience of the conversion calculation the rate of RUR to RUR conversion is added.
func ThailandCBRResponseToRates(body []byte) (map[string]entity.Rate, error) {
	reader := bytes.NewReader(body)
	parser := xml.NewDecoder(reader)
	parser.CharsetReader = charset.NewReaderLabel
	resp := cb_entity.ThailandCBRData{}
	err := parser.Decode(&resp)
	if err != nil {
		return nil, fmt.Errorf("xml unmarshal failed %s", err)
	}
	m := map[string]entity.Rate{}
	for _, r := range resp.Rates {
		if !strings.Contains(r.Title, buyRateString) && !strings.Contains(r.Title, sellRateString) {
			continue
		}
		v := descriptionRegEx.FindStringSubmatch(r.Description)
		if len(v) < 2 {
			continue
		}
		nominal, err := strconv.Atoi(v[1])
		if err != nil {
			continue // unreachable in tests cause this group contains from 0..9 only
		}
		data, ok := m[r.TargetCurrency]
		var rateTargetToBase float64
		if !ok {
			rateTargetToBase = r.Value
		} else {
			// might be issues if rates are with 6 digits. but no problems for 4 digits at all.
			rateTargetToBase = (r.Value + data.RateTargetToBase) / 2
		}
		m[r.TargetCurrency] = entity.Rate{
			Nominal:          nominal,
			RateTargetToBase: rateTargetToBase,
			TargetCurrency:   r.TargetCurrency,
			BaseCurrency:     "THB",
		}
	}
	m["THB"] = entity.Rate{
		Nominal:          1,
		BaseCurrency:     "THB",
		TargetCurrency:   "THB",
		RateTargetToBase: 1,
	}
	return m, nil
}

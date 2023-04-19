package cb

import "encoding/xml"

type RussiaCBRData struct {
	XMLName xml.Name       `xml:"ValCurs"`
	Rates   []RussiaCBRate `xml:"Valute"`
}

type RussiaCBRate struct {
	CurrencyID string `xml:"CharCode"`
	Nominal    int    `xml:"Nominal"`
	Value      string `xml:"Value"`
}

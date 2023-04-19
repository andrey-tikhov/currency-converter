package cb

import "encoding/xml"

type ThailandCBRData struct {
	XMLName xml.Name         `xml:"RDF"`
	Rates   []ThailandCBRate `xml:"item"`
}

type ThailandCBRate struct {
	Title          string  `xml:"title"`
	Description    string  `xml:"description"`
	TargetCurrency string  `xml:"targetCurrency"`
	Value          float64 `xml:"value"`
}

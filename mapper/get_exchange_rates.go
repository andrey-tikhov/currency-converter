package mapper

import (
	"encoding/json"
	"fmt"
	"my_go/entity"
)

func BodyToGetExchangeRatesRequest(body []byte) (*entity.GetExchangeRatesRequest, error) {
	var r entity.GetExchangeRatesRequest
	err := json.Unmarshal(body, &r)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %s", err)
	}
	return &r, nil
}

func GetExchangeRatesResponseToBytes(r *entity.GetExchangeRatesResponse) ([]byte, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal: %s", err)
	}
	return b, nil
}

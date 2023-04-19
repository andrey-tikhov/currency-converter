package gateway

import (
	"context"
	"my_go/entity"
)

type CBGateway interface {
	GetCBRRates(ctx context.Context) (*entity.ExchangeRates, error)
}

package cb_repository

import (
	"context"
	"go.uber.org/fx"
	"my_go/entity"
	"my_go/mapper"
	"my_go/repository"
)

type Controller interface {
	GetCBRates(ctx context.Context, req *entity.GetExchangeRatesRequest) (*entity.GetExchangeRatesResponse, error)
	GetExchangeRate(ctx context.Context, req *entity.GetExchangeRateRequest) (*entity.GetExchangeRateResponse, error)
}

var _ Controller = (*controller)(nil)

type controller struct {
	repository repository.CBR
}

type Params struct {
	fx.In

	Repository repository.CBR
}

func New(p Params) (Controller, error) {
	return &controller{
		repository: p.Repository,
	}, nil
}

func (c *controller) GetCBRates(
	ctx context.Context,
	r *entity.GetExchangeRatesRequest,
) (*entity.GetExchangeRatesResponse, error) {
	req, err := mapper.GetExchangeRatesRequestToGetCBRatesRequest(r)
	if err != nil {
		return nil, err
	}
	data, err := c.repository.GetCBRates(ctx, req)
	if err != nil {
		return nil, err
	}
	resp, err := mapper.GetCBRRatesResponseToGetExchangeRatesResponse(data)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *controller) GetExchangeRate(
	ctx context.Context,
	req *entity.GetExchangeRateRequest,
) (*entity.GetExchangeRateResponse, error) {
	return c.repository.GetExchangeRate(ctx, req)
}

package conversion

import (
	"context"
	"go.uber.org/config"
	"go.uber.org/fx"
	internalconfig "my_go/config"
	repositorycontroller "my_go/controller/cb_repository"
	"my_go/entity"
	"my_go/mapper"
)

const defaults = "defaults"

// Controller is a interface to provide currency conversion data
type Controller interface {
	Convert(ctx context.Context, req *entity.ConvertCurrencyRequest) (*entity.ConvertCurrencyResponse, error)
}

// Compile time check that controller implements Controller interface
var _ Controller = (*controller)(nil)

type controller struct {
	config               internalconfig.Defaults
	repositoryController repositorycontroller.Controller
}

// Params is a container for the Controller dependencies
type Params struct {
	fx.In

	Config               config.Provider
	RepositoryController repositorycontroller.Controller
}

// New is a constructor for the Controller interface
func New(p Params) (Controller, error) {
	var d internalconfig.Defaults
	err := p.Config.Get(defaults).Populate(&d)
	if err != nil {
		return nil, err // unreachable in tests, cause provider is populating from valid yaml.
	}
	return &controller{
		config:               d,
		repositoryController: p.RepositoryController,
	}, nil
}

// Convert loads the conversion rate of provided currencies and amount for the country central bank specified.
// If no country is specified it falls back to default one set in the config
func (c *controller) Convert(
	ctx context.Context,
	req *entity.ConvertCurrencyRequest,
) (*entity.ConvertCurrencyResponse, error) {
	r, err := mapper.ConvertCurrencyRequestToGetExchangeRateRequest(req, c.config.DefaultCB)
	if err != nil {
		return nil, err
	}
	got, err := c.repositoryController.GetExchangeRate(ctx, r)
	if err != nil {
		return nil, err
	}
	res, err := mapper.GetExchangeRateResponseToConvertCurrencyResponse(got)
	if err != nil {
		return nil, err
	}
	return res, nil
}

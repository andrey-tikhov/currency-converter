package repository

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/fx"
	"my_go/entity"
	"my_go/gateway"
	"my_go/gateway/russia"
	"my_go/gateway/thailand"
	"my_go/mapper"
	"sync"
	"time"
)

type CBR interface {
	GetCBRates(ctx context.Context, req *entity.GetCBRatesRequest) (*entity.GetCBRatesResponse, error)
	GetExchangeRate(ctx context.Context, req *entity.GetExchangeRateRequest) (*entity.GetExchangeRateResponse, error)
}

// Compile time check that cbr implements CBR interface
var _ CBR = (*cbr)(nil)

// Params is a container for all CBR dependencies
type Params struct {
	fx.In

	ThaiGateway   thailand.Gateway
	RussiaGateway russia.Gateway
}

type cbr struct {
	sync.RWMutex

	TimeNow    func() time.Time
	Gateways   map[string]gateway.CBGateway    // maps country to respective gateway
	RatesCache map[string]entity.ExchangeRates // maps country to ExchangeRatesObject
}

// New is a constructor for the CBR interface
// Cache implementation is leveraging the fact that fx module that provides this constructor
// is called in lazy manner. Meaning once interface was created it will be re-used.
// Hence in memory cache will be kept in a proper state.
func New(p Params) (CBR, error) {
	return &cbr{
		TimeNow: time.Now,
		Gateways: map[string]gateway.CBGateway{
			entity.Russia:   p.RussiaGateway,
			entity.Thailand: p.ThaiGateway,
		},
		RatesCache: map[string]entity.ExchangeRates{
			entity.Russia:   {},
			entity.Thailand: {},
		},
	}, nil
}

// GetCBRates loads central bank rates for the central bank of the country provided in the request
// if country is not implemented it fails (assuming no central bank = no rates).
func (c *cbr) GetCBRates(ctx context.Context, req *entity.GetCBRatesRequest) (*entity.GetCBRatesResponse, error) {
	if req == nil {
		return nil, errors.New("nil GetCBRatesRequest")
	}
	c.RLock()
	cachedRates, ok := c.RatesCache[req.Country]
	c.RUnlock()

	if !ok || c.needsRefresh(&cachedRates) {
		if err := c.reloadCache(ctx, req.Country); err != nil {
			return nil, fmt.Errorf("failed to reload rates: %s", err)
		}
		c.RLock()
		cachedRates, ok = c.RatesCache[req.Country]
		c.RUnlock()
		if !ok {
			// unreachable in tests because if cache was reloaded successfully the key will be present in map
			return nil, fmt.Errorf("provided Country %s unsupported", req.Country)
		}
	}
	return &entity.GetCBRatesResponse{
		Rates: &cachedRates,
	}, nil
}

func (c *cbr) GetExchangeRate(ctx context.Context, req *entity.GetExchangeRateRequest) (*entity.GetExchangeRateResponse, error) {
	if req == nil {
		return nil, errors.New("nil GetExchangeRateRequest")
	}
	rates, err := c.GetCBRates(ctx, &entity.GetCBRatesRequest{
		Country: req.Country,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to load exchange rates for CB %s, err: %s", req.Country, err)
	}
	resp, err := mapper.CBRRatesAndGetExchangeRateRequestToGetExchangeRateResponse(rates.Rates, req)
	if err != nil {
		return nil, fmt.Errorf("failed to convert rates and request to response: %s", err)
	}
	return resp, nil
}

func (c *cbr) reloadCache(ctx context.Context, country string) error {
	c.Lock()
	defer c.Unlock()
	if rates, ok := c.RatesCache[country]; ok {
		if !c.needsRefresh(&rates) {
			return nil
		}
	}
	gw, ok := c.Gateways[country]
	if !ok {
		return fmt.Errorf("provided country %s unsupported", country)
	}
	rates, err := gw.GetCBRRates(ctx)
	if err != nil {
		return fmt.Errorf("failed to load %s central bank data: %s", country, err)
	}
	if rates == nil {
		return fmt.Errorf("nil rates returned from central bank %s with no error", country)
	}
	c.RatesCache[country] = *rates
	return nil
}

func (c *cbr) needsRefresh(r *entity.ExchangeRates) bool {
	if r == nil {
		return true
	}
	if r.DateLoaded == "" {
		return true
	}
	if c.TimeNow().In(r.TimeZone).Format(entity.DateLayout) == r.DateLoaded {
		return false
	}
	return true
}

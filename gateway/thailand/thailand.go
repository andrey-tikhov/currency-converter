package thailand

import (
	"context"
	"fmt"
	"go.uber.org/config"
	"io"
	internalconfig "my_go/config"
	"my_go/entity"
	"my_go/gateway"
	mapper "my_go/mapper/cb"
	"net/http"
	"time"
)

const configKey = "thailand_cb_config"

// Gateway is an interface that implements shared interface of CBGateway for Thailand
type Gateway interface {
	gateway.CBGateway
}

// Compile time check that thailandCRBGateway implements Gateway interface
var _ Gateway = (*thailandCRBGateway)(nil)

type thailandCRBGateway struct {
	TimeNow func() time.Time
	Config  internalconfig.ThailandCBConfig
}

// New is a constructor for Gateway interface
func New(c config.Provider) (Gateway, error) {
	var cfg internalconfig.ThailandCBConfig
	err := c.Get(configKey).Populate(&cfg)
	if err != nil {
		return nil, err // unreachable in tests, cause provider is populating from valid yaml.
	}
	return &thailandCRBGateway{
		TimeNow: time.Now,
		Config:  cfg,
	}, nil
}

// GetCBRRates returns exchange rates for the bank of Thailand
func (g *thailandCRBGateway) GetCBRRates(ctx context.Context) (*entity.ExchangeRates, error) {
	client := http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequestWithContext(ctx, "GET", g.Config.APIURL, nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(io.LimitReader(res.Body, 1<<20))
	if err != nil {
		return nil, err
	}
	if c := res.StatusCode; c != 200 {
		return nil, fmt.Errorf("unexpected status code %d received from %s", c, g.Config.APIURL)
	}

	m, err := mapper.ThailandCBRResponseToRates(body)
	if err != nil {
		return nil, err
	}

	tz, err := time.LoadLocation(g.Config.Timezone)
	if err != nil {
		return nil, fmt.Errorf("bad location %s provided in the config, err %s", g.Config.Timezone, err)
	}
	dateStr := g.TimeNow().In(tz).Format(entity.DateLayout)
	return &entity.ExchangeRates{
		Country:    entity.Thailand,
		TimeZone:   tz,
		DateLoaded: dateStr,
		Rates:      m,
	}, nil
}

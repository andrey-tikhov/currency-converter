package integration_tests

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	uberconfig "go.uber.org/config"
	"go.uber.org/fx"
	"my_go/controller"
	"my_go/entity"
	"my_go/gateway/russia"
	"my_go/gateway/thailand"
	"my_go/handler"
	"my_go/logger"
	russiagatewaymock "my_go/mocks/gateway/russia"
	thailandgatewaymock "my_go/mocks/gateway/thailand"
	"my_go/repository"
	"net"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"
)

// This test ensures that caching logic + lazy interface provision by fx works as expected
// for 2 concurrent calls initial cache is created for first call and for 2nd call cache is used
// this is enforced by mock behaviour based on MaxTimes and MinTimes provided.
func TestCacheUsage(t *testing.T) {
	ruTZ, _ := time.LoadLocation("Europe/Moscow")
	ctrl := gomock.NewController(t)
	type expectedGatewayCalls struct {
		Russia   int
		Thailand int
	}
	t.Run("test that for 2 concurrent calls network is called only once", func(t *testing.T) {
		NewTestRUCBGateway := func() russia.Gateway {
			gw := russiagatewaymock.NewMockGateway(ctrl)
			gw.
				EXPECT().
				GetCBRRates(gomock.Any()).
				MaxTimes(1).
				MinTimes(1).
				Return(&entity.ExchangeRates{
					Country:    "russia",
					DateLoaded: time.Now().In(ruTZ).Format(entity.DateLayout),
					TimeZone:   ruTZ,
					Rates:      map[string]entity.Rate{},
				}, nil)
			return gw
		}
		NewTestTHCBGateway := func() thailand.Gateway {
			gw := thailandgatewaymock.NewMockGateway(ctrl)
			return gw
		}
		NewMux := func(lc fx.Lifecycle) *http.ServeMux {
			mux := http.NewServeMux()
			server := &http.Server{
				Addr:    "127.0.0.1:8000",
				Handler: mux,
			}
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					ln, err := net.Listen("tcp", server.Addr)
					if err != nil {
						return err
					}
					go server.Serve(ln)
					return nil
				},
				OnStop: func(ctx context.Context) error {
					return server.Shutdown(ctx)
				},
			})
			return mux
		}
		NewConfig := func() uberconfig.Provider {
			configOption := uberconfig.Source(strings.NewReader(`{"defaults":{"default_cb": "russia"}}`))
			provider, _ := uberconfig.NewYAML(configOption)
			return provider
		}
		Register := func(mux *http.ServeMux, h handler.Handler) {
			mux.HandleFunc("/get_exchange_rates", h.GetCBRates)
			mux.HandleFunc("/convert", h.ConvertCurrency)
			mux.HandleFunc("/hello", h.Hello)
		}
		app := fx.New(
			fx.Provide(NewTestRUCBGateway),
			fx.Provide(NewTestTHCBGateway),
			fx.Provide(NewMux),
			fx.Provide(NewConfig),
			handler.Module,
			controller.Module,
			logger.Module,
			repository.Module,
			fx.Invoke(Register),
		)
		startCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		err := app.Start(startCtx)
		assert.NoError(t, err)
		wg := &sync.WaitGroup{}
		wg.Add(2)
		go func() {
			req, _ := http.NewRequest(
				"POST",
				"http://localhost:8000/get_exchange_rates",
				strings.NewReader(`{"country":"russia"}`),
			)
			client := http.Client{
				Timeout: time.Second * 10,
			}
			_, err := client.Do(req)
			assert.NoError(t, err)
			wg.Done()
		}()
		go func() {
			req, _ := http.NewRequest(
				"POST",
				"http://localhost:8000/get_exchange_rates",
				strings.NewReader(`{"country":"russia"}`),
			)
			client := http.Client{
				Timeout: time.Second * 10,
			}
			_, err := client.Do(req)
			assert.NoError(t, err)
			wg.Done()
		}()
		wg.Wait()
		stopCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = app.Stop(stopCtx)
		assert.NoError(t, err)
	})
}

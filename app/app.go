package app

import (
	"go.uber.org/fx"
	"my_go/config"
	"my_go/controller"
	"my_go/gateway/russia"
	"my_go/gateway/thailand"
	"my_go/handler"
	"my_go/logger"
	"my_go/repository"
	"net/http"
)

var Module = fx.Options(
	handler.Module,
	config.Module,
	repository.Module,
	controller.Module,
	logger.Module,
	fx.Provide(russia.New),
	fx.Provide(thailand.New),
	fx.Invoke(StartAndListen),
)

func StartAndListen(h handler.Handler) {
	mux := http.NewServeMux()
	mux.HandleFunc("/get_exchange_rates", h.GetCBRates)
	mux.HandleFunc("/convert", h.ConvertCurrency)
	mux.HandleFunc("/hello", h.Hello)
	http.ListenAndServe(":8000", mux)
}

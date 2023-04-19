package controller

import (
	"go.uber.org/fx"
	"my_go/controller/cb_repository"
	"my_go/controller/conversion"
)

var Module = fx.Options(
	fx.Provide(cb_repository.New),
	fx.Provide(conversion.New),
)

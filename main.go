package main

import (
	"go.uber.org/fx"
	"my_go/app"
)

func opts() fx.Option {
	return fx.Options(
		app.Module,
	)
}

func main() {
	fx.New(opts()).Run()
}

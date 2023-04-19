package main

import (
	"fmt"
	"go.uber.org/fx"
	"my_go/app"
)

func opts() fx.Option {
	return fx.Options(
		app.Module,
	)
}

func main() {
	fmt.Println("tst")
	fx.New(opts()).Run()
}

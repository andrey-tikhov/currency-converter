package handler

import (
	"context"
	"fmt"
	"my_go/gateway/thailand"
)

func TestInvoke(g thailand.Gateway) {
	res, err := g.GetCBRRates(context.Background())
	fmt.Printf("res %+v err %s\n", res, err)
	fmt.Println("Hello world")
}

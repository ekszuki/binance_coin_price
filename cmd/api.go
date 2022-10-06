package main

import (
	"context"

	"br.com.sygnux/binance/pkg/infra/startup"
)

func main() {
	ctx := context.Background()

	apiStartup := startup.NewAPIStartup(ctx)
	apiStartup.Initialize()

}

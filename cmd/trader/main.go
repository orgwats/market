package main

import (
	"wats/config"
	"wats/internal/trader"
)

func main() {
	c := config.LoadConfig()
	t := trader.NewTrader(c)

	t.Start()
}

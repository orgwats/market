package chart

import (
	"wats/internal/chart/candle"
	"wats/internal/chart/indicators"
	"wats/internal/database"
)

type Chart struct {
	Candles    *candle.CandleBuffer
	Indicators *indicators.Indicators
}

func NewChart(db *database.Database, symbol string) *Chart {
	cb := candle.NewCandleBuffer(db, symbol)
	i := indicators.NewIndicators(cb)

	return &Chart{
		Candles:    cb,
		Indicators: i,
	}
}

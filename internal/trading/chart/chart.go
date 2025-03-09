package chart

import (
	"context"
	"sync"
	"wats/internal/database"
	"wats/internal/trading/chart/candle"
	"wats/internal/trading/chart/indicators"
	stream "wats/internal/trading/stream"
)

type Chart struct {
	ctx    context.Context
	symbol string

	CandleBuffer *candle.CandleBuffer
	Indicators   *indicators.Indicators
}

func NewChart(ctx context.Context, db *database.Database, symbol string) *Chart {
	c := &Chart{
		ctx:    ctx,
		symbol: symbol,

		CandleBuffer: candle.NewCandleBuffer(200),
		Indicators:   indicators.NewIndicators(),
	}

	cds := db.GetCandles(symbol, 200)
	c.CandleBuffer.Init(cds)

	return c
}

func (c *Chart) Run() {
	ch, stop := stream.NewStream(c.symbol)

	var once sync.Once
	for {
		select {
		case cd, ok := <-ch:
			if !ok {
				return
			}

			once.Do(func() {
				c.CandleBuffer.AddCandle(cd)
			})

			c.CandleBuffer.UpdateLastCandle(cd)
			c.Indicators.Update(c.CandleBuffer.GetCandles())

			if cd.Closed {
				c.CandleBuffer.AddCandle(cd)
				// DB INSERT
			}
		case <-c.ctx.Done():
			stop()
			return
		}
	}
}

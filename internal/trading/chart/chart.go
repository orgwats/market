package chart

import (
	"context"
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

		CandleBuffer: candle.NewCandleBuffer(30),
		Indicators:   indicators.NewIndicators(),
	}

	cds := db.GetCandles(symbol)
	c.CandleBuffer.Init(cds)

	return c
}

func (c *Chart) Run() {
	ch, stop := stream.NewStream(c.symbol)

	for {
		select {
		case cd, ok := <-ch:
			if !ok {
				return
			}

			if cd.Closed {
				c.CandleBuffer.AddCandle(cd)
				// TODO: DB INSERT 로직 필요
			} else {
				c.CandleBuffer.UpdateLastCandle(cd)
			}

			c.Indicators.Update(c.CandleBuffer.GetCandles())
		case <-c.ctx.Done():
			stop()
			return
		}
	}
}

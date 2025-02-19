package trading

import (
	"context"
	"wats/internal/database"
	"wats/internal/trading/analyzer"
	"wats/internal/trading/chart"
)

type Trader struct {
	ctx    context.Context
	cancel context.CancelFunc

	symbol string
	chart  *chart.Chart

	db *database.Database

	// services
	analyzer *analyzer.Analyzer
}

func NewTrader(db *database.Database, symbol string) *Trader {
	ctx, cancel := context.WithCancel(context.Background())

	c := chart.NewChart(ctx, db, symbol)
	a := analyzer.NewAnalyzer(ctx, c)

	return &Trader{
		ctx:    ctx,
		cancel: cancel,

		symbol: symbol,
		chart:  c,

		db: db,

		analyzer: a,
	}
}

func (t *Trader) Start() {
	go t.chart.Run()
	go t.analyzer.Strat()
}

func (t *Trader) Stop() {
	t.cancel()
}

package analyzer

import (
	"context"
	"wats/internal/trading/chart"
)

type Analyzer struct {
	ctx context.Context

	chart *chart.Chart
}

func NewAnalyzer(ctx context.Context, chart *chart.Chart) *Analyzer {
	return &Analyzer{
		ctx: ctx,

		chart: chart,
	}
}

func (a *Analyzer) Strat() {
}

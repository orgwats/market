package analyzer

import (
	"context"
	"log"
	"wats/internal/types"
)

type Analyzer struct {
	ctx context.Context

	channel *types.TraderChannel
}

func NewAnalyzer(ctx context.Context, ch *types.TraderChannel) *Analyzer {
	return &Analyzer{
		ctx: ctx,

		channel: ch,
	}
}

func (a *Analyzer) Strat() {
	// TODO: 임시 선언 > 추후 설정에서 불러오는 방식으로 수정
	symbols := map[string]bool{
		"TRUMPUSDT": true,
	}

	for e := range a.channel.Stream {
		for k := range symbols {
			log.Printf("symbol: %s, price: %s", e[k].Symbol, e[k].MarkPrice)
		}
	}
}

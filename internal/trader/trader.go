package trader

import (
	"context"
	"log"
	"wats/config"
	"wats/internal/analyzer"
	"wats/internal/streamer"
	"wats/internal/types"

	"github.com/adshao/go-binance/v2/futures"
)

type Trader struct {
	ctx    context.Context
	cancel context.CancelFunc

	// services
	streamer *streamer.Streamer
	analyzer *analyzer.Analyzer

	channel *types.TraderChannel
}

func NewTrader(config *config.Config) *Trader {
	ctx, cancel := context.WithCancel(context.Background())

	// 채널 초기화
	ch := &types.TraderChannel{
		Stream: make(chan map[string]futures.WsMarkPriceEvent),
		Order:  make(chan struct{}),
		Done:   make(chan error),
	}

	s := streamer.NewStreamer(ctx, ch)
	a := analyzer.NewAnalyzer(ctx, ch)

	return &Trader{
		ctx:    ctx,
		cancel: cancel,

		streamer: s,
		analyzer: a,

		channel: ch,
	}
}

func (t *Trader) Start() {
	log.Println("[trader] starting...")
	go t.streamer.Start()
	go t.analyzer.Strat()
}

func (t *Trader) Stop() {
	log.Println("[trader] stopping...")
	t.cancel()
	log.Println("[trader] stopped.")
}

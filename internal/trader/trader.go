package trader

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"wats/config"
	"wats/internal/streamer"
)

type Trader interface {
	Start()
}

type TraderImpl struct {
	config *config.Config

	ctx    context.Context
	cancel context.CancelFunc
}

func NewTrader(config *config.Config) *TraderImpl {
	t := &TraderImpl{
		config: config,
	}

	t.ctx, t.cancel = context.WithCancel(context.Background())

	return t
}

func (t *TraderImpl) Start() {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		waitForStopSignal()
		t.cancel()
	}()

	config := t.config
	ctx := t.ctx

	s := streamer.NewStreamer(ctx, config.WebsocketURL)

	go s.Start(&wg)

	// 모든 goroutine이 종료될 때까지 대기
	wg.Wait()
	log.Println("Stopped wats.")
}

func waitForStopSignal() {
	// OS 시그널 처리
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// 시그널 대기
	<-sigCh
}

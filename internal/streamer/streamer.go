package streamer

import (
	"context"
	"errors"
	"log"
	"time"
	"wats/internal/types"

	"github.com/adshao/go-binance/v2/futures"
)

type Streamer struct {
	ctx context.Context

	channel *types.TraderChannel
}

func NewStreamer(ctx context.Context, ch *types.TraderChannel) *Streamer {
	return &Streamer{
		ctx: ctx,

		channel: ch,
	}
}

func (s *Streamer) Start() {
	doneC, stopC, err := futures.WsAllMarkPriceServeWithRate(
		1*time.Second,
		s.handleServe,
		s.handleError,
	)
	if err != nil {
		s.channel.Done <- errors.New("[streamer] failed to websocket connect")
		return
	}

	log.Println("[streamer] websoket connected successfully")

	select {
	case <-s.ctx.Done():
		log.Println("[streamer] received stop signal, close websocket.")
		close(stopC)
		return
	case <-doneC:
		log.Println("[streamer] websoket closed by server.")
		s.channel.Done <- nil
		return
	}
}

func (s *Streamer) handleServe(events futures.WsAllMarkPriceEvent) {
	// TODO: 임시 선언 > 추후 설정에서 불러오는 방식으로 수정
	symbols := map[string]bool{
		"TRUMPUSDT": true,
	}

	m := make(map[string]futures.WsMarkPriceEvent, 3)

	for _, e := range events {
		if symbols[e.Symbol] {
			m[e.Symbol] = *e
		}
	}

	s.channel.Stream <- m
}

func (s *Streamer) handleError(err error) {
	log.Printf("[streamer] error: %v\n", err)
}

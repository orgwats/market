package candle

import (
	"wats/internal/database"
	s "wats/internal/stream"
	"wats/internal/types"
)

type RingBuffer struct {
	buf  []*types.Candle
	head int
	size int
}

type CandleBuffer struct {
	ch   chan *types.Candle
	ring *RingBuffer
}

func NewCandleBuffer(db *database.Database, symbol string) *CandleBuffer {
	cb := &CandleBuffer{
		ch: s.NewStream(symbol),
		ring: &RingBuffer{
			buf:  make([]*types.Candle, 30),
			head: 0,
			size: 30,
		},
	}

	for _, c := range db.GetCandles(symbol) {
		// Candle ring buffer 초기화
		cb.AddCandle(c)
	}

	return cb
}

func (cb *CandleBuffer) ReceiveCandle() *types.Candle {
	return <-cb.ch
}

func (cb *CandleBuffer) GetCandles() []*types.Candle {
	return cb.ring.buf
}

func (cb *CandleBuffer) AddCandle(c *types.Candle) {
	cb.ring.buf[cb.ring.head] = c
	cb.ring.head = (cb.ring.head + 1) % cb.ring.size
}

func (cb *CandleBuffer) UpdateLastCandle(c *types.Candle) {
	idx := (cb.ring.head - 1 + cb.ring.size) % cb.ring.size
	cb.ring.buf[idx] = c
}

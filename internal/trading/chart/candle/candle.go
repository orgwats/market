package candle

import (
	"wats/internal/types"
)

type CandleBuffer struct {
	buf  []*types.Candle
	head int
	size int
}

func NewCandleBuffer(size int) *CandleBuffer {
	return &CandleBuffer{
		buf:  make([]*types.Candle, size),
		head: 0,
		size: size,
	}
}

func (cb *CandleBuffer) Init(cds []*types.Candle) {
	for _, c := range cds {
		cb.AddCandle(c)
	}
}

func (cb *CandleBuffer) GetCandles() []*types.Candle {
	return cb.buf
}

func (cb *CandleBuffer) AddCandle(c *types.Candle) {
	cb.buf[cb.head] = c
	cb.head = (cb.head + 1) % cb.size
}

func (cb *CandleBuffer) UpdateLastCandle(c *types.Candle) {
	idx := (cb.head - 1 + cb.size) % cb.size
	cb.buf[idx] = c
}

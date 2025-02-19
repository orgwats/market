package indicators

import "wats/internal/types"

type Indicators struct {
	RSI chan float64
}

func NewIndicators() *Indicators {
	return &Indicators{
		RSI: make(chan float64),
	}
}

func (i *Indicators) Update(candles []*types.Candle) {
	i.RSI <- CalculateRSI(candles, 14)
}

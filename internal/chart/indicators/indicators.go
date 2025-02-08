package indicators

import (
	"log"
	"wats/internal/chart/candle"
)

type Indicators struct {
	RSI chan float64
}

func NewIndicators(cb *candle.CandleBuffer) *Indicators {
	i := &Indicators{
		RSI: make(chan float64),
	}

	go func() {
		for {
			c := cb.ReceiveCandle()
			log.Println(c)
			// 첫 스트림 데이터를 추가해주기 위해 once.Do 필요

			cb.AddCandle(c)
			log.Println(c, CalculateRSI(cb.GetCandles(), 14))
		}
	}()

	return i
}

package indicators

import (
	"log"
	"math"
	"wats/internal/types"
)

func CalculateRSI(candles []*types.Candle, period int) float64 {
	if len(candles) < period {
		log.Println("Not enough data to calculate RSI")
		return 0
	}

	var gains, losses float64

	// 평균 상승/하락폭 계산
	for i := 1; i < period; i++ {
		diff := candles[i].Close - candles[i-1].Close
		if diff > 0 {
			gains += diff
		} else {
			losses -= -diff
		}
	}

	// 평균 값
	avgGain := gains / float64(period)
	avgLoss := losses / float64(period)

	for i := period; i < len(candles); i++ {
		diff := candles[i].Close - candles[i-1].Close
		if diff > 0 {
			avgGain = ((avgGain * float64(period-1)) + diff) / float64(period)
			avgLoss = ((avgLoss * float64(period-1)) + 0) / float64(period)
		} else {
			avgGain = ((avgGain * float64(period-1)) + 0) / float64(period)
			avgLoss = ((avgLoss * float64(period-1)) - diff) / float64(period)
		}
	}

	// RSI 계산
	if avgLoss == 0 {
		return 100
	}
	rs := avgGain / avgLoss
	rsi := 100 - (100 / (1 + rs))

	return math.Min(rsi, 100)
}

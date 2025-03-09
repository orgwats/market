package indicators

import (
	"log"
	"wats/internal/types"
)

func CalculateRSI(candles []*types.Candle, period int) float64 {
	if len(candles) < period {
		log.Println("Not enough data to calculate RSI")
		return 0
	}

	var gains, losses []float64

	for i := 1; i < len(candles); i++ {
		diff := candles[i].Close - candles[i-1].Close
		if diff > 0 {
			gains = append(gains, diff)
			losses = append(losses, 0)
		} else {
			gains = append(gains, 0)
			losses = append(losses, -diff)
		}
	}

	au := ComputeWellesWiderMA(gains, float64(period))
	ad := ComputeWellesWiderMA(losses, float64(period))

	rs := au / ad

	return (rs / (1 + rs)) * 100
}

func ComputeWellesWiderMA(prices []float64, period float64) float64 {
	k := 1.0 / period
	ma := make([]float64, len(prices))
	ma[0] = prices[0]

	for i := 1; i < len(prices); i++ {
		ma[i] = (prices[i]*k + ma[i-1]*(1-k))
	}

	return ma[len(prices)-1]
}

package streame

import (
	"log"
	"strconv"
	"wats/internal/types"

	"github.com/adshao/go-binance/v2/futures"
)

func NewStream(symbol string) chan *types.Candle {
	ch := make(chan *types.Candle)

	go func() {
		doneC, _, _ := futures.WsKlineServe(
			symbol,
			"1m",
			func(event *futures.WsKlineEvent) {
				ch <- parseEvent(event)
			},
			func(err error) {
				log.Printf("[stream] error: %v\n", err)
			},
		)
		<-doneC
	}()

	return ch
}

func parseEvent(event *futures.WsKlineEvent) *types.Candle {
	parseFloat := func(s string) float64 {
		val, _ := strconv.ParseFloat(s, 64)
		return val
	}

	return &types.Candle{
		Symbol:              event.Symbol,
		OpenTime:            event.Kline.StartTime,
		Open:                parseFloat(event.Kline.Open),
		High:                parseFloat(event.Kline.High),
		Low:                 parseFloat(event.Kline.Low),
		Close:               parseFloat(event.Kline.Close),
		Volume:              parseFloat(event.Kline.Volume),
		CloseTime:           event.Kline.EndTime,
		QuoteVolume:         parseFloat(event.Kline.QuoteVolume),
		Count:               int(event.Kline.TradeNum),
		TakerBuyVolume:      parseFloat(event.Kline.ActiveBuyVolume),
		TakerBuyQuoteVolume: parseFloat(event.Kline.ActiveBuyQuoteVolume),
		Closed:              event.Kline.IsFinal,
	}
}

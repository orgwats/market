package stream

import (
	"log"

	"github.com/orgwats/stream/internal/config"

	"github.com/adshao/go-binance/v2/futures"
	pb "github.com/orgwats/idl/gen/go/market"
)

func NewStream(cfg *config.Config) (ch chan *pb.Candle) {
	ch = make(chan *pb.Candle)
	doneC, _ := connectKlineWebsocket(cfg.Symbols, ch)

	go func() {
		<-doneC
		close(ch)
	}()

	// TODO: 추가 필요
	// stop = func() {
	// 	close(stopC)
	// }

	return ch
}

func connectKlineWebsocket(symbols []string, ch chan *pb.Candle) (doneC, stopC chan struct{}) {
	log.Println(symbols)
	symbolIntervalPair := make(map[string]string)

	for _, symbol := range symbols {
		symbolIntervalPair[symbol] = "1m"
	}

	doneC, stopC, _ = futures.WsCombinedKlineServe(
		symbolIntervalPair,
		func(event *futures.WsKlineEvent) {
			ch <- parseEvent(event)
		},
		func(err error) {
			log.Printf("WsKlineServe error: %v\n", err)
		},
	)

	return doneC, stopC
}

func parseEvent(event *futures.WsKlineEvent) *pb.Candle {
	return &pb.Candle{
		Symbol:              event.Symbol,
		OpenTime:            event.Kline.StartTime,
		Open:                event.Kline.Open,
		High:                event.Kline.High,
		Low:                 event.Kline.Low,
		Close:               event.Kline.Close,
		Volume:              event.Kline.Volume,
		CloseTime:           event.Kline.EndTime,
		QuoteVolume:         event.Kline.QuoteVolume,
		Count:               event.Kline.TradeNum,
		TakerBuyVolume:      event.Kline.ActiveBuyVolume,
		TakerBuyQuoteVolume: event.Kline.ActiveBuyQuoteVolume,
		Closed:              event.Kline.IsFinal,
	}
}

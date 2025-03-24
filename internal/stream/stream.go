package stream

import (
	"log"

	"github.com/orgwats/stream/internal/config"

	"github.com/adshao/go-binance/v2/futures"
	pb "github.com/orgwats/idl/gen/go/market"
)

type Stream struct {
	Ch   chan *pb.Candle
	Stop func()
}

func NewStream(cfg *config.Config) (*Stream, error) {
	ch := make(chan *pb.Candle)
	doneC, stopC, err := connectKlineWebsocket(cfg.Symbols, ch)
	if err != nil {
		return nil, err
	}

	stream := &Stream{
		Ch: ch,
		Stop: func() {
			close(stopC)
		},
	}

	go func() {
		<-doneC
		close(ch)
	}()

	return stream, err
}

func connectKlineWebsocket(symbols []string, ch chan *pb.Candle) (doneC, stopC chan struct{}, err error) {
	symbolIntervalPair := make(map[string]string)

	for _, symbol := range symbols {
		symbolIntervalPair[symbol] = "1m"
	}

	doneC, stopC, err = futures.WsCombinedKlineServe(
		symbolIntervalPair,
		func(event *futures.WsKlineEvent) {
			ch <- parseEvent(event)
		},
		func(err error) {
			log.Printf("WsKlineServe error: %v\n", err)
		},
	)
	if err != nil {
		return nil, nil, err
	}

	return doneC, stopC, nil
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

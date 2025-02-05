package types

import "github.com/adshao/go-binance/v2/futures"

type TraderChannel struct {
	Stream chan map[string]futures.WsMarkPriceEvent
	Order  chan struct{}
	Done   chan error
}

type Candle struct {
	Symbol              string
	OpenTime            int64
	Open                float64
	High                float64
	Low                 float64
	Close               float64
	Volume              float64
	CloseTime           int64
	QuoteVolume         float64
	Count               int
	TakerBuyVolume      float64
	TakerBuyQuoteVolume float64
	Closed              bool
}

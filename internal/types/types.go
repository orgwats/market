package types

import "github.com/adshao/go-binance/v2/futures"

type TraderChannel struct {
	Stream chan map[string]futures.WsMarkPriceEvent
	Order  chan struct{}
	Done   chan error
}

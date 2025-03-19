package hub

import (
	"sync"

	"github.com/google/uuid"
	pb "github.com/orgwats/idl/gen/go/market"
)

type Hub struct {
	subscribers map[string]map[uuid.UUID]chan *pb.Candle
	mu          sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		subscribers: make(map[string]map[uuid.UUID]chan *pb.Candle),
	}
}

func (h *Hub) AddSubscriber(symbol string) (uuid.UUID, <-chan *pb.Candle) {
	h.mu.Lock()
	defer h.mu.Unlock()

	id, _ := uuid.NewUUID()
	ch := make(chan *pb.Candle, 1)

	if h.subscribers[symbol] == nil {
		h.subscribers[symbol] = make(map[uuid.UUID]chan *pb.Candle)
	}

	h.subscribers[symbol][id] = ch
	return id, ch
}

func (h *Hub) RemoveSubscriber(symbol string, id uuid.UUID) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if ch, ok := h.subscribers[symbol][id]; ok {
		close(ch)
		delete(h.subscribers[symbol], id)
	}
}

func (h *Hub) Broadcast(symbol string, c *pb.Candle) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, ch := range h.subscribers[symbol] {
		select {
		case ch <- c:
		default:
			select {
			case <-ch:
			default:
			}
			ch <- c
		}
	}
}

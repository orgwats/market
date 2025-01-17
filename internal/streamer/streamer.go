package streamer

import (
	"context"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Streamer interface {
	Start(wg *sync.WaitGroup)
}

type StreamerImpl struct {
	context context.Context
	url     string
}

func NewStreamer(context context.Context, url string) *StreamerImpl {
	return &StreamerImpl{
		context: context,
		url:     url,
	}
}

func (s *StreamerImpl) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	// 1. WebSocket 연결
	conn, _, err := websocket.DefaultDialer.Dial(s.url, nil)

	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	defer conn.Close()
	log.Println("Connected to:", s.url)

	// 2. 메시지 수신 루프
	for {
		select {
		case <-s.context.Done():
			log.Printf("Streamer Context canceled. Exiting.")
			conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			return
		default:
			_, message, err := conn.ReadMessage()

			if err != nil {
				log.Printf("Streamer Read error: %v", err)
				return
			}

			log.Println("message : ", message)
		}
	}
}

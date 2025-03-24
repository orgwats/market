package gapi

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/orgwats/market/internal/config"
	db "github.com/orgwats/market/internal/db/sqlc"
	"github.com/orgwats/market/internal/hub"
	"github.com/orgwats/market/internal/stream"

	pb "github.com/orgwats/idl/gen/go/market"
)

type Server struct {
	pb.UnimplementedMarketServer

	// TODO: 임시
	cfg   *config.Config
	store db.Store

	binanceClient *futures.Client
	stream        *stream.Stream
	hub           *hub.Hub
}

func NewServer(cfg *config.Config, store db.Store) *Server {
	hub := hub.NewHub()
	client := futures.NewClient(cfg.BinanceApiKey, cfg.BinanceSecretKey)

	return &Server{
		cfg:           cfg,
		store:         store,
		binanceClient: client,
		hub:           hub,
	}
}

func (s *Server) Run() {
	err := s.sync()
	if err != nil {
		log.Fatal("cannot sync candle data:", err)
	}

	stream, err := stream.NewStream(s.cfg)
	if err != nil {
		log.Fatal("cannot create stream:", err)
	}

	s.stream = stream

	ch := make(chan *pb.Candle)

	go func() {
		for c := range s.stream.Ch {
			s.hub.Broadcast(c.Symbol, c)

			if c.Closed {
				ch <- c
			}
		}
	}()

	go func() {
		for c := range ch {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			err := s.store.SaveCandle(ctx, db.SaveCandleParams{
				Symbol:              c.Symbol,
				OpenTime:            c.OpenTime,
				Open:                c.Open,
				High:                c.High,
				Low:                 c.Low,
				Close:               c.Close,
				Volume:              c.Volume,
				CloseTime:           c.CloseTime,
				QuoteVolume:         c.QuoteVolume,
				Count:               c.Count,
				TakerBuyVolume:      c.TakerBuyVolume,
				TakerBuyQuoteVolume: c.TakerBuyQuoteVolume,
			})
			if err != nil {
				// 필요에 따라 재시도 로직
			}
			cancel()
		}
	}()
}

func (s *Server) Stop() {
	s.stream.Stop()
}

func (s *Server) sync() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	baseService := *s.binanceClient.NewKlinesService().Interval("1m").Limit(200)

	ch := make(chan db.Candle)

	var wg sync.WaitGroup
	for _, symbol := range s.cfg.Symbols {
		wg.Add(1)
		go func(service futures.KlinesService, symbol string) {
			defer wg.Done()

			svc := service.Symbol(symbol)

			lastestCandle, err := s.store.GetLatestCandle(ctx, symbol)
			if err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					// 예외 처리 필요
				}
			} else {
				svc.StartTime(lastestCandle.OpenTime + 1)
			}

			klines, err := svc.Do(ctx)
			if err != nil {
				// 필요에 따라 재시도 로직
			}

			for _, kline := range klines {
				ch <- db.Candle{
					Symbol:              symbol,
					OpenTime:            kline.OpenTime,
					Open:                kline.Open,
					High:                kline.High,
					Low:                 kline.Low,
					Close:               kline.Close,
					Volume:              kline.Volume,
					CloseTime:           kline.CloseTime,
					QuoteVolume:         kline.QuoteAssetVolume,
					Count:               kline.TradeNum,
					TakerBuyVolume:      kline.TakerBuyBaseAssetVolume,
					TakerBuyQuoteVolume: kline.TakerBuyQuoteAssetVolume,
				}
			}

		}(baseService, symbol)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	var candles []db.Candle

	for candle := range ch {
		candles = append(candles, candle)
	}

	err := s.store.SaveCandles(ctx, candles)
	if err != nil {
		return err
	}

	log.Printf("SYNCED!!! : %d Rows", len(candles))

	return nil
}

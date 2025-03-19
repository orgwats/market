package gapi

import (
	"context"
	"log"

	db "github.com/orgwats/stream/internal/db/sqlc"

	pb "github.com/orgwats/idl/gen/go/market"
)

func (s *Server) GetCandles(ctx context.Context, req *pb.GetCandlesRequest) (*pb.GetCandlesResponse, error) {
	arg := db.GetCandlesParams{
		Symbol: req.Symbol,
		Limit:  req.Limit,
	}

	dbCandles, err := s.store.GetCandles(ctx, arg)
	if err != nil {
		log.Println(err)
	}

	var pbCandles []*pb.Candle

	for _, c := range dbCandles {
		pbCandle := &pb.Candle{
			Symbol:              c.Symbol,
			OpenTime:            c.OpenTime,
			Open:                c.Open,
			High:                c.High,
			Low:                 c.Low,
			Close:               c.Close,
			Volume:              c.Volume,
			CloseTime:           c.CloseTime,
			QuoteVolume:         c.QuoteVolume,
			Count:               int64(c.Count),
			TakerBuyVolume:      c.TakerBuyVolume,
			TakerBuyQuoteVolume: c.TakerBuyQuoteVolume,
			Closed:              true,
		}
		pbCandles = append(pbCandles, pbCandle)
	}

	return &pb.GetCandlesResponse{Candles: pbCandles}, nil
}

package gapi

import (
	pb "github.com/orgwats/idl/gen/go/market"
)

func (s *Server) NewCandleStream(req *pb.NewCandleStreamRequest, stream pb.Market_NewCandleStreamServer) error {
	ctx := stream.Context()

	// 구독자 추가: 고유 ID와 채널을 받아옴
	subscriberID, subCh := s.hub.AddSubscriber(req.Symbol)
	defer s.hub.RemoveSubscriber(req.Symbol, subscriberID)

	// 구독 채널로부터 계속 데이터를 읽어 전송
	for {
		select {
		case <-ctx.Done():
			// 클라이언트 연결 종료 시 종료
			return ctx.Err()
		case c, ok := <-subCh:
			if !ok {
				// 허브에서 채널이 닫히면 종료
				return nil
			}

			candle := &pb.Candle{
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
				Closed:              c.Closed,
			}

			if err := stream.Send(candle); err != nil {
				return err
			}
		}
	}
}

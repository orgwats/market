package trader

import (
	"context"
	"log"
	"strconv"
	"sync"
	"time"
	"wats/internal/database"
	"wats/internal/trading/analyzer"
	"wats/internal/trading/chart"
	"wats/internal/trading/market"
	"wats/internal/types"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
)

type Trader struct {
	ctx    context.Context
	cancel context.CancelFunc

	db *database.Database

	symbol string
}

func NewTrader(db *database.Database, symbol string) *Trader {
	ctx, cancel := context.WithCancel(context.Background())

	return &Trader{
		ctx:    ctx,
		cancel: cancel,

		db: db,

		symbol: symbol,
	}
}

func (t *Trader) Start() {
	// Candle 데이터 DB sync 작업
	t.syncCandleData()

	c := chart.NewChart(t.ctx, t.db, t.symbol)
	a := analyzer.NewAnalyzer(t.ctx, c)

	go c.Run()
	go a.Strat()
}

func (t *Trader) Stop() {
	t.cancel()
}

func (t *Trader) syncCandleData() {
	today := time.Now()
	yesterday := today.AddDate(0, 0, -1)

	// DB에서 가장 최신 캔들 조회
	lc := t.db.GetLatestCandle(t.symbol)

	var start, end time.Time

	if lc == nil {
		// yesterday - 29 ~ yesterday
		start = yesterday.AddDate(0, 0, -29)
		end = yesterday
	} else if d := time.UnixMilli(int64(lc.CloseTime)); d.Before(yesterday) {
		// 마지막 캔들 데이터 ~ yesterday
		start = time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
		end = yesterday
	}

	if !start.IsZero() && !end.IsZero() {
		m := market.NewMarketData(t.symbol)
		fpCh := make(chan string)

		// 캔들 데이터 다운로드
		var wgDownload sync.WaitGroup
		for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
			wgDownload.Add(1)
			go func(date time.Time) {
				defer wgDownload.Done()
				select {
				case <-t.ctx.Done():
					return
				default:
					fp, err := m.DownloadFile(date.Year(), int(date.Month()), date.Day())
					if err != nil {
						log.Printf("%d-%02d-%02d 데이터 다운로드 실패: %s", date.Year(), int(date.Month()), date.Day(), err)
						return
					}
					fpCh <- fp
				}
			}(d)
		}

		go func() {
			wgDownload.Wait()
			close(fpCh)
		}()

		const dbWorkers = 5

		var wgDB sync.WaitGroup
		for i := 0; i < dbWorkers; i++ {
			wgDB.Add(1)
			go func() {
				defer wgDB.Done()
				for fp := range fpCh {
					select {
					case <-t.ctx.Done():
						return
					default:
						err := t.db.SaveCandlesFromCSV(t.symbol, fp)
						if err != nil {
							log.Printf("%s 데이터 저장 실패: %v\n", fp, err)
							continue
						}
					}
				}
			}()
		}

		wgDB.Wait()
	}

	client := binance.NewFuturesClient(
		"",
		"",
	)

	var cds []*types.Candle

	lc = t.db.GetLatestCandle(t.symbol)
	startTime := lc.CloseTime

	for {
		klines, err := client.NewKlinesService().
			// cfg 파일로 분리 필요 (임시)
			Symbol("ETHUSDT").
			Interval("1m").
			StartTime(startTime).
			Limit(1500).
			Do(t.ctx)
		if err != nil {
			log.Println("API 요청 실패 : ", err)
			break
		}

		isBreak := false

		for i, k := range klines {
			c := convertKlineToCandle(k)

			if i == len(klines)-1 {
				lcSecond := time.UnixMilli(int64(lc.CloseTime)).Second()
				cSecond := time.UnixMilli(int64(c.CloseTime)).Second()

				if lcSecond == cSecond {
					isBreak = true
				} else {
					startTime = lc.CloseTime
					cds = append(cds, c)
					lc = c
				}
			} else {
				cds = append(cds, c)
			}
		}

		if isBreak {
			break
		}
	}

	err := t.db.SaveCandles(t.symbol, cds)
	if err != nil {
		log.Println(err)
	}

	log.Println("sync 작업 종료.")
}

func convertKlineToCandle(k *futures.Kline) *types.Candle {
	open, _ := strconv.ParseFloat(k.Open, 64)
	high, _ := strconv.ParseFloat(k.High, 64)
	low, _ := strconv.ParseFloat(k.Low, 64)
	close, _ := strconv.ParseFloat(k.Close, 64)
	volume, _ := strconv.ParseFloat(k.Volume, 64)
	quoteVolume, _ := strconv.ParseFloat(k.QuoteAssetVolume, 64)
	takerBuyVolume, _ := strconv.ParseFloat(k.TakerBuyBaseAssetVolume, 64)
	takerBuyQuoteVolume, _ := strconv.ParseFloat(k.TakerBuyQuoteAssetVolume, 64)

	return &types.Candle{
		Symbol:              "",
		OpenTime:            k.OpenTime,
		Open:                open,
		High:                high,
		Low:                 low,
		Close:               close,
		Volume:              volume,
		CloseTime:           k.CloseTime,
		QuoteVolume:         quoteVolume,
		Count:               int(k.TradeNum),
		TakerBuyVolume:      takerBuyVolume,
		TakerBuyQuoteVolume: takerBuyQuoteVolume,
	}
}

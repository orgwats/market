package trader

import (
	"context"
	"log"
	"sync"
	"time"
	"wats/internal/database"
	"wats/internal/trading/analyzer"
	"wats/internal/trading/chart"
	"wats/internal/trading/market"
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
						log.Printf("%d-%02d-%02d 데이터 다운로드 실패: %w", date.Year(), int(date.Month()), date.Day(), err)
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
}

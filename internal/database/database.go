package database

import (
	"database/sql"
	"log"
	"wats/internal/types"

	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	db *sql.DB
}

func NewDatabase() *Database {
	// TODO: 환경 변수로 이동
	dsn := "root:123456@tcp(localhost:3306)/binance"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Database connection failed:", err)
	}

	return &Database{db: db}
}

// GetCandles: Candle 데이터를 조회하고 배열로 반환
func (d *Database) GetCandles(symbol string) []*types.Candle {
	query := `
		SELECT *
		FROM candle
		WHERE symbol = ?
		ORDER BY open_time DESC
		LIMIT 30;
	`

	var candles []*types.Candle

	rows, err := d.db.Query(query, symbol)
	if err != nil {
		// 에러 처리 로직 추가 필요
		return candles
	}
	defer rows.Close()

	for rows.Next() {
		candle := &types.Candle{Closed: true}
		err := rows.Scan(&candle.Symbol, &candle.OpenTime, &candle.Open, &candle.High, &candle.Low, &candle.Close, &candle.Volume, &candle.CloseTime, &candle.QuoteVolume, &candle.Count, &candle.TakerBuyVolume, &candle.TakerBuyQuoteVolume)
		if err != nil {
			// 에러 처리 로직 추가 필요
			return candles
		}
		candles = append(candles, candle)
	}
	return candles
}

func (d *Database) GetLatestCandle(symbol string) *types.Candle {
	query := `
		SELECT *
		FROM candle
		WHERE symbol = ?
		ORDER BY open_time DESC
		LIMIT 1;
	`

	row := d.db.QueryRow(query, symbol)
	candle := &types.Candle{Closed: true}
	err := row.Scan(&candle.Symbol, &candle.OpenTime, &candle.Open, &candle.High, &candle.Low, &candle.Close, &candle.Volume, &candle.CloseTime, &candle.QuoteVolume, &candle.Count, &candle.TakerBuyVolume, &candle.TakerBuyQuoteVolume)
	if err != nil {
		return nil
	}

	return candle
}

func (d *Database) Close() {
	d.db.Close()
}

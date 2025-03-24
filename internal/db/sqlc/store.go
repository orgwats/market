package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type Store interface {
	Querier
	SaveCandles(ctx context.Context, candles []Candle) error
	SaveCandlesFromCSV(ctx context.Context, symbol, filePath string) error
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *SQLStore {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

func (s *SQLStore) SaveCandlesFromCSV(ctx context.Context, symbol, filePath string) error {
	query := fmt.Sprintf(`
		LOAD DATA LOCAL INFILE './%s'
		INTO TABLE candles
		FIELDS TERMINATED BY ','
		LINES TERMINATED BY '\n'
		IGNORE 1 LINES
		(open_time, open, high, low, close, volume, close_time, quote_volume, count, taker_buy_volume, taker_buy_quote_volume, @ignore)
		SET symbol = '%s'
	`, filePath, symbol)

	_, err := s.db.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (s *SQLStore) SaveCandles(ctx context.Context, candles []Candle) error {
	const batchSize = 1000
	for start := 0; start < len(candles); start += batchSize {
		end := start + batchSize
		if end > len(candles) {
			end = len(candles)
		}

		query, args := generateBulkInsertQuery(candles[start:end])
		_, err := s.db.ExecContext(ctx, query, args...)
		if err != nil {
			return err
		}
	}
	return nil
}

func generateBulkInsertQuery(chunk []Candle) (query string, arg []interface{}) {
	var sb strings.Builder
	sb.WriteString("INSERT INTO candles (symbol, open_time, open, high, low, close, volume, close_time, count, quote_volume, taker_buy_volume, taker_buy_quote_volume) VALUES ")

	placeholders := make([]string, 0, len(chunk))
	args := make([]interface{}, 0, len(chunk)*12)

	for _, c := range chunk {
		placeholders = append(placeholders, "(?,?,?,?,?,?,?,?,?,?,?,?)")
		args = append(args,
			c.Symbol,
			c.OpenTime,
			c.Open,
			c.High,
			c.Low,
			c.Close,
			c.Volume,
			c.CloseTime,
			c.Count,
			c.QuoteVolume,
			c.TakerBuyVolume,
			c.TakerBuyQuoteVolume,
		)
	}

	sb.WriteString(strings.Join(placeholders, ","))

	return sb.String(), args
}

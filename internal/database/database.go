package database

import (
	"database/sql"
	"fmt"
	"log"
	"wats/internal/types"

	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	db *sql.DB
}

func NewDatabase() *Database {
	// TODO: 환경 변수로 이동
	dsn := "root:123456@tcp(localhost:3306)/binance?allowAllFiles=true"

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

func (d *Database) SaveCandlesFromCSV(symbol, filePath string) error {
	query := fmt.Sprintf(`
		LOAD DATA LOCAL INFILE './%s'
		INTO TABLE candle
		FIELDS TERMINATED BY ','
		LINES TERMINATED BY '\n'
		IGNORE 1 LINES
		(open_time, open, high, low, close, volume, close_time, quote_volume, count, taker_buy_volume, taker_buy_quote_volume, @ignore)
		SET symbol = '%s'
	`, filePath, symbol)

	_, err := d.db.Exec(query)
	if err != nil {
		return err
	}


func (d *Database) SaveCandles(symbol string, cds []*types.Candle) error {
	const batchSize = 1000
	for start := 0; start < len(cds); start += batchSize {
		end := start + batchSize
		if end > len(cds) {
			end = len(cds)
		}
		chunk := cds[start:end]
		if err := d.bulkInsert(symbol, chunk); err != nil {
			return err
		}
	}
	return nil
}

func (d *Database) bulkInsert(symbol string, chunk []*types.Candle) error {
	var sb strings.Builder
	sb.WriteString("INSERT INTO candle (symbol, open_time, open, high, low, close, volume, close_time, count, quote_volume, taker_buy_volume, taker_buy_quote_volume) VALUES ")

	placeholders := make([]string, 0, len(chunk))
	args := make([]interface{}, 0, len(chunk)*12)

	for _, c := range chunk {
		placeholders = append(placeholders, "(?,?,?,?,?,?,?,?,?,?,?,?)")
		args = append(args,
			symbol,
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
	query := sb.String()

	result, err := d.db.Exec(query, args...)
	if err != nil {
		return err
	}
	c, _ := result.RowsAffected()
	log.Printf("bulkInsert 완료 : %d행 추가", c)
	return err
}

func (d *Database) Close() {
	d.db.Close()
}

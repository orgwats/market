package market

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"wats/internal/types"
)

type MarketData struct {
	symbol string
}

func NewMarketData(symbol string) *MarketData {
	return &MarketData{
		symbol: symbol,
	}
}

func (m *MarketData) GetCandles(year, month, day int) ([]*types.Candle, error) {
	i := "monthly"
	d := fmt.Sprintf("%d-%02d", year, month)
	if day > 0 {
		i = "daily"
		d = fmt.Sprintf("%s-%02d", d, day)
	}

	url := fmt.Sprintf(
		"https://data.binance.vision/data/futures/um/%s/klines/%s/1m/%s-1m-%s.zip",
		i, m.symbol, m.symbol, d,
	)

	fp, err := downloadFile(url)
	if err != nil {
		return nil, err
	}
	defer os.Remove(fp)

	zrc, err := zip.OpenReader(fp)
	if err != nil {
		return nil, err
	}
	defer zrc.Close()

	irc, err := zrc.File[0].Open()
	if err != nil {
		return nil, err
	}
	defer irc.Close()

	cds, err := m.parseCSV(irc)
	if err != nil {
		return nil, err
	}

	return cds, nil
}

func downloadFile(url string) (filePath string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", err
	}

	tf, err := os.CreateTemp("", "download-*.zip")
	if err != nil {
		return "", err
	}

	_, err = io.Copy(tf, resp.Body)
	if err != nil {
		tf.Close()
		return "", err
	}
	tf.Close()

	return tf.Name(), nil
}

func (m *MarketData) parseCSV(rc io.ReadCloser) ([]*types.Candle, error) {
	r := csv.NewReader(rc)
	r.FieldsPerRecord = -1 // 레코드당 필드 수 제한 두지 않음

	var cds []*types.Candle

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		openTime, _ := strconv.ParseInt(record[0], 10, 64)
		open, _ := strconv.ParseFloat(record[1], 64)
		high, _ := strconv.ParseFloat(record[2], 64)
		low, _ := strconv.ParseFloat(record[3], 64)
		close, _ := strconv.ParseFloat(record[4], 64)
		volume, _ := strconv.ParseFloat(record[5], 64)
		closeTime, _ := strconv.ParseInt(record[6], 10, 64)
		quoteVolume, _ := strconv.ParseFloat(record[7], 64)
		count, _ := strconv.Atoi(record[8])
		takerBuyVolume, _ := strconv.ParseFloat(record[9], 64)
		takerBuyQuoteVolume, _ := strconv.ParseFloat(record[10], 64)

		cds = append(cds, &types.Candle{
			Symbol:              m.symbol,
			OpenTime:            openTime,
			Open:                open,
			High:                high,
			Low:                 low,
			Close:               close,
			Volume:              volume,
			CloseTime:           closeTime,
			QuoteVolume:         quoteVolume,
			Count:               count,
			TakerBuyVolume:      takerBuyVolume,
			TakerBuyQuoteVolume: takerBuyQuoteVolume,
			Closed:              true,
		})
	}

	return cds, nil
}

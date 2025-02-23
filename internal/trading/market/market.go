package market

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type MarketData struct {
	symbol string
}

func NewMarketData(symbol string) *MarketData {
	return &MarketData{symbol: symbol}
}

func (m *MarketData) DownloadFile(year, month, day int) (filePath string, err error) {
	url := fmt.Sprintf(
		"https://data.binance.vision/data/futures/um/daily/klines/%s/1m/%s-1m-%d-%02d-%02d.zip",
		m.symbol, m.symbol, year, month, day,
	)

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
	tfPath := tf.Name()

	_, err = io.Copy(tf, resp.Body)
	if err != nil {
		tf.Close()
		return "", err
	}
	tf.Close()
	defer os.Remove(tfPath)

	zrc, err := zip.OpenReader(tfPath)
	if err != nil {
		return "", err
	}
	defer zrc.Close()

	irc, err := zrc.File[0].Open()
	if err != nil {
		return "", err
	}
	defer irc.Close()

	sUrl := strings.Split(url, "/")
	fileName := strings.Split(sUrl[len(sUrl)-1], ".")[0] + ".csv"

	f, err := os.Create(fileName)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(f, irc)
	if err != nil {
		f.Close()
		return "", err
	}
	f.Close()

	return f.Name(), nil
}

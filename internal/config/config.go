package config

// TODO: 임시
type Config struct {
	Symbols          []string
	DBDriver         string
	DBSource         string
	BinanceApiKey    string
	BinanceSecretKey string
}

// TODO: 임시
func LoadConfig() (*Config, error) {
	return &Config{
		Symbols:          []string{"BTCUSDT", "ETHUSDT", "XRPUSDT"},
		DBDriver:         "mysql",
		DBSource:         "root:123456@tcp(172.17.0.2:3306)/binance?allowAllFiles=true",
		BinanceApiKey:    "",
		BinanceSecretKey: "",
	}, nil
}

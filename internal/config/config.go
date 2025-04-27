package config

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Host struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type Database struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	Password string `json:"password"`
	Dbname   string `json:"dbname"`
	Host
}

type ServiceConfig struct {
	Host
	Database Database `json:"database"`
}

type CommonConfig struct {
	Kafka Host `json:"kafka"`
}

type Config struct {
	Common  CommonConfig `json:"common"`
	Service struct {
		Market struct {
			ServiceConfig
			Symbols          []string `json:"symbols"`
			BinanceApiKey    string   `json:"binanceApiKey"`
			BinanceSecretKey string   `json:"binanceSecretKey"`
		} `json:"market"`
	} `json:"service"`
}

func LoadConfig() (*Config, error) {
	url := fmt.Sprintf("%s/config?service=market", os.Getenv("CONFIG_API_URL"))
	resp, err := http.Get(url)
	if err != nil {
		// 에러 처리
		return nil, err
	}
	defer resp.Body.Close()

	var config *Config
	err = json.NewDecoder(resp.Body).Decode(&config)
	if err != nil {
		// 에러 처리
		return nil, err
	}

	return config, nil
}

// 데이터베이스 추가 시, 해당 함수에서 관리
func (d *Database) DSN() string {
	switch d.Type {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?allowAllFiles=true",
			d.Username,
			d.Password,
			d.Host.Host,
			d.Host.Port,
			d.Dbname,
		)
	default:
		return ""
	}
}

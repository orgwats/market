package config

import (
	"flag"
	"os"

	"github.com/naoina/toml"
)

type Config struct {
	WebsocketURL string
}

var configPath = flag.String("config", "./config.toml", "Path to the config TOML file")

func LoadConfig() *Config {
	flag.Parse()

	c := new(Config)

	if file, err := os.Open(*configPath); err != nil {
		panic(err)
	} else if err = toml.NewDecoder(file).Decode(c); err != nil {
		panic(err)
	} else {
		return c
	}
}

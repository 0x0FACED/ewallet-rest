package server

import "ewallet/internal/app/store"

type Config struct {
	BindAddress string `toml:"bind_address"`
	LogLevel    string `toml:"log_level"`
	StoreCfg    *store.Config
}

func NewConfig() *Config {
	return &Config{
		StoreCfg: store.NewConfig(),
	}
}

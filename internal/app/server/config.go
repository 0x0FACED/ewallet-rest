package server

type Config struct {
	BindAddress string `toml:"bind_address"`
	DatabaseUrl string `toml:"database_url"`
}

func NewConfig() *Config {
	return &Config{}
}

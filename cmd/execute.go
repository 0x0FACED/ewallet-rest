package cmd

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/MisterGnida/ewallet-rest/internal/app/server"
	"log"
)

var (
	configPath string
)

func init() {
	flag.StringVar(
		&configPath,
		"config-path",
		"config/server.toml",
		"path to config",
	)
}

func Execute() {
	flag.Parse()

	config := server.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}
	s := server.New(config)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}

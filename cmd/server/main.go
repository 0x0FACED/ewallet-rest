package server

import "flag"

var flagConfig = flag.String("config", "./configs/server.toml", "path to config file")

func main() {
	flag.Parse()

}

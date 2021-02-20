package main

import (
	"falcon-seed/internal/api"
	"falcon-seed/pkg/config"
	"flag"
)

func main() {
	cfgPath := flag.String("c", "./cmd/server/conf.local.yaml", "Path to config file")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	checkErr(err)

	checkErr(api.Start(cfg))
}

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

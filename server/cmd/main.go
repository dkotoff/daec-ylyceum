package main

import (
	"log"
	"os"

	"github.com/dkotoff/daec-ylyceum/server/app"
	"github.com/dkotoff/daec-ylyceum/server/config"
	"github.com/dkotoff/daec-ylyceum/server/logger"
)

func main() {
	conf, err := config.LoadFromEnv()
	if err != nil {
		logger.Error("Error to parse config")
		os.Exit(1)
	}

	app, err := app.New(conf)
	if err != nil {
		log.Fatal("failed to read config")
		os.Exit(1)
	}
	err = app.Run()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

}

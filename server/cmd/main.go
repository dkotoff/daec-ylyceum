package main

import (
	"log"
	"os"

	"github.com/dkotoff/daec-ylyceum/server/config"
	"github.com/dkotoff/daec-ylyceum/server/internal/app"
	"github.com/dkotoff/daec-ylyceum/server/logger"
)

func main() {

	logger.Info("Read config...")
	conf, err := config.LoadFromEnv()
	if err != nil {
		logger.Fatal("Failed to read config")
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

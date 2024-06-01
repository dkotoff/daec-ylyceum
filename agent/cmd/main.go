package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dkotoff/daec-ylyceum/agent/app"
	"github.com/dkotoff/daec-ylyceum/agent/config"
)

func main() {

	config, err := config.LoadFromEnv()
	if err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	app, err := app.New(config)
	if err != nil {
		log.Fatalf("Failed to init app: %v", err)
	}

	app.Run()
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	<-exit
	app.Stop()

}

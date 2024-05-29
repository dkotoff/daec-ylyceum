package main

import (
	"log"
	"os"

	"github.com/dkotoff/daec-ylyceum/server/app"
	"github.com/dkotoff/daec-ylyceum/server/config"
)

func main() {

	log.Print("Read config...")
	conf, err := config.LoadFromEnv()
	if err != nil {
		log.Fatal("Failed to read config")
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

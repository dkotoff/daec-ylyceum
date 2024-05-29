package main

import (
	"log"
	"os"

	"github.com/dkotoff/daec-ylyceum/server/internal/app"
)

func main() {

	app, err := app.New()
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

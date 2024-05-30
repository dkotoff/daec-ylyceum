package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/dkotoff/daec-ylyceum/agent/app"
)

func main() {
	app := new(app.App)

	app.Run()
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	<-exit
	app.Stop()

}

package app

import (
	"github.com/dkotoff/daec-ylyceum/agent/config"
	computingservice "github.com/dkotoff/daec-ylyceum/agent/internal/computing_service"
	"github.com/dkotoff/daec-ylyceum/agent/logger"
)

type App struct {
	service *computingservice.ExpressionService
	config  *config.Config
}

func New(cfg *config.Config) (*App, error) {
	return &App{config: cfg}, nil
}

func (a *App) Run() {
	a.service = computingservice.NewExpressionService(a.config)
	logger.Info("Run agent with computing_power=%d", a.config.ComputingPower)
	a.service.Run()

}

func (a *App) Stop() {
	a.service.Stop()
	logger.Info("Stop goroutines")
}

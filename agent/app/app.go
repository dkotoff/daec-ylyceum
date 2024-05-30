package app

import computingservice "github.com/dkotoff/daec-ylyceum/agent/internal/computing_service"

type App struct {
	service *computingservice.ExpressionService
}

func (a *App) Run() {
	a.service = computingservice.NewExpressionService()
	a.service.Run(4)

}

func (a *App) Stop() {
	a.service.Stop()
}

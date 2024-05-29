package app

import (
	"net/http"

	expressionservice "github.com/dkotoff/daec-ylyceum/server/internal/expression_service"
	"github.com/go-chi/chi/v5"
)

type App struct {
	server *http.Server
}

func New() (*App, error) {

	app := new(App)

	service := expressionservice.NewExpressionService()

	router := chi.NewRouter()

	router.Route("/api/v1", func(r chi.Router) {
		r.HandleFunc("/expressions", service.ExpressionsHandler)
		r.HandleFunc("/calculate", service.CalculateHandler)
		r.HandleFunc("/expression/{id}", service.ExpressionHandler)

	})

	router.Route("/internal", func(r chi.Router) {
		r.Get("/task", service.GetTask)
		r.Post("/task", service.PostTask)
	})

	app.server = &http.Server{
		Handler: router,
		Addr:    ":5002",
	}

	return app, nil
}

func (a *App) Run() error {
	err := a.server.ListenAndServe()
	return err
}

package app

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/dkotoff/daec-ylyceum/server/config"
	expressionservice "github.com/dkotoff/daec-ylyceum/server/internal/expression_service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
)

type App struct {
	server *http.Server
}

func New(conf *config.Config) (*App, error) {

	app := new(App)

	service := expressionservice.NewExpressionService(conf)
	router := chi.NewRouter()
	logger := httplog.NewLogger("ExpressionOrchestrator", httplog.Options{
		LogLevel:       slog.LevelWarn,
		Concise:        false,
		RequestHeaders: false,
	})

	router.Use(httplog.RequestLogger(logger))

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
	log.Print("Start listen and serve")
	return err
}

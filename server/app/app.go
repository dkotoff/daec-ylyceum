package app

import (
	"html/template"
	"net/http"
	"os"

	"github.com/dkotoff/daec-ylyceum/server/config"
	expressionservice "github.com/dkotoff/daec-ylyceum/server/internal/expression_service"
	"github.com/dkotoff/daec-ylyceum/server/internal/web"
	"github.com/dkotoff/daec-ylyceum/server/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type App struct {
	server *http.Server
}

func New(conf *config.Config) (*App, error) {
	app := new(App)

	service := expressionservice.NewExpressionService(conf)
	router := chi.NewRouter()
	router.Use(middleware.SetHeader("Content-Type:", "application/json"))

	if _, err := os.Stat("./templates"); !os.IsNotExist(err) {
		templates := make(map[string]*template.Template)
		templates["index.html"] = template.Must(template.ParseFiles("templates/index.html"))
		templates["expression.html"] = template.Must(template.ParseFiles("templates/expression.html"))
		webApp := web.NewWebApp(templates, service)
		router.HandleFunc("/", webApp.IndexHandler)
		router.Post("/submit-expression", webApp.SubmitExpression)
		router.Get("/get-expression/{id}", webApp.GetExpressionById)
	} else {
		logger.Error("Templates folder not found Web Application not running")
	}

	router.Route("/api/v1", func(r chi.Router) {
		r.HandleFunc("/expressions", service.ExpressionsHandler)
		r.HandleFunc("/calculate", service.CalculateHandler)
		r.HandleFunc("/expressions/{id}", service.ExpressionHandler)

	})

	router.Route("/internal", func(r chi.Router) {
		r.Get("/task", service.GetTaskHandler)
		r.Post("/task", service.PostTaskHandler)
	})

	app.server = &http.Server{
		Handler: router,
		Addr:    ":" + conf.ServerPort,
	}

	return app, nil
}

func (a *App) Run() error {
	logger.Info("Start listen and serve at http://localhost%v", a.server.Addr)
	err := a.server.ListenAndServe()
	return err
}

package web

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"

	expressionservice "github.com/dkotoff/daec-ylyceum/server/internal/expression_service"
	"github.com/go-chi/chi/v5"
)

type WebApp struct {
	templates    map[string]*template.Template
	expr_service *expressionservice.ExpressionsService
}

func NewWebApp(templates map[string]*template.Template,
	expr_service *expressionservice.ExpressionsService) *WebApp {
	return &WebApp{
		templates:    templates,
		expr_service: expr_service,
	}
}

func (wa *WebApp) IndexHandler(w http.ResponseWriter, r *http.Request) {
	json, err := json.Marshal(wa.expr_service.GetExperssions())
	if err != nil {
		http.Error(w, "Failed to parse expressions", 500)
	}
	if json == nil {
		json = make([]byte, 0)
	}

	templ := wa.templates["index.html"]
	templ.ExecuteTemplate(w,
		"index.html",
		&map[string]template.JS{"Expressions": template.JS(json)})

}

func (wa *WebApp) SubmitExpression(w http.ResponseWriter, r *http.Request) {
	expr := r.PostFormValue("expression")
	expressionId, err := wa.expr_service.AddExpression(expr)

	if err != nil {
		log.Fatal(err)
	}
	expression, ok := wa.expr_service.GetExpression(expressionId)
	if !ok {
		log.Fatal("not found")
	}
	tmpl := wa.templates["expression.html"]
	tmpl.ExecuteTemplate(w, "expression.html", expression)
}

func (wa *WebApp) GetExpressionById(w http.ResponseWriter, r *http.Request) {
	expression_id_string := chi.URLParam(r, "id")
	expression_id, err := strconv.Atoi(expression_id_string)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	expression, ok := wa.expr_service.GetExpression(expression_id)
	if !ok {
		return
	}

	tmpl := wa.templates["expression.html"]
	tmpl.ExecuteTemplate(w, "expression.html", expression)
}

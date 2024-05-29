package expressionservice

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (s *ExpressionsService) CalculateHandler(w http.ResponseWriter, r *http.Request) {

	buff, _ := io.ReadAll(r.Body)
	var request map[string]string
	err := json.Unmarshal(buff, &request)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}

	id, err := s.AddExpression(request["expression"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}

	resp_schema := CalculateResponseSchema{Id: id}
	resp_map := map[string]CalculateResponseSchema{"expression": resp_schema}

	response, err := json.Marshal(&resp_map)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}
	w.Write(response)
}

func (s *ExpressionsService) ExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	var result []ExpressionResponseSchema
	for _, expression := range s.expressions {
		schema := ExpressionResponseSchema{
			Id:     expression.id,
			Status: s.tasks[expression.head_task_id].status,
			Result: s.tasks[expression.head_task_id].result,
		}
		result = append(result, schema)

	}
	result_byte, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write(result_byte)
}

func (s *ExpressionsService) ExpressionHandler(w http.ResponseWriter, r *http.Request) {

	expression_id_string := chi.URLParam(r, "id")
	expression_id, err := strconv.Atoi(expression_id_string)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	expression, ok := s.expressions[expression_id]
	if !ok {
		http.Error(w, "Expression not found", http.StatusNotFound)
		return
	}

	answer := make(map[string]ExpressionResponseSchema)

	answer["expression"] = ExpressionResponseSchema{
		Id:     expression_id,
		Status: s.tasks[expression.head_task_id].status,
		Result: s.tasks[expression.head_task_id].result,
	}

	answer_buff, err := json.Marshal(answer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(answer_buff)
}

func (s *ExpressionsService) GetTask(w http.ResponseWriter, r *http.Request) {
	task, ok := s.GetUnfinishedTask()
	if !ok {
		http.Error(w, "No tasks", http.StatusNotFound)
		return
	}

	out_buff, err := json.Marshal(task)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
	}

	w.Write(out_buff)
}

func (s *ExpressionsService) PostTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	buff, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}

	var request TaskSchemaRequest

	json.Unmarshal(buff, &request)
	if err != nil {
		http.Error(w, http.StatusText(422), http.StatusUnprocessableEntity)
	}

	ok := s.SetTaskResult(request.Id, request.Result)
	if !ok {
		http.Error(w, http.StatusText(404), http.StatusNotFound)
	}
}

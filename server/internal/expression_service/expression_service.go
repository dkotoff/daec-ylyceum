package expressionservice

import (
	"errors"

	"github.com/dengsgo/math-engine/engine"
	"github.com/dkotoff/daec-ylyceum/server/config"
	"github.com/dkotoff/daec-ylyceum/server/logger"
)

type ExpressionsService struct {
	expressions map[int]*Expression
	tasks       map[int]*Task
	config      *config.Config
}

func NewExpressionService(conf *config.Config) *ExpressionsService {
	return &ExpressionsService{
		expressions: make(map[int]*Expression),
		tasks:       make(map[int]*Task),
		config:      conf,
	}
}

func (s *ExpressionsService) AddExpression(expr string) (int, error) {

	// Создает новое выражение и разбирает на задачи

	if len(expr) <= 0 {
		return 0, errors.New("Expression length <= 0")
	}
	tokens, err := engine.Parse(expr)

	if err != nil {
		return 0, err
	}

	// Создание дерева выраженийы

	ast := engine.NewAST(tokens, expr)

	if ast.Err != nil {
		return 0, ast.Err
	}

	// Разбор дерева на задачи

	tree := ast.ParseExpression()
	if ast.Err != nil {
		return 0, ast.Err
	}

	head := s.CreateTasksFromTree(tree)

	expression := NewExpression()
	expression.head_task_id = head
	expression.expr = expr

	s.expressions[expression.id] = expression

	logger.Debug("Added new expression: %s", expr)
	return expression.id, nil

}

func (s *ExpressionsService) CreateTasksFromTree(node engine.ExprAST) int {

	// Если узел дерева это число - просто создаем заврешенную задачу и выходим

	if _, ok := node.(engine.NumberExprAST); ok {
		task := NewTask()
		task.status = Complete
		task.result = float64(node.(engine.NumberExprAST).Val)
		s.tasks[task.id] = task
		return task.id
	}

	// Иначе если узел операция создает задачу и продолжает рекурсивно проходить по дереву
	task := NewTask()
	task.left = s.CreateTasksFromTree(node.(engine.BinaryExprAST).Lhs)
	task.right = s.CreateTasksFromTree(node.(engine.BinaryExprAST).Rhs)
	task.operation = node.(engine.BinaryExprAST).Op
	s.tasks[task.id] = task
	logger.Debug("Added new task LeftId:%d RightId:%d Operation:%s", task.left, task.right, task.operation)

	return task.id
}

func (s *ExpressionsService) SetTaskResult(id int, result float64) bool {

	_, ok := s.tasks[id]
	if !ok {
		return false
	}

	logger.Debug("Task %d calculated result %f", id, result)

	s.tasks[id].result = result
	s.tasks[id].status = Complete

	return true
}

func (s *ExpressionsService) GetUnfinishedTask() (TaskSchema, bool) {
	for _, task := range s.tasks {
		if task.status == Complete || task.status == InProgress {
			continue
		}
		left, _ := s.tasks[task.left]
		if left.status != Complete {
			continue
		}

		right, _ := s.tasks[task.right]
		if right.status != Complete {
			continue
		}

		var op_time int
		switch task.operation {
		case "+":
			op_time = s.config.TimeAddition
		case "-":
			op_time = s.config.TimeSubtraction
		case "/":
			op_time = s.config.TimeDivision
		case "*":
			op_time = s.config.TimeMultiplication
		}
		task.status = InProgress
		return TaskSchema{
			Id:             task.id,
			Arg1:           left.result,
			Arg2:           right.result,
			Operation:      task.operation,
			Operation_time: op_time,
		}, true
	}
	return TaskSchema{}, false
}

func (s *ExpressionsService) GetExpression(id int) (ExpressionResponseSchema, bool) {
	expression, ok := s.expressions[id]

	if !ok {
		return ExpressionResponseSchema{}, false
	}

	return ExpressionResponseSchema{
		Id:     expression.id,
		Status: s.tasks[expression.head_task_id].status,
		Result: s.tasks[expression.head_task_id].result,
	}, true
}

func (s *ExpressionsService) GetExperssions() []ExpressionResponseSchema {
	var result []ExpressionResponseSchema
	for _, expression := range s.expressions {
		schema := ExpressionResponseSchema{
			Id:     expression.id,
			Status: s.tasks[expression.head_task_id].status,
			Result: s.tasks[expression.head_task_id].result,
		}
		result = append(result, schema)

	}
	return result
}

func (s *ExpressionsService) GetExpressionResult(id int) (float64, bool) {

	expression, ok := s.expressions[id]
	if !ok {
		return 0, false
	}

	// if !s.tasks[expression.head_task_id].status {
	// 	return 0, false
	// }

	// Очищаем пул задач для выражения если оно решено и записываем ответ в выражение
	if task, _ := s.tasks[expression.head_task_id]; task.status == Complete {
		expression.result = task.result
		expression.status = true
		s.DeleteTasksRecursive(task.id)
	}

	return expression.result, expression.status
}

func (s *ExpressionsService) DeleteTasksRecursive(id int) {
	task, ok := s.tasks[id]
	if !ok {
		return
	}

	// Если идентификатор задачи равен -1 (что свойственно для задач из 1 числа) просто удаляем и выходим
	if task.left == -1 {
		delete(s.tasks, id)
	}

	s.DeleteTasksRecursive(task.left)
	s.DeleteTasksRecursive(task.right)
	delete(s.tasks, id)
}

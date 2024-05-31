package expressionservice

import (
	"errors"

	"github.com/dengsgo/math-engine/engine"
	"github.com/dkotoff/daec-ylyceum/server/config"
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

	if len(expr) <= 0 {
		return 0, errors.New("Expression length <= 0")
	}
	tokens, err := engine.Parse(expr)

	if err != nil {
		return 0, err
	}

	ast := engine.NewAST(tokens, expr)

	if ast.Err != nil {
		return 0, ast.Err
	}

	tree := ast.ParseExpression()
	if ast.Err != nil {
		return 0, ast.Err
	}

	head := s.CreateTasksFromTree(tree)

	expression := NewExpression()
	expression.head_task_id = head
	expression.expr = expr

	s.expressions[expression.id] = expression
	return expression.id, nil

}

func (s *ExpressionsService) CreateTasksFromTree(node engine.ExprAST) int {

	if _, ok := node.(engine.NumberExprAST); ok {
		task := NewTask()
		task.status = true
		task.result = int(node.(engine.NumberExprAST).Val)
		s.tasks[task.id] = task
		return task.id
	}

	task := NewTask()
	task.left = s.CreateTasksFromTree(node.(engine.BinaryExprAST).Lhs)
	task.right = s.CreateTasksFromTree(node.(engine.BinaryExprAST).Rhs)
	task.operation = node.(engine.BinaryExprAST).Op
	s.tasks[task.id] = task

	return task.id
}

func (s *ExpressionsService) SetTaskResult(id int, result int) bool {

	_, ok := s.tasks[id]
	if !ok {
		return false
	}

	s.tasks[id].result = result
	s.tasks[id].status = true

	return true
}

func (s *ExpressionsService) GetUnfinishedTask() (TaskSchema, bool) {
	for _, task := range s.tasks {
		if task.status == true {
			continue
		}
		left, _ := s.tasks[task.left]
		if left.status != true {
			continue
		}

		right, _ := s.tasks[task.right]
		if right.status != true {
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

func (s *ExpressionsService) GetExpressionResult(id int) (int, bool) {

	expression, ok := s.expressions[id]
	if !ok {
		return 0, false
	}

	// if !s.tasks[expression.head_task_id].status {
	// 	return 0, false
	// }

	// Очищаем пул задач для выражения если оно решено и записываем ответ в выражение
	if task, _ := s.tasks[expression.head_task_id]; task.status {
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

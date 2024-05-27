package service

import (
	"github.com/dengsgo/math-engine/engine"
)

type ExpressionsService struct {
	expressions map[int]*Expression
	tasks       map[int]*Task
}

func (s *ExpressionsService) AddExpression(id int, expr string) error {
	tokens, err := engine.Parse(expr)

	if err != nil {
		return err
	}

	ast := engine.NewAST(tokens, expr)

	if ast.Err != nil {
		return ast.Err
	}

	tree := ast.ParseExpression()
	if ast.Err != nil {
		return ast.Err
	}

	head := s.CreateTasksFromTree(tree)

	expression := NewExpression(id)
	expression.head_task_id = head
	expression.expr = expr

	s.expressions[expression.id] = expression
	return nil

}

func (s *ExpressionsService) CreateTasksFromTree(node engine.ExprAST) int {

	if _, ok := node.(engine.NumberExprAST); ok {
		task := NewTask()
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

func (s *ExpressionsService) SetTaskResult(id int, result int) {
	s.tasks[id].result = result
	s.tasks[id].status = true
}

func (s *ExpressionsService) GetUnfinishedTasks() []Task {

	result := make([]Task, 0)

	for _, task := range s.tasks {
		if !task.status {
			result = append(result, *task)
		}
	}

	return result
}

func (s *ExpressionsService) GetExpressionResult(id int) (int, bool) {

	expression, ok := s.expressions[id]
	if !ok {
		return 0, false
	}

	if !s.tasks[expression.head_task_id].status {
		return 0, false
	}

	return s.tasks[expression.head_task_id].result, true
}

package service

type Expression struct {
	id           int
	head_task_id int
	expr         string
}

func NewExpression(id int) *Expression {
	return &Expression{
		id: id,
	}
}

package expressionservice

var expressionCounter = 0

type Expression struct {
	id           int
	head_task_id int
	expr         string
	status       bool
	result       int
}

func NewExpression() *Expression {
	defer func() {
		expressionCounter++
	}()
	return &Expression{
		id:     expressionCounter,
		status: false,
		result: 0,
	}
}

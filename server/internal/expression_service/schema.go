package expressionservice

type CalculateResponseSchema struct {
	Id int `json:"id"`
}

type ExpressionResponseSchema struct {
	Id     int     `json:"id"`
	Status bool    `json:"status"`
	Result float64 `json:"result"`
}

type TaskSchema struct {
	Id             int     `json:"id"`
	Arg1           float64 `json:"arg1"`
	Arg2           float64 `json:"arg2"`
	Operation      string  `json:"operation"`
	Operation_time int     `json:"operation_time"`
}

type TaskSchemaRequest struct {
	Id     int     `json:"id"`
	Result float64 `json:"result"`
}

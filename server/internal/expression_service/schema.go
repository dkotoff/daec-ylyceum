package expressionservice

type CalculateResponseSchema struct {
	Id int `json:"id"`
}

type ExpressionResponseSchema struct {
	Id     int  `json:"id"`
	Status bool `json:"status"`
	Result int  `json:"result"`
}

type TaskSchema struct {
	Id             int    `json:"id"`
	Arg1           int    `json:"arg1"`
	Arg2           int    `json:"arg2"`
	Operation      string `json:"operation"`
	Operation_time int    `json:"operation_time"`
}

type TaskSchemaRequest struct {
	Id     int `json:"id"`
	Result int `json:"result"`
}

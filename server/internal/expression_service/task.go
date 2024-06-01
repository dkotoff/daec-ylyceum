package expressionservice

var id_counter int = 0

type Task struct {
	id        int
	operation string
	left      int
	right     int
	result    float64
	status    bool
}

func NewTask() *Task {
	id_counter++
	return &Task{id: id_counter,
		left: -1, operation: "+", right: -1, result: 0, status: false}
}

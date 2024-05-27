package service

var id_counter int = 1

type Task struct {
	id        int
	operation string
	left      int
	right     int
	result    int
	status    bool
}

func NewTask() *Task {
	defer func() {
		id_counter++
	}()
	return &Task{id: id_counter}
}

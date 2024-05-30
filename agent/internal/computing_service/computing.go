package computingservice

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type TaskSchemaResponse struct {
	Id     int `json:"id"`
	Result int `json:"result"`
}

type TaskSchema struct {
	Id             int    `json:"id"`
	Arg1           int    `json:"arg1"`
	Arg2           int    `json:"arg2"`
	Operation      string `json:"operation"`
	Operation_time int    `json:"operation_time"`
}

var ops = map[string]func(int, int) int{
	"+": func(a, b int) int { return a + b },
	"-": func(a, b int) int { return a - b },
	"*": func(a, b int) int { return a * b },
	"/": func(a, b int) int { return a / b },
}

type ExpressionService struct {
	input_chan  chan Task
	output_chan chan Task
	client      *http.Client
	stop_chan   chan interface{}
}

func NewExpressionService() *ExpressionService {
	return &ExpressionService{
		input_chan:  make(chan Task),
		output_chan: make(chan Task),
		client:      &http.Client{},
		stop_chan:   make(chan interface{}),
	}
}

func (s *ExpressionService) RunComputers(count int, stop chan interface{}) {

	for i := 0; i <= count; i++ {
		go func() {
			for true {
				select {
				case <-stop:
					return
				case task := <-s.input_chan:
					time.Sleep(time.Millisecond * time.Duration(task.operation_time))
					task.result = ops[task.operation](task.arg1, task.arg2)
					s.output_chan <- task
				}
			}
		}()
	}
}

func (s *ExpressionService) RunRequestsLoop(stop chan interface{}) {
	go func() {
		for true {
			time.Sleep(time.Duration(1) * time.Second)
			select {
			case <-stop:
				return
			default:
				resp, err := s.client.Get("http://localhost:5002/internal/task")
				if err != nil {
					log.Printf("Error at request to orcestrator: %v", err)
					continue
				}
				if resp.StatusCode != 200 {
					log.Printf("No tasks")
					continue
				}
				defer resp.Body.Close()
				buff, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Printf("Error at request to orcestrator: %v", err)
					continue
				}

				var taskSchema TaskSchema
				err = json.Unmarshal(buff, &taskSchema)
				if err != nil {
					log.Printf("Error to validate taks: %v", err)
				}
				s.input_chan <- Task{id: taskSchema.Id,
					arg1:           taskSchema.Arg1,
					arg2:           taskSchema.Arg2,
					operation_time: taskSchema.Operation_time,
					operation:      taskSchema.Operation}

			}
		}
	}()
}

func (s *ExpressionService) RunResponsesLoop(stop chan interface{}) {
	go func() {
		for true {
			select {
			case <-stop:
				return
			case task := <-s.output_chan:
				response_schema := TaskSchemaResponse{Id: task.id, Result: task.result}

				buff, err := json.Marshal(response_schema)
				if err != nil {
					log.Printf("Error to validate data: %v", err)
				}
				r := bytes.NewReader(buff)
				_, err = s.client.Post("http://localhost:5002/internal/task", "application/json", r)
				if err != nil {
					log.Printf("Error to validate data: %v", err)
				}
			}
		}
	}()
}

func (s *ExpressionService) Run(computing_power int) {
	log.Print("Starting service")
	s.RunComputers(computing_power, s.stop_chan)
	s.RunRequestsLoop(s.stop_chan)
	s.RunResponsesLoop(s.stop_chan)
}

func (s *ExpressionService) Stop() {
	s.stop_chan <- 0
}

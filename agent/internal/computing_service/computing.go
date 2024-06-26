package computingservice

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/dkotoff/daec-ylyceum/agent/config"
	"github.com/dkotoff/daec-ylyceum/agent/logger"
)

type TaskSchemaResponse struct {
	Id     int     `json:"id"`
	Result float64 `json:"result"`
}

type TaskSchema struct {
	Id             int     `json:"id"`
	Arg1           float64 `json:"arg1"`
	Arg2           float64 `json:"arg2"`
	Operation      string  `json:"operation"`
	Operation_time int     `json:"operation_time"`
}

var ops = map[string]func(float64, float64) float64{
	"+": func(a, b float64) float64 { return a + b },
	"-": func(a, b float64) float64 { return a - b },
	"*": func(a, b float64) float64 { return a * b },
	"/": func(a, b float64) float64 { return a / b },
}

type ExpressionService struct {
	input_chan  chan Task
	output_chan chan Task
	client      *http.Client
	stop_chan   chan interface{}
	config      *config.Config
}

func NewExpressionService(cfg *config.Config) *ExpressionService {
	return &ExpressionService{
		input_chan:  make(chan Task),
		output_chan: make(chan Task),
		client:      &http.Client{},
		stop_chan:   make(chan interface{}),
		config:      cfg,
	}
}

func (s *ExpressionService) RunComputers(count int, stop chan interface{}) {
	for i := 0; i < count; i++ {
		go worker(i, stop, s.input_chan, s.output_chan)
	}
}

func worker(num int, stop chan interface{}, input chan Task, output chan Task) {
	for {
		select {
		case task, ok := <-input:
			if !ok {
				return
			}
			time.Sleep(time.Millisecond * time.Duration(task.operation_time))
			if task.operation == "/" && task.arg2 == 0 {
				task.result = 0
				output <- task
				logger.Debug("Zero division at task %d calculated with result %f", task.id, task.result)
			} else {
				task.result = ops[task.operation](task.arg1, task.arg2)
				logger.Debug("Task %d calculated with result %f at goroutine %d", task.id, task.result, num)
				output <- task
			}
		case <-stop:
			return
		}
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
				resp, err := s.client.Get("http://localhost:" + s.config.ServerPort + "/internal/task")
				if err != nil {
					logger.Error("Failed to connect orcestrator: %v", err)
					continue
				}
				if resp.StatusCode == 404 {
					continue
				}

				defer resp.Body.Close()
				buff, err := io.ReadAll(resp.Body)
				if err != nil {
					logger.Error("Error at request to orcestrator: %v", err)
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
				_, err = s.client.Post("http://localhost:"+s.config.ServerPort+"/internal/task", "application/json", r)
				if err != nil {
					log.Printf("Error to validate data: %v", err)
				}
			}
		}
	}()
}

func (s *ExpressionService) Run() {
	log.Print("Starting service")
	s.RunComputers(s.config.ComputingPower, s.stop_chan)
	s.RunRequestsLoop(s.stop_chan)
	s.RunResponsesLoop(s.stop_chan)
}

func (s *ExpressionService) Stop() {
	s.stop_chan <- 0
}

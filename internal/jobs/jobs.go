package jobs

import (
	"fmt"
)

type Handler func(data map[string]string) (string, error)

type Job struct {
	Id      string
	Handler Handler
	Data    map[string]string
}

type JobResult struct {
	Id       string
	Message  string
	WorkerId int
	Success  bool
}

func worker(id int, jobs <-chan Job, results chan<- JobResult) {
	for j := range jobs {
		fmt.Println("worker", id, "started  job", j.Id)

		if j.Handler == nil {
			results <- JobResult{
				Message:  "no handler func passed",
				Success:  false,
				Id:       j.Id,
				WorkerId: id,
			}
			return
		}

		msg, err := j.Handler(j.Data)
		if err != nil {
			results <- JobResult{
				Message:  err.Error(),
				Success:  false,
				Id:       j.Id,
				WorkerId: id,
			}
			return
		}

		results <- JobResult{
			Message:  msg,
			Id:       j.Id,
			WorkerId: id,
			Success:  true,
		}
	}
}

func Start(workers int) (chan Job, chan JobResult) {
	const numJobs = 5
	jobs := make(chan Job, numJobs)
	results := make(chan JobResult, numJobs)

	// start workers
	for w := 1; w <= workers; w++ {
		go worker(w, jobs, results)
	}

	return jobs, results
}

func Stop(jobs chan Job, results chan JobResult) {
	close(jobs)
	close(results)
}

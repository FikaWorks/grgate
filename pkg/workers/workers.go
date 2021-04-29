package workers

import (
	"fmt"
)

// Worker process jobs from a queue
type Worker struct {
	ID     int
	Job    chan Job
	Queue  chan chan Job
	Cancel chan struct{}
}

// NewWorker return a worker which process jobs from a queue
func NewWorker(id int, queue chan chan Job, cancel chan struct{}) *Worker {
	return &Worker{
		ID:     id,
		Job:    make(chan Job),
		Queue:  queue,
		Cancel: cancel,
	}
}

// Start worker
func (w *Worker) Start() {
	go func() {
		for {
			w.Queue <- w.Job

			select {
			case work := <-w.Job:
				fmt.Printf("[worker #%d] received work request %s\n", w.ID, work.Repository)
				work.Process()
				fmt.Printf("[worker #%d] finished work for %s\n", w.ID, work.Repository)
			case <-w.Cancel:
				fmt.Printf("[worker #%d]: stopping\n", w.ID)
				return
			}
		}
	}()
}

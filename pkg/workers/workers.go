package workers

import (
	"github.com/rs/zerolog/log"
)

// Worker process jobs from a queue
type Worker struct {
	ID     int
	Job    chan *Job
	Queue  chan chan *Job
	Cancel chan struct{}
}

// NewWorker return a worker which process jobs from a queue
func NewWorker(id int, queue chan chan *Job, cancel chan struct{}) *Worker {
	return &Worker{
		ID:     id,
		Job:    make(chan *Job),
		Queue:  queue,
		Cancel: cancel,
	}
}

// Start worke, process job from the queue when they arrive
func (w *Worker) Start() {
	go func() {
		for {
			w.Queue <- w.Job

			select {
			case work := <-w.Job:
        log.Debug().
          Int("worker", w.ID).
          Str("owner", work.Owner).
          Str("repository", work.Repository).
          Msg("Processing work item from queue")
				work.Process()
        log.Debug().
          Int("worker", w.ID).
          Str("owner", work.Owner).
          Str("repository", work.Repository).
          Msg("Completed work from queue")
			case <-w.Cancel:
        log.Info().
          Int("worker", w.ID).
          Msg("Stopping worker queue")
				return
			}
		}
	}()
}

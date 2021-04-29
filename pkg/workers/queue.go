package workers

import (
	"github.com/rs/zerolog/log"
)

// WorkerPool dispatch jobs to a pool of workers
type WorkerPool struct {

  // Job queue
  JobQueue chan *Job

  // WorkerQueue is the job queue of a worker
  WorkerQueue chan chan *Job
  Workers []*Worker
}

// NewWorkerPool return a WorkerPool to process jobs
func NewWorkerPool(workerCount int, cancel chan struct{}) *WorkerPool {
  workers := []*Worker{}
  workerQueue := make(chan chan *Job, workerCount)

	for i := 0; i < workerCount; i++ {
    log.Info().Msgf("Initialising worker %d", i+1)
		worker := NewWorker(i+1, workerQueue, cancel)
    workers = append(workers, worker)
	}

  return &WorkerPool{
    JobQueue: make(chan *Job, 100),
	  WorkerQueue: workerQueue,
    Workers: workers,
  }
}

// Start the worker pool
func (wp *WorkerPool) Start() {
  for _, worker := range wp.Workers {
    log.Info().
      Int("worker", worker.ID).
      Msg("Starting worker")
    worker.Start()
  }

	go func() {
		for {
			select {
			case job := <-wp.JobQueue:
				go func() {
					worker := <-wp.WorkerQueue
					worker <- job
				}()
			}
		}
	}()
}

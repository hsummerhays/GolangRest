package worker

import (
	"context"
	"log/slog"
	"sync"
)

// Job represents a unit of work to be processed by the worker pool
type Job func(ctx context.Context) error

// Pool represents a worker pool
type Pool struct {
	jobs    chan Job
	wg      sync.WaitGroup
	workers int
}

// NewPool creates a new worker pool
func NewPool(workers int, bufferSize int) *Pool {
	return &Pool{
		jobs:    make(chan Job, bufferSize),
		workers: workers,
	}
}

// Start boots up the workers
func (p *Pool) Start(ctx context.Context) {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go func(workerID int) {
			defer p.wg.Done()
			slog.Info("Worker started", "worker_id", workerID)
			for {
				select {
				case <-ctx.Done():
					slog.Info("Worker shutting down due to context cancellation", "worker_id", workerID)
					return
				case job, ok := <-p.jobs:
					if !ok {
						slog.Info("Worker shutting down because job channel closed", "worker_id", workerID)
						return
					}
					if err := job(ctx); err != nil {
						slog.Error("Job failed", "worker_id", workerID, "error", err)
					}
				}
			}
		}(i)
	}
}

// Submit adds a job to the queue
func (p *Pool) Submit(job Job) {
	p.jobs <- job
}

// Stop closes the job channel and waits for all workers to finish processing
func (p *Pool) Stop() {
	close(p.jobs)
	p.wg.Wait()
}

package attack

import (
	"context"
	"sync"
	"time"

	"github.com/falconjonz/falcon_rdp/config"
	"github.com/falconjonz/falcon_rdp/logger"
)

type WorkItem struct {
	Target     config.Target
	Credential config.Credential
	RetryCount int
}

type WorkResult struct {
	Target      config.Target
	Credential  config.Credential
	Success     bool
	Error       string
	Duration    time.Duration
	Timestamp   time.Time
}

type WorkerPool struct {
	workers     int
	workChan    chan WorkItem
	resultsChan chan WorkResult
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
	logger      *logger.Logger
	cfg         *config.Config
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(workers int, cfg *config.Config, log *logger.Logger) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	return &WorkerPool{
		workers:     workers,
		workChan:    make(chan WorkItem, workers*2),
		resultsChan: make(chan WorkResult, workers*2),
		ctx:         ctx,
		cancel:      cancel,
		logger:      log,
		cfg:         cfg,
	}
}

// Start starts the worker pool
func (wp *WorkerPool) Start(handler func(WorkItem) WorkResult) {
	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go wp.worker(handler)
	}
}

// worker processes work items
func (wp *WorkerPool) worker(handler func(WorkItem) WorkResult) {
	defer wp.wg.Done()

	for {
		select {
		case <-wp.ctx.Done():
			return
		case work, ok := <-wp.workChan:
			if !ok {
				return
			}

			// Execute work with timeout
			ctx, cancel := context.WithTimeout(wp.ctx, wp.cfg.Attack.Timeout)
			defer cancel()

			start := time.Now()
			result := handler(work)
			result.Duration = time.Since(start)
			result.Timestamp = time.Now()

			select {
			case <-ctx.Done():
				return
			case wp.resultsChan <- result:
			}
		}
	}
}

// Submit submits a work item
func (wp *WorkerPool) Submit(item WorkItem) error {
	select {
	case <-wp.ctx.Done():
		return wp.ctx.Err()
	case wp.workChan <- item:
		return nil
	}
}

// Results returns the results channel
func (wp *WorkerPool) Results() <-chan WorkResult {
	return wp.resultsChan
}

// Stop stops the worker pool
func (wp *WorkerPool) Stop() {
	wp.cancel()
	close(wp.workChan)
}

// Wait waits for all workers to finish
func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
	close(wp.resultsChan)
}

// WaitFor waits for all workers with timeout
func (wp *WorkerPool) WaitFor(timeout time.Duration) error {
	done := make(chan struct{})
	go func() {
		wp.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-time.After(timeout):
		return context.DeadlineExceeded
	}
}

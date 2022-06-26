package controlflow

import (
	"context"
	"sync"

	"github.com/hashicorp/go-multierror"
)

type JobFunc func(int) error

type WorkerPool struct {
	Jobs, Workers int
}

func (pool WorkerPool) Execute(ctx context.Context, jobFunc JobFunc) error {
	return pool.dispatch(jobFunc).run(ctx)
}

func (pool *WorkerPool) dispatch(f JobFunc) *jobDispatcher {
	Assert(pool.Workers > 0 || pool.Jobs == 0, "invalid number of workers requested = %d", pool.Workers)

	dispatcher := jobDispatcher{
		f:        f,
		results:  make(chan error, pool.Jobs),
		jobsChan: make(chan int),
		jobs:     pool.Jobs,
	}
	dispatcher.waitGroup.Add(pool.Workers)

	for workerIdx := 0; workerIdx < pool.Workers; workerIdx++ {
		Go(dispatcher.worker)
	}
	return &dispatcher
}

func (dispatcher *jobDispatcher) run(ctx context.Context) error {
Loop:
	for jobIdx := 0; jobIdx < dispatcher.jobs; jobIdx++ {
		select {
		case <-ctx.Done():
			break Loop
		case dispatcher.jobsChan <- jobIdx:
		}
	}
	close(dispatcher.jobsChan)

	return multierror.Append(dispatcher.collect(), ctx.Err()).ErrorOrNil()
}

type jobDispatcher struct {
	f         JobFunc
	results   chan error
	jobsChan  chan int
	jobs      int
	waitGroup sync.WaitGroup
}

func (dispatcher *jobDispatcher) collect() *multierror.Error {
	// wait for workers to finish as they may still write pending results
	dispatcher.waitGroup.Wait()
	close(dispatcher.results)

	var mErr *multierror.Error
	for err := range dispatcher.results {
		mErr = multierror.Append(mErr, err)
	}
	return mErr
}

func (dispatcher *jobDispatcher) worker() {
	defer dispatcher.waitGroup.Done()
	for job := range dispatcher.jobsChan {
		dispatcher.results <- dispatcher.f(job)
	}
}

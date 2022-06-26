package worker

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	l "raftt.io/bananas/pkg/logging"
	sw "raftt.io/bananas/pkg/timeutil/stopwatch"
)

const (
	defaultSleepOnRetry = time.Minute * 2
	maxRetriesOnFailure = 2
)

type Worker struct {
	Job    func(context.Context) error
	Desc   string
	Period time.Duration
}

func (worker Worker) logger() logrus.FieldLogger {
	return l.Logger.WithField("worker", worker.Desc)
}

func (worker Worker) Loop(ctx context.Context) {
	for {
		worker.retryOnError(ctx, defaultSleepOnRetry)
		worker.logger().WithField("sleepingFor", worker.Period).Debug("Going to sleep")
		select {
		case <-time.After(worker.Period):
		case <-ctx.Done():
			worker.logger().WithError(ctx.Err()).Debug("stopping")
			return
		}
	}
}

func (worker *Worker) do(ctx context.Context) error {
	defer sw.Stopwatch(worker.Desc, logrus.Fields{"worker": worker.Desc})()
	return worker.Job(ctx)
}

func (worker *Worker) retryOnError(ctx context.Context, sleepOnRetry time.Duration) {
	for retryIdx := 0; retryIdx < maxRetriesOnFailure; retryIdx++ {
		if retryIdx > 0 {
			time.Sleep(sleepOnRetry)
			worker.logger().WithField("retry", retryIdx).Debug("Retrying")
		}
		if err := worker.do(ctx); err != nil {
			worker.logger().WithError(err).Error("Failed")
		} else {
			break
		}
	}
}

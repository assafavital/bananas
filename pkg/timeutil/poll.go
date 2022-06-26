package timeutil

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/pkg/errors"
)

// Retry retries op() with a constant interval.
// Convenience wrapper around backoff.Retry(). See docs for more info
func Retry(ctx context.Context, interval time.Duration, op func() error) error {
	return RetryNotify(ctx, interval, op, nil)
}

func RetryNotify(ctx context.Context, interval time.Duration, op func() error, notify backoff.Notify) error {
	err := backoff.RetryNotify(op, backoff.WithContext(backoff.NewConstantBackOff(interval), ctx), notify)
	if errors.Is(err, context.Canceled) {
		err = errors.WithStack(err)
	}
	return err
}

// NotFinished is an error used to indicate that an operation is still ongoing.
type NotFinished struct {
	err error
}

func (p *NotFinished) Unwrap() error {
	return p.err
}

func (p *NotFinished) Is(err error) bool {
	_, is := err.(*NotFinished)
	return is
}

func (p *NotFinished) Error() string {
	if p.err != nil {
		return fmt.Sprintf("timeout: %v", p.err)
	}
	return "timeout"
}

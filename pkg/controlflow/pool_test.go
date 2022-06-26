package controlflow

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestPool(t *testing.T) {
	type poolArgs struct {
		jobs, workers int
	}
	tests := []struct {
		name     string
		f        JobFunc
		pool     poolArgs
		wantErrs int
		timeout  time.Duration
	}{
		{
			name: "zero len",
			f: func(i int) error {
				panic("should not be called")
			},
			pool: poolArgs{0, 0},
		},
		{
			name: "one len",
			f: func(i int) error {
				return nil
			},
			pool: poolArgs{1, 1},
		},
		{
			name: "one err",
			f: func(i int) error {
				return errors.New("a")
			},
			pool:     poolArgs{1, 1},
			wantErrs: 1,
		},
		{
			name: "keep calling even with err",
			f: func(i int) error {
				return errors.New("a")
			},
			pool:     poolArgs{5, 5},
			wantErrs: 5,
		},
		{
			name: "only some errors",
			f: func(i int) error {
				if i < 2 {
					return errors.New("a")
				}
				return nil
			},
			pool:     poolArgs{5, 5},
			wantErrs: 2,
		},
		{
			name: "timeout returns",
			f: func(i int) error {
				time.Sleep(time.Second * 1)
				return nil
			},
			pool:     poolArgs{1, 1},
			wantErrs: 1,
			timeout:  time.Millisecond * 100,
		},
	}
	workerCounts := []int{1, 2, 5, 7, 23}
	for _, tt := range tests {
		for _, workers := range workerCounts {
			t.Run(fmt.Sprintf("%s-%d_workers", tt.name, workers), func(t *testing.T) {
				// Hack to allow the test to manipulate workers
				pool := WorkerPool{Jobs: tt.pool.jobs, Workers: workers}
				ctx := context.Background()
				var cancel func()
				if tt.timeout > 0 {
					ctx, cancel = context.WithTimeout(ctx, tt.timeout)
					defer cancel()
				}
				err := pool.Execute(ctx, tt.f)
				if (err != nil) != (tt.wantErrs > 0) {
					t.Errorf("Pool() error = %v, wantErrs %v", err, tt.wantErrs)
				}

				if tt.wantErrs > 0 {
					merr, ok := err.(*multierror.Error)
					if !ok {
						t.Errorf("wrong error type: %+v", merr)
					} else if len(merr.Errors) != tt.wantErrs {
						t.Errorf("Pool() errors = %d, wantErrs %v", len(merr.Errors), tt.wantErrs)
					}

				}
			})
		}
	}

	emptyFunc := func(int) error { return nil }

	assert.Panics(t, func() { // workers count is 0
		_ = WorkerPool{Workers: 0, Jobs: 10}.Execute(context.Background(), emptyFunc)
	})

	// But if the length is also 0, then 0 workers are fine
	PanicIfErr(WorkerPool{Workers: 0, Jobs: 0}.Execute(context.Background(), emptyFunc))

	// Use a canceled context
	cancelled, cancel := context.WithCancel(context.Background())
	cancel()
	err := WorkerPool{Workers: 1, Jobs: 1}.Execute(cancelled, emptyFunc)
	if err == nil {
		t.Errorf("Expected an error calling Execute with a canceled context")
	}

	// cancel during return
	cancelled, cancel = context.WithCancel(context.Background())
	err = WorkerPool{Workers: 1, Jobs: 2}.Execute(cancelled, func(int) error {
		cancel()
		return nil
	})
	if err == nil {
		t.Errorf("Expected error calling Execute when context is cancelled during run")
	}
}

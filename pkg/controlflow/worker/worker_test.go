package worker

import (
	"context"
	"testing"
	"time"

	"github.com/pkg/errors"
)

func TestWorker_retryOnError(t *testing.T) {
	tests := []struct {
		name             string
		failUntill, want int
	}{
		{
			name:       "Doesnt fail",
			failUntill: 0,
			want:       0,
		},
		{
			name:       "Always fail",
			failUntill: 1000,
			want:       maxRetriesOnFailure,
		},
		{
			name:       "Fails once",
			failUntill: 1,
			want:       1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			counter := 0
			job := &Worker{
				Job: func(_ context.Context) error {
					if counter >= tt.failUntill {
						return nil
					}
					counter++
					return errors.New("error")
				},
			}
			job.retryOnError(context.Background(), time.Millisecond)
			if counter != tt.want {
				t.Errorf("retryOnError() = %v, want %v", counter, tt.want)
			}
		})
	}
}

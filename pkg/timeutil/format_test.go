package timeutil

import (
	"testing"
	"time"
)

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name string
		d    time.Duration
		want string
	}{
		{"very quick", time.Microsecond, "<10ms"},
		{"sorta quick", 123 * time.Millisecond, "0.12s"},
		{"rounded", 1003 * time.Millisecond, "1s"},
		{"kinda slow", 4300 * time.Millisecond, "4.3s"},
		{"kinda slow rounded", 4321 * time.Millisecond, "4.3s"},
		{"very slow", 90 * time.Second, "1m30s"},
		{"very slow rounded", 90*time.Second + 153*time.Millisecond, "1m30.2s"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatDuration(tt.d); got != tt.want {
				t.Errorf("FormatDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

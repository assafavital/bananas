package timeutil

import (
	"fmt"
	"time"
)

func FormatDuration(d time.Duration) string {
	switch {
	case d < 10*time.Millisecond:
		return "<10ms"
	case d < time.Second:
		return fmt.Sprintf("%vs", d.Round(time.Millisecond*10).Seconds())
	default:
		return d.Round(time.Millisecond * 100).String()

	}
}

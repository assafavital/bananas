package controlflow

import (
	"testing"
)

func TestGo(t *testing.T) {
	tests := []struct {
		name   string
		f      func()
		panics bool
	}{
		{"nothing happens", func() {}, false},
		{"panics", func() { panic("TEST - should panic") }, true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			panicChan := make(chan struct{})
			noPanicChan := make(chan struct{})
			abort = func() {
				close(panicChan)
			}

			Go(func() {
				tt.f()
				close(noPanicChan)
			})
			var panicked bool
			select {
			case <-noPanicChan:
				panicked = false
			case <-panicChan:
				panicked = true
			}

			if panicked != tt.panics {
				t.Errorf("expected panic=%v, but panicked=%v", tt.panics, panicked)
			}
		})
	}
}

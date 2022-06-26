package controlflow

import (
	"runtime/debug"

	l "raftt.io/bananas/pkg/logging"
)

// defined as a variable for testing purposes
var abort = func() { l.Logger.Exit(1) }

func catchPanic() {
	if r := recover(); r != nil {
		l.Logger.WithField("recoveredFrom", r).Error("Panicked!")
		l.Logger.Debug(string(debug.Stack()))
		abort()
	}
}

func Go(f func()) {
	go func() {
		defer catchPanic()
		f()
	}()
}

func Execute(f func() error) error {
	defer catchPanic()
	return f()
}

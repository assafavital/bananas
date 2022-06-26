package controlflow

import (
	"fmt"

	l "raftt.io/bananas/pkg/logging"
)

// Assert checks a condition and panics with msg if it is false
func Assert(condition bool, format string, args ...interface{}) {
	if !condition {
		panic(fmt.Sprintf(format, args...))
	}
}

func PanicIfErr(err error) {
	if err != nil {
		l.Logger.WithError(err).Error("Panicking")
		panic(err)
	}
}

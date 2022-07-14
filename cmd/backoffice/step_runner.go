package main

import (
	"reflect"
	"runtime"

	"gorm.io/gorm"
	l "raftt.io/bananas/pkg/logging"
)

type step func(*gorm.DB) error

func run(step step) {
	if err := step(db); err != nil {
		l.Logger.Errorf("%q failed: [%v]", runtime.FuncForPC(reflect.ValueOf(step).Pointer()).Name(), err)
	}
}

package main

import "gorm.io/gorm"

type step func(*gorm.DB) error

func run(step step) error { return step(db) }

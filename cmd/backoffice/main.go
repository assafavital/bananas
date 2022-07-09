package main

import (
	"gorm.io/gorm"
	"raftt.io/bananas/cmd/backoffice/database/session"
	"raftt.io/bananas/cmd/backoffice/demo"
	l "raftt.io/bananas/pkg/logging"
)

var db *gorm.DB

func init() {
	dbSession, err := session.CurrentSession()
	if err != nil {
		panic(err)
	}
	db = dbSession.DB
}

func main() {
	l.Logger.Info("👋 WELCOME TO BACKOFFICE 👋")
	if err := run(demo.CreateUsers); err != nil {
		l.Logger.Error(err)
	}
	//if err := run(demo.MigrateUsers); err != nil {
	//	l.Logger.Error(err)
	//}
	if err := run(demo.BetterMigrateUsers); err != nil {
		l.Logger.Error(err)
	}
}

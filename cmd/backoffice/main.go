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
	run(demo.CreateUsers)
	run(demo.BetterMigrateUsers)
	l.Logger.Info("🐷 That's all Folks!")
}

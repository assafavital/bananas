package demo

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"raftt.io/bananas/cmd/backoffice/types"
	l "raftt.io/bananas/pkg/logging"
)

const (
	theOneAndOnly    = "bibi"
	theOnlyFemaleOne = "golda"
)

func createPMTable(db *gorm.DB) error {
	return db.
		Migrator().
		CreateTable(&types.PrimeMinister{})
}

func createPM(db *gorm.DB, name string) error {
	return errors.Wrap(
		db.
			Create(&types.PrimeMinister{Name: name}).
			Error,
		"failed to create primeMinister",
	)
}

func findPM(db *gorm.DB, name string) (*types.PrimeMinister, error) {
	var foundPM types.PrimeMinister
	if err := db.
		Where(&types.PrimeMinister{Name: name}).
		First(&foundPM).
		Error; err != nil {
		return nil, errors.WithStack(err)
	}
	return &foundPM, nil
}

func BabySteps(db *gorm.DB) error {
	l.Logger.Debug("🏃 Running 'Baby Steps' tutorial...")

	// 1. Define `prime_ministers` table
	if err := createPMTable(db); err != nil {
		return err
	}

	// 2. Create first PrimeMinister
	if err := createPM(db, theOneAndOnly); err != nil {
		return err
	}
	if err := createPM(db, theOnlyFemaleOne); err != nil {
		return err
	}

	// 3. Query a PrimeMinister
	found, err := findPM(db, theOneAndOnly)
	if err != nil {
		return err
	}
	l.Logger.WithField("foundPM", found).Info("found a prime minister")
	return nil
}

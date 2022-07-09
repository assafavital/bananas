package types

import "gorm.io/gorm"

type PrimeMinister struct {
	gorm.Model
	Name string
}

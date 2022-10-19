package model

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	IDTransaction string
	IDCustomer    uint
	Currency      string
	Amount        uint
	Status        string
}

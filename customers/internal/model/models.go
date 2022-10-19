package model

import "gorm.io/gorm"

type (
	Customer struct {
		gorm.Model
		Balance       uint
		Version       uint
		Email         string
		IDTransaction string
		Status        string
	}
	ProcessedTransaction struct {
		gorm.Model
		IDTransaction string
	}
	CustomerS struct {
		gorm.Model
		Balance       uint
		Version       uint
		IDTransaction string
		Status        string
	}
)

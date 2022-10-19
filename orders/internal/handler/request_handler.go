package handler

import "gorm.io/gorm"

type Handler struct {
	OrderHandler OrderHandler
}

func InitHandler(db *gorm.DB) *Handler {
	return &Handler{
		OrderHandler: *NewOrderHandler(db),
	}
}

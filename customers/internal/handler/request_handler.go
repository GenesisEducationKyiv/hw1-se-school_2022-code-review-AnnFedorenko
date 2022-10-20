package handler

import "gorm.io/gorm"

type Handler struct {
	CustomerHandler            CustomerHandler
	CustomerTransactionHandler CustomerTransactionHandler
}

func InitHandler(db *gorm.DB) *Handler {
	return &Handler{
		CustomerHandler:            *NewCustomerHandler(db),
		CustomerTransactionHandler: *NewCustomerTransactionHandler(db),
	}
}

package handler

import (
	"customers/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const globalID = "gid"

type CustomerHandler struct {
	db *gorm.DB
}

type registerCustomerRequest struct {
	EmailAddress string `json:"emailAddress"`
}

func NewCustomerHandler(db *gorm.DB) *CustomerHandler {
	return &CustomerHandler{db}
}

func (h *CustomerHandler) RegisterCustomer(c *gin.Context) interface{} {
	registerCustomerRequest := struct {
		Email string `json:"Address"`
	}{}
	transactionId := c.Query("gid")

	err := c.BindJSON(&registerCustomerRequest)
	if err != nil {
		return err
	}

	return h.db.
		Create(&model.Customer{
			IDTransaction: transactionId,
			Email:         registerCustomerRequest.Email,
			Balance:       100,
			Status:        "created",
		}).
		Error
}

func (h *CustomerHandler) RegisterCustomerCompensate(c *gin.Context) interface{} {
	transactionId := c.Query("gid")

	return h.db.
		Model(&model.Customer{}).
		Where("id_transaction = ?", transactionId).
		Update("status", "canceled").
		Limit(1).
		Error
}

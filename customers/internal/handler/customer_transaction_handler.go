package handler

import (
	"customers/internal/model"

	"github.com/dtm-labs/client/dtmcli"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CustomerTransactionHandler struct {
	db *gorm.DB
}

func NewCustomerTransactionHandler(db *gorm.DB) *CustomerTransactionHandler {
	return &CustomerTransactionHandler{db}
}

func (h *CustomerTransactionHandler) WithdrawMoney(c *gin.Context) interface{} {
	withdrawRequest := struct {
		IdCustomer uint `json:"idCustomer"`
		Amount     uint `json:"amount"`
	}{}
	transactionId := c.Query("gid")

	err := c.BindJSON(&withdrawRequest)
	if err != nil {
		return dtmcli.ErrFailure
	}

	// transaction
	err = h.db.Transaction(func(tx *gorm.DB) error {
		// find customer
		var customer model.Customer
		err = tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&customer, withdrawRequest.IdCustomer).
			Error
		if err != nil {
			return dtmcli.ErrFailure
		}
		// check balance
		if customer.Balance < withdrawRequest.Amount {
			return dtmcli.ErrFailure
		}

		// change
		customer.Version = customer.Version + 1
		customer.Balance -= withdrawRequest.Amount

		// save
		res := tx.Model(&model.Customer{}).
			Where("id = ? AND version = ?", customer.ID, customer.Version-1).
			Save(customer)
		if res.RowsAffected != 1 {
			return dtmcli.ErrOngoing
		}
		tx.Save(&model.ProcessedTransaction{IDTransaction: transactionId})

		return nil
	})

	return err
}

func (h *CustomerTransactionHandler) WithdrawMoneyCompensate(c *gin.Context) interface{} {
	compensateRequest := struct {
		IdCustomer uint `json:"idCustomer"`
		Amount     uint `json:"amount"`
	}{}
	transactionId := c.Query("gid")

	err := c.BindJSON(&compensateRequest)
	if err != nil {
		return dtmcli.ErrFailure
	}

	// transaction
	err = h.db.Transaction(func(tx *gorm.DB) error {
		// filter regular cases
		var pt model.ProcessedTransaction
		err = tx.Where(&model.ProcessedTransaction{IDTransaction: transactionId}).First(&pt).Error
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		// change customer
		var customer model.Customer
		err = tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&customer, compensateRequest.IdCustomer).
			Error
		if err != nil {
			return dtmcli.ErrFailure
		}

		customer.Version = customer.Version + 1
		customer.Balance = customer.Balance + compensateRequest.Amount

		// save
		result := tx.
			Model(&model.Customer{}).
			Where("id = ? AND version = ?", customer.ID, customer.Version-1).
			Updates(customer)
		if result.RowsAffected == 0 {
			return dtmcli.ErrOngoing
		}

		return nil
	})

	return err
}

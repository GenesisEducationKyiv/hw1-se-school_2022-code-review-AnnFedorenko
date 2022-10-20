package model

type Email struct {
	Address       string `form:"email" binding:"required"`
	TransactionID string
}

package service

import (
	"errors"
	"fmt"
	"net/mail"
	"net/smtp"
	"rate-api/config"
)

var ErrEmailSubscribed = errors.New("Email already subscribed")
var ErrEmailNotValid = errors.New("Email address is not valid")

type Email struct {
	Address string `form:"email" binding:"required"`
}

func AddEmail(newEmail Email) error {
	if !isEmailValid(newEmail.Address) {
		return ErrEmailNotValid
	}

	isEmailSubscribed, err := isEmailSubscribed(newEmail.Address)
	if err != nil {
		return err
	}
	if isEmailSubscribed {
		return ErrEmailSubscribed
	}

	return appendToFile(newEmail.Address, config.Cfg.EmailStorage)
}

func isEmailSubscribed(address string) (bool, error) {
	isFileExist, err := isFileExist(config.Cfg.EmailStorage)
	if err != nil {
		return false, err
	}
	if !isFileExist {
		return isFileExist, nil
	}

	isEmailExist, err := isStringExist(address, config.Cfg.EmailStorage)
	if err != nil {
		return isEmailExist, err
	}
	return isEmailExist, nil
}

func isEmailValid(address string) bool {
	_, err := mail.ParseAddress(address)
	return err == nil
}

func SendEmails() error {
	var cfg = config.Cfg
	user := cfg.SMTPUsername
	password := cfg.SMTPPassword
	addr := fmt.Sprintf("%s:%d", cfg.SMTPHost, cfg.SMTPPort)
	host := cfg.SMTPHost
	sender := "BTC rate app"

	receiver := readFileToArray(config.Cfg.EmailStorage)
	if receiver == nil {
		return nil
	}

	rate, err := GetRateFromBinance()
	if err != nil {
		return err
	}

	msg := []byte("From: Bitcoin rate helper\r\n" +
		"Subject: BTCUAH Rate\r\n\r\n" +
		rate.Price +
		"\r\n")

	auth := smtp.PlainAuth("", user, password, host)

	if err = smtp.SendMail(addr, auth, sender, receiver, msg); err != nil {
		return err
	}

	return nil
}

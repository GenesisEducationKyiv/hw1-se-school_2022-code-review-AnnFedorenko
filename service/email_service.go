package service

import (
	"errors"
	"net/mail"
	"rate-api/model"
	"rate-api/repository"
)

var ErrEmailSubscribed = errors.New("email already subscribed")
var ErrEmailNotValid = errors.New("email address is not valid")

type EmailService struct {
	repo repository.EmailRepository
}

func NewEmailService(repo repository.EmailRepository) *EmailService {
	return &EmailService{repo: repo}
}

func (s *EmailService) AddEmail(newEmail model.Email) error {
	if !isEmailValid(newEmail.Address) {
		return ErrEmailNotValid
	}

	isEmailSubscribed, err := s.repo.IsExist(newEmail)
	if err != nil {
		return err
	}

	if isEmailSubscribed {
		return ErrEmailSubscribed
	}

	return s.repo.Add(newEmail)
}

func isEmailValid(address string) bool {
	_, err := mail.ParseAddress(address)
	return err == nil
}

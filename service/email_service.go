package service

import (
	"net/mail"
	"rate-api/handler"
	"rate-api/mailerror"
	"rate-api/model"
)

type EmailService struct {
	repo EmailRepositoryInterface
}

type EmailRepositoryInterface interface {
	IsExist(email model.Email) (bool, error)
	Add(email model.Email) error
	GetAllEmails() []string
}

func NewEmailService(repo EmailRepositoryInterface) handler.EmailServiceInterface {
	return &EmailService{repo: repo}
}

func (s *EmailService) AddEmail(newEmail model.Email) error {
	if !isEmailValid(newEmail.Address) {
		return mailerror.ErrEmailNotValid
	}

	isEmailSubscribed, err := s.repo.IsExist(newEmail)
	if err != nil {
		return err
	}

	if isEmailSubscribed {
		return mailerror.ErrEmailSubscribed
	}

	return s.repo.Add(newEmail)
}

func (s *EmailService) GetAllEmails() []string {
	return s.repo.GetAllEmails()
}

func isEmailValid(address string) bool {
	_, err := mail.ParseAddress(address)
	return err == nil
}

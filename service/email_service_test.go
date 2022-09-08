package service_test

import (
	"log"
	"os"
	"rate-api/config"
	"rate-api/model"
	"rate-api/repository"
	"rate-api/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddEmailNotValid(t *testing.T) {
	config.LoadConfig()
	r := repository.NewEmailRepository(config.Cfg.EmailStorage)
	s := service.NewEmailService(*r)
	email := model.Email{Address: "invalid"}

	err := s.AddEmail(email)
	assert.EqualError(t, err, service.ErrEmailNotValid.Error())
}

func TestAddEmailDublicate(t *testing.T) {
	config.LoadConfig()
	r := repository.NewEmailRepository(config.Cfg.EmailStorage)
	s := service.NewEmailService(*r)
	email := model.Email{Address: config.Cfg.SMTPUsername}

	initTestFile()
	err := s.AddEmail(email)

	assert.EqualError(t, err, service.ErrEmailSubscribed.Error())
	t.Cleanup(func() { deleteTestFile() })
}

func deleteTestFile() {
	err := os.Remove(config.Cfg.EmailStorage)
	if err != nil {
		log.Panic(err)
	}
}

func initTestFile() {
	email := config.Cfg.SMTPUsername

	file, err := os.OpenFile(config.Cfg.EmailStorage, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()
	if _, err := file.WriteString(email + "\n"); err != nil {
		log.Panic(err)
	}
}

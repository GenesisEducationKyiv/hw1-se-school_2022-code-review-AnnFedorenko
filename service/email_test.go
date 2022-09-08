package service_test

import (
	"bufio"
	"log"
	"os"
	"rate-api/config"
	"rate-api/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddEmail(t *testing.T) {
	config.LoadConfig()
	email := service.Email{Address: config.Cfg.SMTPUsername}

	err := service.AddEmail(email)
	if err != nil {
		t.Fatalf("Unexpected error:\n %v", err)
	}

	assert.True(t, checkFileContainsString(email.Address))
	t.Cleanup(func() { deleteTestFile() })
}

func TestAddEmailNotValid(t *testing.T) {
	config.LoadConfig()
	email := service.Email{Address: "invalid"}

	err := service.AddEmail(email)
	assert.EqualError(t, err, service.ErrEmailNotValid.Error())
}

func TestAddEmailDublicate(t *testing.T) {
	config.LoadConfig()
	email := service.Email{Address: config.Cfg.SMTPUsername}

	initTestFile()
	err := service.AddEmail(email)

	assert.EqualError(t, err, service.ErrEmailSubscribed.Error())
	t.Cleanup(func() { deleteTestFile() })
}

func TestSendEmailsSuccessIntegration(t *testing.T) {
	config.LoadConfig()
	initTestFile()
	err := service.SendEmails()

	assert.NoError(t, err)
	t.Cleanup(func() { deleteTestFile() })
}

func TestSendEmailsFailIntegration(t *testing.T) {
	config.LoadConfig()
	config.Cfg.SMTPHost = "dummy"
	initTestFile()

	err := service.SendEmails()

	assert.Error(t, err)
	t.Cleanup(func() { deleteTestFile() })
}

func checkFileContainsString(str string) bool {
	file, err := os.Open(config.Cfg.EmailStorage)
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() == str {
			return true
		}
	}

	if err := scanner.Err(); err != nil {
		return false
	}

	return false
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

package service_test

import (
	"rate-api/config"
	"rate-api/repository"
	"rate-api/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendEmailsSuccessIntegration(t *testing.T) {
	config.LoadConfig()
	emailRepo := repository.NewEmailRepository(config.Cfg.EmailStorage)
	rateRepo := repository.NewRateRepository()
	rateServ := service.NewRateService(*rateRepo)
	emailSendServ := service.NewEmailSendService(*emailRepo, *rateServ)
	initTestFile()
	err := emailSendServ.SendEmails()

	assert.NoError(t, err)
	t.Cleanup(func() { deleteTestFile() })
}

func TestSendEmailsFailIntegration(t *testing.T) {
	config.LoadConfig()
	config.LoadConfig()
	emailRepo := repository.NewEmailRepository(config.Cfg.EmailStorage)
	rateRepo := repository.NewRateRepository()
	rateServ := service.NewRateService(*rateRepo)
	emailSendServ := service.NewEmailSendService(*emailRepo, *rateServ)
	config.Cfg.SMTPHost = "dummy"
	initTestFile()

	err := emailSendServ.SendEmails()

	assert.Error(t, err)
	t.Cleanup(func() { deleteTestFile() })
}

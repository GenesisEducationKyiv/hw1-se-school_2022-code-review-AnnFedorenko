package service_test

import (
	"rate-api/config"
	"rate-api/repository"
	"rate-api/service"
	"rate-api/service/rate"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type EmailSendRepositoryMock struct {
	mock.Mock
}

func TestSendEmailsSuccessIntegration(t *testing.T) {
	config.LoadConfig()
	emailServ := service.NewEmailService(repository.NewEmailRepository(config.Cfg.EmailStorage))
	rateServ := rate.NewBinanceRateService()
	emailSendServ := service.NewEmailSendService(emailServ, rateServ)
	initTestFile()

	err := emailSendServ.SendEmails()

	assert.NoError(t, err)
	t.Cleanup(func() { deleteTestFile() })
}

func TestSendEmailsFailIntegration(t *testing.T) {
	config.LoadConfig()
	emailServ := service.NewEmailService(repository.NewEmailRepository(config.Cfg.EmailStorage))
	rateServ := rate.NewBinanceRateService()
	emailSendServ := service.NewEmailSendService(emailServ, rateServ)
	config.Cfg.SMTPHost = "dummy"
	initTestFile()

	err := emailSendServ.SendEmails()

	assert.Error(t, err)
	t.Cleanup(func() { deleteTestFile() })
}

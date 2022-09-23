package service_test

import (
	"errors"
	"log"
	"rate-api/config"
	"rate-api/handler"
	"rate-api/mailclient"
	"rate-api/model"
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

type EmailClientMock struct {
	mock.Mock
}

func (c *EmailClientMock) Send(from string, to []string, msg []byte) error {
	args := c.Called(from, to, msg)
	return args.Error(0)
}

type EmailServiceMock struct {
	mock.Mock
}

func (s *EmailServiceMock) AddEmail(email model.Email) error {
	args := s.Called(email)
	return args.Error(0)
}

func (s *EmailServiceMock) GetAllEmails() []string {
	args := s.Called()
	arr, ok := args.Get(0).([]string)
	if !ok {
		if arr == nil {
			return nil
		}
		log.Panic("Assertion failed")
	}
	return arr
}

type RateServiceInterface interface {
	GetRate() (model.Rate, error)
	SetNext(next *RateServiceInterface)
}

type RateServiceMock struct {
	mock.Mock
}

func (s *RateServiceMock) GetRate() (model.Rate, error) {
	args := s.Called()
	rate, ok := args.Get(0).(model.Rate)
	if !ok {
		log.Panic("Assertion failed")
	}
	return rate, args.Error(1)
}

func (s *RateServiceMock) SetNext(next *handler.RateServiceInterface) {
	s.Called(next)
}

func TestSendEmailsSuccess(t *testing.T) {
	rs := new(RateServiceMock)
	es := new(EmailServiceMock)
	c := new(EmailClientMock)
	emailSendServ := service.NewEmailSendService(es, rs, c)
	emails := []string{"example@gmail.com"}

	es.Mock.On("GetAllEmails").Return(emails)
	rs.Mock.On("GetRate").Return(model.Rate{}, nil)
	c.Mock.On("Send", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	err := emailSendServ.SendEmails()

	assert.NoError(t, err)
	es.AssertNumberOfCalls(t, "GetAllEmails", 1)
	rs.AssertNumberOfCalls(t, "GetRate", 1)
	c.AssertNumberOfCalls(t, "Send", 1)
}

func TestSendEmailsNoReceiver(t *testing.T) {
	rs := new(RateServiceMock)
	es := new(EmailServiceMock)
	c := new(EmailClientMock)
	emailSendServ := service.NewEmailSendService(es, rs, c)

	es.Mock.On("GetAllEmails").Return(nil)
	err := emailSendServ.SendEmails()

	assert.NoError(t, err)
	es.AssertNumberOfCalls(t, "GetAllEmails", 1)
	rs.AssertNumberOfCalls(t, "GetRate", 0)
	c.AssertNumberOfCalls(t, "Send", 0)
}

func TestSendEmailsFailGetRate(t *testing.T) {
	rs := new(RateServiceMock)
	es := new(EmailServiceMock)
	c := new(EmailClientMock)
	emailSendServ := service.NewEmailSendService(es, rs, c)
	emails := []string{"example@gmail.com"}

	es.Mock.On("GetAllEmails").Return(emails)
	rs.Mock.On("GetRate").Return(model.Rate{}, rate.ErrRateFieldMissed)
	err := emailSendServ.SendEmails()

	assert.Error(t, err)
	assert.ErrorIs(t, err, rate.ErrRateFieldMissed)
	es.AssertNumberOfCalls(t, "GetAllEmails", 1)
	rs.AssertNumberOfCalls(t, "GetRate", 1)
	c.AssertNumberOfCalls(t, "Send", 0)
}

func TestSendEmailsFailSend(t *testing.T) {
	rs := new(RateServiceMock)
	es := new(EmailServiceMock)
	c := new(EmailClientMock)
	emailSendServ := service.NewEmailSendService(es, rs, c)
	emails := []string{"example@gmail.com"}

	es.Mock.On("GetAllEmails").Return(emails)
	rs.Mock.On("GetRate").Return(model.Rate{}, nil)
	c.Mock.On("Send", mock.Anything, mock.Anything, mock.Anything).Return(errors.New(""))
	err := emailSendServ.SendEmails()

	assert.Error(t, err)
	es.AssertNumberOfCalls(t, "GetAllEmails", 1)
	rs.AssertNumberOfCalls(t, "GetRate", 1)
	c.AssertNumberOfCalls(t, "Send", 1)
}

func TestSendEmailsSuccessIntegration(t *testing.T) {
	config.LoadConfig()
	emailServ := service.NewEmailService(repository.NewEmailRepository(config.Cfg.EmailStorage))
	rateServ := rate.NewBinanceRateService()
	emailSendServ := service.NewEmailSendService(emailServ, rateServ, mailclient.NewSMTPClient(config.Cfg))
	initTestFile()

	err := emailSendServ.SendEmails()

	assert.NoError(t, err)
	t.Cleanup(func() { deleteTestFile() })
}

func TestSendEmailsFailIntegration(t *testing.T) {
	config.LoadConfig()
	config.Cfg.SMTPHost = "dummy"
	emailServ := service.NewEmailService(repository.NewEmailRepository(config.Cfg.EmailStorage))
	rateServ := rate.NewBinanceRateService()
	emailSendServ := service.NewEmailSendService(emailServ, rateServ, mailclient.NewSMTPClient(config.Cfg))
	initTestFile()

	err := emailSendServ.SendEmails()

	assert.Error(t, err)
	t.Cleanup(func() { deleteTestFile() })
}

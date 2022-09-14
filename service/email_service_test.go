package service_test

import (
	"log"
	"os"
	"rate-api/config"
	"rate-api/mailerror"
	"rate-api/model"
	"rate-api/repository"
	"rate-api/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type EmailRepositoryMock struct {
	mock.Mock
}

func (s *EmailRepositoryMock) IsExist(email model.Email) (bool, error) {
	args := s.Called(email)
	return args.Bool(0), args.Error(1)
}

func (s *EmailRepositoryMock) Add(email model.Email) error {
	args := s.Called(email)
	return args.Error(0)
}

func (s *EmailRepositoryMock) GetAllEmails() []string {
	args := s.Called()
	arr, ok := args.Get(0).([]string)

	if ok == false {
		log.Panic("Assertion failed")
	}
	return arr
}

var testEmail = "test@gmail.com"

func TestAddEmail(t *testing.T) {
	r := new(EmailRepositoryMock)
	s := service.NewEmailService(r)
	email := model.Email{Address: testEmail}

	r.Mock.On("IsExist", email).Return(false, nil)
	r.Mock.On("Add", email).Return(nil)
	err := s.AddEmail(email)

	assert.NoError(t, err)
}
func TestAddEmailNotValid(t *testing.T) {
	r := new(EmailRepositoryMock)
	s := service.NewEmailService(r)
	email := model.Email{Address: "invalid"}

	err := s.AddEmail(email)

	assert.EqualError(t, err, mailerror.ErrEmailNotValid.Error())
}

func TestGetAllEmails(t *testing.T) {
	repo := new(EmailRepositoryMock)
	serv := service.NewEmailService(repo)
	var strArr []string

	repo.Mock.On("GetAllEmails").Return(strArr)
	serv.GetAllEmails()

	repo.AssertNumberOfCalls(t, "GetAllEmails", 1)
}

func TestAddEmailDublicate(t *testing.T) {
	repo := new(EmailRepositoryMock)
	serv := service.NewEmailService(repo)
	email := model.Email{Address: testEmail}

	repo.Mock.On("IsExist", email).Return(true, nil)
	err := serv.AddEmail(email)

	repo.AssertNumberOfCalls(t, "IsExist", 1)
	assert.EqualError(t, err, mailerror.ErrEmailSubscribed.Error())
}

func TestAddEmailDublicateIntegration(t *testing.T) {
	config.LoadConfig()
	r := repository.NewEmailRepository(config.Cfg.EmailStorage)
	s := service.NewEmailService(r)
	email := model.Email{Address: config.Cfg.SMTPUsername}

	initTestFile()
	err := s.AddEmail(email)

	assert.EqualError(t, err, mailerror.ErrEmailSubscribed.Error())
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

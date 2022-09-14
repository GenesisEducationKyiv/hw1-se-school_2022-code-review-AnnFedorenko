package repository_test

import (
	"bufio"
	"log"
	"os"
	"rate-api/config"
	"rate-api/model"
	"rate-api/repository"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddEmail(t *testing.T) {
	config.LoadConfig()
	email := model.Email{Address: config.Cfg.SMTPUsername}
	repo := repository.NewEmailRepository(config.Cfg.EmailStorage)

	err := repo.Add(email)
	if err != nil {
		t.Fatalf("Unexpected error:\n %v", err)
	}

	assert.True(t, checkFileContainsString(email.Address))
	t.Cleanup(func() { deleteTestFile() })
}

func TestIsExistTrue(t *testing.T) {
	config.LoadConfig()
	email := model.Email{Address: config.Cfg.SMTPUsername}
	repo := repository.NewEmailRepository(config.Cfg.EmailStorage)
	initTestFile()

	isExist, err := repo.IsExist(email)
	if err != nil {
		t.Fatalf("Unexpected error:\n %v", err)
	}

	assert.True(t, isExist)
	t.Cleanup(func() { deleteTestFile() })
}

func TestIsExistFalse(t *testing.T) {
	config.LoadConfig()
	email := model.Email{Address: config.Cfg.SMTPUsername}
	repo := repository.NewEmailRepository(config.Cfg.EmailStorage)
	initEmptyTestFile()

	isExist, err := repo.IsExist(email)
	if err != nil {
		t.Fatalf("Unexpected error:\n %v", err)
	}

	assert.False(t, isExist)
	t.Cleanup(func() { deleteTestFile() })
}

func TestIsExistFileNotExist(t *testing.T) {
	config.LoadConfig()
	email := model.Email{Address: config.Cfg.SMTPUsername}
	repo := repository.NewEmailRepository(config.Cfg.EmailStorage)

	isExist, err := repo.IsExist(email)
	if err != nil {
		t.Fatalf("Unexpected error:\n %v", err)
	}

	assert.False(t, isExist)
}

func TestGetAllUsers(t *testing.T) {
	config.LoadConfig()
	email := model.Email{Address: config.Cfg.SMTPUsername}
	repo := repository.NewEmailRepository(config.Cfg.EmailStorage)
	initTestFile()

	users := repo.GetAllEmails()

	assert.Len(t, users, 1)
	assert.Equal(t, email.Address, users[0])
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

func initEmptyTestFile() {
	file, err := os.OpenFile(config.Cfg.EmailStorage, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()
}

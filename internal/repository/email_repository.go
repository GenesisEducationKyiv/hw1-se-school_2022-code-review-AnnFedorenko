package repository

import (
	"bufio"
	"errors"
	"os"
	"rate-api/internal/model"
	"rate-api/internal/service"
)

type (
	EmailRepository struct {
		filePath string
	}
)

const fileMode = 0600

func NewEmailRepository(filePath string) service.EmailRepositoryInterface {
	return &EmailRepository{filePath: filePath}
}

func (r *EmailRepository) IsExist(email model.Email) (bool, error) {
	isFileExist, err := isFileExist(r.filePath)
	if err != nil {
		return false, err
	}
	if !isFileExist {
		return isFileExist, nil
	}

	file, err := os.Open(r.filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() == email.Address {
			return true, nil
		}
	}

	return false, nil
}

func (r *EmailRepository) Add(email model.Email) error {
	file, err := os.OpenFile(r.filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, fileMode)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(email.Address + "\n")

	return err
}

func (r *EmailRepository) GetAllEmails() []string {
	file, err := os.Open(r.filePath)
	if err != nil {
		return nil
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}

func isFileExist(fileName string) (bool, error) {
	_, err := os.Stat(fileName)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

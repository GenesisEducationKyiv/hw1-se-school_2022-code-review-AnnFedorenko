package service

import (
	"errors"
	"rate-api/model"
	"rate-api/repository"
)

var ErrRateFieldMissed = errors.New("rate field is missing")

type RateService struct {
	repo repository.BinanceRateRepository
}

func NewRateService(repo repository.BinanceRateRepository) *RateService {
	return &RateService{repo: repo}
}

func (s *RateService) GetRate() (model.Rate, error) {
	return s.repo.GetRate()
}

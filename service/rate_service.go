package service

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"rate-api/config"
	"rate-api/model"
	"rate-api/router"
)

var ErrRateFieldMissed = errors.New("rate field is missing")

type RateService struct {
}

func NewRateService() router.RateServiceInterface {
	return &RateService{}
}

func (s *RateService) GetRate() (model.Rate, error) {
	var newRate model.Rate
	resp, err := http.Get(config.Cfg.BtcURL)
	if err != nil {
		return newRate, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return newRate, err
	}

	if err := json.Unmarshal(body, &newRate); err != nil {
		return newRate, err
	}

	if newRate.Price == "" {
		return newRate, ErrRateFieldMissed
	}

	return newRate, nil
}

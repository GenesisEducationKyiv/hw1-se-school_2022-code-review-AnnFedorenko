package rate

import (
	"errors"
	"io"
	"net/http"
	"rate-api/config"
	"rate-api/model"
	"rate-api/router"
)

var ErrRateFieldMissed = errors.New("rate field is missing")

type CoingateRateService struct {
	next   *router.RateServiceInterface
	source string
}

func NewCoingateRateService() router.RateServiceInterface {
	return &CoingateRateService{source: "coingate"}
}

func (s *CoingateRateService) GetRate() (model.Rate, error) {
	rate, err := s.getRate()
	rate.Source = s.source
	if err != nil {
		if s.next == nil {
			return rate, err
		}
		next := *s.next
		rate, err = next.GetRate()
	}
	return rate, err
}

func (s *CoingateRateService) getRate() (model.Rate, error) {
	var newRate model.Rate
	resp, err := http.Get(config.Cfg.CoingateURL)
	if err != nil {
		return newRate, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return newRate, err
	}

	newRate.Price = string(body)
	if newRate.Price == "" {
		return newRate, ErrRateFieldMissed
	}

	return newRate, nil
}

func (s *CoingateRateService) SetNext(next *router.RateServiceInterface) {
	s.next = next
}

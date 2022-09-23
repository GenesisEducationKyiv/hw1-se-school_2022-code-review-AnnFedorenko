package rate

import (
	"encoding/json"
	"io"
	"net/http"
	"rate-api/config"
	"rate-api/model"
	"rate-api/router"
)

type BinanceRateService struct {
	next   *router.RateServiceInterface
	source string
}

func NewBinanceRateService() router.RateServiceInterface {
	return &BinanceRateService{source: "binance"}
}

func (s *BinanceRateService) GetRate() (model.Rate, error) {
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

func (s *BinanceRateService) getRate() (model.Rate, error) {
	var newRate model.Rate
	resp, err := http.Get(config.Cfg.BinanceURL)
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

func (s *BinanceRateService) SetNext(next *router.RateServiceInterface) {
	s.next = next
}

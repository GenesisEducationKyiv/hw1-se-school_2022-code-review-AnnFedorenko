package rate

import (
	"io"
	"net/http"
	"rate-api/config"
	"rate-api/model"
	"rate-api/router"

	"github.com/tidwall/gjson"
)

type CoinbaseRateService struct {
	next   *router.RateServiceInterface
	source string
}

func NewCoinbaseRateService() router.RateServiceInterface {
	return &CoinbaseRateService{source: "coinbase"}
}

func (s *CoinbaseRateService) GetRate() (model.Rate, error) {
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

func (s *CoinbaseRateService) getRate() (model.Rate, error) {
	var newRate model.Rate
	resp, err := http.Get(config.Cfg.CoinbaseURL)
	if err != nil {
		return newRate, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return newRate, err
	}

	price := gjson.GetBytes(body, "data.amount").String()
	if price == "" {
		return newRate, ErrRateFieldMissed
	}
	newRate.Price = price

	return newRate, nil
}

func (s *CoinbaseRateService) SetNext(next *router.RateServiceInterface) {
	s.next = next
}

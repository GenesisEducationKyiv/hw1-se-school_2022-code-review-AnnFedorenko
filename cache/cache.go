package cache

import (
	"rate-api/config"
	"rate-api/model"
	"rate-api/router"
	"time"

	"github.com/patrickmn/go-cache"
)

var rateKey = "rate"

type CacheRateService struct {
	c    cache.Cache
	serv *router.RateServiceInterface
}

func NewCacheRateService(serv router.RateServiceInterface) router.RateServiceInterface {
	duration := time.Duration(config.Cfg.LogDuration) * time.Minute
	c := cache.New(duration, duration+duration)
	return &CacheRateService{c: *c,
		serv: &serv}
}

func (s *CacheRateService) GetRate() (model.Rate, error) {
	cacheRate := s.getRate()
	if cacheRate.Price != "" {
		return cacheRate, nil
	}

	if s.serv == nil {
		return cacheRate, nil
	}
	serv := *s.serv
	newRate, err := serv.GetRate()
	if err != nil {
		return newRate, err
	}

	s.setRate(newRate)

	return newRate, nil
}

func (s *CacheRateService) SetNext(next *router.RateServiceInterface) {
	s.serv = next
}

func (s *CacheRateService) getRate() model.Rate {
	var rate model.Rate
	cache, ok := s.c.Get(rateKey)
	if !ok {
		return rate
	}

	rate, assertOk := cache.(model.Rate)
	if !assertOk {
		return rate
	}

	return rate
}

func (s *CacheRateService) setRate(rate model.Rate) {
	s.c.Set(rateKey, rate, cache.DefaultExpiration)
}

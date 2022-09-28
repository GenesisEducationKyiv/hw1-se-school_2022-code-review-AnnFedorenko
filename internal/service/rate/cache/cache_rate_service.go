package cache

import (
	"rate-api/config"
	"rate-api/internal/handler"
	"rate-api/internal/model"
	"time"

	"github.com/patrickmn/go-cache"
)

var rateKey = "rate"

type CacheRateServiceProxy struct {
	c    cache.Cache
	serv *handler.RateServiceInterface
}

func NewCacheRateServiceProxy(serv handler.RateServiceInterface) handler.RateServiceInterface {
	duration := time.Duration(config.Cfg.LogDuration) * time.Minute
	c := cache.New(duration, duration+duration)
	return &CacheRateServiceProxy{c: *c,
		serv: &serv}
}

func (s *CacheRateServiceProxy) GetRate() (model.Rate, error) {
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

func (s *CacheRateServiceProxy) SetNext(next *handler.RateServiceInterface) {
	s.serv = next
}

func (s *CacheRateServiceProxy) getRate() model.Rate {
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

func (s *CacheRateServiceProxy) setRate(rate model.Rate) {
	s.c.Set(rateKey, rate, cache.DefaultExpiration)
}

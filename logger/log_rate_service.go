package logger

import (
	"fmt"
	"rate-api/model"
	"rate-api/router"

	log "github.com/sirupsen/logrus"
)

type LogRateService struct {
	serv *router.RateServiceInterface
}

func NewLogRateService(serv router.RateServiceInterface) router.RateServiceInterface {
	return &LogRateService{
		serv: &serv}
}

func (s *LogRateService) GetRate() (model.Rate, error) {
	serv := *s.serv
	newRate, err := serv.GetRate()
	if err == nil {
		log.Info(fmt.Sprintf("%s - Response: %s", newRate.Source, newRate.Price))
	}

	return newRate, err
}

func (s *LogRateService) SetNext(next *router.RateServiceInterface) {
	s.serv = next
}

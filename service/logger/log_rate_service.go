package logger

import (
	"fmt"
	"rate-api/handler"
	"rate-api/model"

	log "github.com/sirupsen/logrus"
)

type LogRateService struct {
	serv *handler.RateServiceInterface
}

func NewLogRateService(serv handler.RateServiceInterface) handler.RateServiceInterface {
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

func (s *LogRateService) SetNext(next *handler.RateServiceInterface) {
	s.serv = next
}

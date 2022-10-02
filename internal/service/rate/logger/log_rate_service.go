package logger

import (
	"fmt"
	"rate-api/internal/handler"
	"rate-api/internal/model"

	log "github.com/sirupsen/logrus"
)

type LogRateServiceDecorator struct {
	serv *handler.RateServiceInterface
}

func NewLogRateServiceDecorator(serv handler.RateServiceInterface) handler.RateServiceInterface {
	return &LogRateServiceDecorator{
		serv: &serv}
}

func (s *LogRateServiceDecorator) GetRate() (model.Rate, error) {
	serv := *s.serv
	newRate, err := serv.GetRate()
	if err == nil {
		log.Info(fmt.Sprintf("%s - Response: %s", newRate.Source, newRate.Price))
	}

	return newRate, err
}

func (s *LogRateServiceDecorator) SetNext(next *handler.RateServiceInterface) {
	s.serv = next
}

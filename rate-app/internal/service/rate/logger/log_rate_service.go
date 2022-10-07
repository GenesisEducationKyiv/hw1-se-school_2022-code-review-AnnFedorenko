package logger

import (
	"fmt"
	"rate-api/internal/handler"
	"rate-api/internal/logger"
	"rate-api/internal/model"
)

type LogRateServiceDecorator struct {
	serv *handler.RateServiceInterface
	log  logger.LoggerInterface
}

func NewLogRateServiceDecorator(serv handler.RateServiceInterface,
	log logger.LoggerInterface) handler.RateServiceInterface {
	return &LogRateServiceDecorator{
		serv: &serv,
		log:  log}
}

func (s *LogRateServiceDecorator) GetRate() (model.Rate, error) {
	serv := *s.serv
	newRate, err := serv.GetRate()
	if err == nil {
		s.log.Info(fmt.Sprintf("%s - Response: %s", newRate.Source, newRate.Price))
	}

	return newRate, err
}

func (s *LogRateServiceDecorator) SetNext(next *handler.RateServiceInterface) {
	s.serv = next
}

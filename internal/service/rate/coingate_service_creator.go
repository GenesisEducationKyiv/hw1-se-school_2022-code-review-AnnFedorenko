package rate

import "rate-api/internal/handler"

type CoingateRateCreator struct {
}

func NewCoingateRateCreator() RateFactory {
	return CoingateRateCreator{}
}

func (c CoingateRateCreator) GetRateService() handler.RateServiceInterface {
	return NewCoingateRateService()
}

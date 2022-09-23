package rate

import "rate-api/router"

type CoingateRateCreator struct {
}

func NewCoingateRateCreator() RateFactory {
	return CoingateRateCreator{}
}

func (c CoingateRateCreator) GetRateService() router.RateServiceInterface {
	return NewCoingateRateService()
}

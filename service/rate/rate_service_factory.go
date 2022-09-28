package rate

import "rate-api/handler"

type RateFactory interface {
	GetRateService() handler.RateServiceInterface
}

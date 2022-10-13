package rate

import "rate-api/internal/handler"

type RateFactory interface {
	GetRateService() handler.RateServiceInterface
}

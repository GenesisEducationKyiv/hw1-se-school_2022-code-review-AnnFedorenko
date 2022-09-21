package rate

import (
	"rate-api/router"
)

type RateFactory interface {
	GetRateService() router.RateServiceInterface
}

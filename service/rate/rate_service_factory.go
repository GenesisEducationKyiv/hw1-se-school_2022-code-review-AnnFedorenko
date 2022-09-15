package rate

import (
	"log"
	"rate-api/router"
)

func GetRateService(servType string) router.RateServiceInterface {
	switch servType {
	case "binance":
		return NewBinanceRateService()
	case "coinbase":
		return NewCoinbaseRateService()
	case "coingate":
		return NewCoingateRateService()
	default:
		log.Fatal("Unknown rate provider")
	}
	return nil
}

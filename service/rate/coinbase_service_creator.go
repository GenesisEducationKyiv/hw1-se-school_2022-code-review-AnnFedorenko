package rate

import "rate-api/router"

type CoinbaseRateCreator struct {
}

func NewCoinbaseRateCreator() RateFactory {
	return CoinbaseRateCreator{}
}

func (c CoinbaseRateCreator) GetRateService() router.RateServiceInterface {
	return NewCoinbaseRateService()
}

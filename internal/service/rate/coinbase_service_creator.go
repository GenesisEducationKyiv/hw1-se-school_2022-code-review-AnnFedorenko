package rate

import "rate-api/internal/handler"

type CoinbaseRateCreator struct {
}

func NewCoinbaseRateCreator() RateFactory {
	return CoinbaseRateCreator{}
}

func (c CoinbaseRateCreator) GetRateService() handler.RateServiceInterface {
	return NewCoinbaseRateService()
}

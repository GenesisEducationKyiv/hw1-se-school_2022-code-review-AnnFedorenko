package rate

import "rate-api/router"

type BinanceRateCreator struct {
}

func NewBinanceRateCreator() RateFactory {
	return BinanceRateCreator{}
}

func (c BinanceRateCreator) GetRateService() router.RateServiceInterface {
	return NewBinanceRateService()
}

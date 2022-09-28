package rate

import "rate-api/handler"

type BinanceRateCreator struct {
}

func NewBinanceRateCreator() RateFactory {
	return BinanceRateCreator{}
}

func (c BinanceRateCreator) GetRateService() handler.RateServiceInterface {
	return NewBinanceRateService()
}

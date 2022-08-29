package service

import (
	"encoding/json"
	"io"
	"net/http"
)

type Rate struct {
	Price string
}

const BinanceURL = "https://api.binance.com/api/v3/ticker/price?symbol=BTCUAH"

func GetRateFromBinance() (Rate, error) {
	var newRate Rate
	resp, err := http.Get(BinanceURL)
	if err != nil {
		return newRate, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return newRate, err
	}

	if err := json.Unmarshal(body, &newRate); err != nil {
		return newRate, err
	}

	return newRate, nil
}

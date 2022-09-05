package service

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"rate-api/config"
)

var ErrRateFieldMissed = errors.New("Rate field is missing")

type Rate struct {
	Price string
}

func GetRateFromBinance() (Rate, error) {
	var newRate Rate
	resp, err := http.Get(config.Cfg.BtcURL)
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

	if newRate.Price == "" {
		return newRate, ErrRateFieldMissed
	}

	return newRate, nil
}

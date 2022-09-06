package service

import (
	"encoding/json"
	"io"
	"net/http"
	"rate-api/config"
)

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

	return newRate, nil
}

package service_test

import (
	"net/http"
	"rate-api/config"
	"rate-api/service"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGetRateFromBinance(t *testing.T) {
	config.LoadConfig()
	expectedRate := "772755.00000000"

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", config.Cfg.BtcURL,
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"symbol": "BTCUAH",
				"price":  expectedRate,
			})
			return resp, err
		},
	)

	rate, _ := service.GetRateFromBinance()

	assert.Equal(t, expectedRate, rate.Price)
}

func TestGetRateFromBinanceMissing(t *testing.T) {
	config.LoadConfig()
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", config.Cfg.BtcURL,
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"symbol": "BTCUAH",
			})
			return resp, err
		},
	)

	_, err := service.GetRateFromBinance()

	assert.Error(t, err)
	assert.ErrorIs(t, err, service.ErrRateFieldMissed)
}

func TestGetRateFromBinanceIntegration(t *testing.T) {
	config.LoadConfig()

	rate, err := service.GetRateFromBinance()

	assert.NoError(t, err)
	assert.NotEmpty(t, rate.Price)
}

func TestGetRateFromBinanceFailedIntegration(t *testing.T) {
	config.Cfg.BtcURL = "https://dummy"

	_, err := service.GetRateFromBinance()

	assert.Error(t, err)
}

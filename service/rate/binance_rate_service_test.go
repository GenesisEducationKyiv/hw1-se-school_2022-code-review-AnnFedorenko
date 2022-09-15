package rate_test

import (
	"net/http"
	"rate-api/config"
	"rate-api/service/rate"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

var binance = "binance"

func TestGetRate(t *testing.T) {
	config.LoadConfig()
	rateServ := rate.GetRateService(binance)

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

	rate, _ := rateServ.GetRate()

	assert.Equal(t, expectedRate, rate.Price)
}

func TestGetRateMissing(t *testing.T) {
	config.LoadConfig()
	rateServ := rate.GetRateService(binance)
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

	_, err := rateServ.GetRate()

	assert.Error(t, err)
	assert.ErrorIs(t, err, rate.ErrRateFieldMissed)
}

func TestGetRateIntegration(t *testing.T) {
	config.LoadConfig()
	rateServ := rate.GetRateService(binance)

	rate, err := rateServ.GetRate()

	assert.NoError(t, err)
	assert.NotEmpty(t, rate.Price)
}

func TestGetRateFailedIntegration(t *testing.T) {
	config.Cfg.BtcURL = "https://dummy"
	rateServ := rate.GetRateService(binance)

	_, err := rateServ.GetRate()

	assert.Error(t, err)
}

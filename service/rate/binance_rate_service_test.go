package rate_test

import (
	"net/http"
	"rate-api/config"
	"rate-api/service/rate"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGetRate(t *testing.T) {
	config.LoadConfig()
	rateServ := rate.NewBinanceRateService()

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", config.Cfg.BinanceURL,
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
	rateServ := rate.NewBinanceRateService()
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", config.Cfg.BinanceURL,
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
	rateServ := rate.NewBinanceRateService()

	rate, err := rateServ.GetRate()

	assert.NoError(t, err)
	assert.NotEmpty(t, rate.Price)
}

func TestGetRateFailedIntegration(t *testing.T) {
	config.Cfg.BinanceURL = dummyURL
	rateServ := rate.NewBinanceRateService()

	_, err := rateServ.GetRate()

	assert.Error(t, err)
}

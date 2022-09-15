package rate_test

import (
	"fmt"
	"net/http"
	"rate-api/config"
	"rate-api/service/rate"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

var coinbase = "coinbase"

func TestGetCoinbaseRate(t *testing.T) {
	config.LoadConfig()
	rateServ := rate.GetRateService(coinbase)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", config.Cfg.CoinbaseURL,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200,
				fmt.Sprintf(`{"data":{"amount":"%s"}}`, expectedRate))
			return resp, nil
		},
	)

	rate, _ := rateServ.GetRate()
	assert.Equal(t, expectedRate, rate.Price)
}

func TestGetCoinbaseRateMissing(t *testing.T) {
	config.LoadConfig()
	rateServ := rate.GetRateService(coinbase)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", config.Cfg.CoinbaseURL,
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

func TestGetCoinbaseRateIntegration(t *testing.T) {
	config.LoadConfig()
	rateServ := rate.GetRateService(coinbase)

	rate, err := rateServ.GetRate()

	assert.NoError(t, err)
	assert.NotEmpty(t, rate.Price)
}

func TestGetCoinbaseRateFailedIntegration(t *testing.T) {
	config.Cfg.CoinbaseURL = dummyURL
	rateServ := rate.GetRateService(coinbase)

	_, err := rateServ.GetRate()

	assert.Error(t, err)
}

package rate_test

import (
	"net/http"
	"rate-api/config"
	"rate-api/service/rate"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

var coingate = "coingate"

func TestGetCoingateRate(t *testing.T) {
	config.LoadConfig()
	rateServ := rate.GetRateService(coingate)

	expectedRate := "772755.00000000"

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", config.Cfg.CoingateURL,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, expectedRate)
			return resp, nil
		},
	)

	rate, _ := rateServ.GetRate()
	assert.Equal(t, expectedRate, rate.Price)
}

func TestGetCoingateRateMissing(t *testing.T) {
	config.LoadConfig()
	rateServ := rate.GetRateService(coingate)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", config.Cfg.CoingateURL,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, "")
			return resp, nil
		},
	)

	_, err := rateServ.GetRate()
	assert.Error(t, err)
	assert.ErrorIs(t, err, rate.ErrRateFieldMissed)
}

func TestGetCoingateRateIntegration(t *testing.T) {
	config.LoadConfig()
	rateServ := rate.GetRateService(coingate)

	rate, err := rateServ.GetRate()

	assert.NoError(t, err)
	assert.NotEmpty(t, rate.Price)
}

func TestGetCoingateRateFailedIntegration(t *testing.T) {
	config.Cfg.CoingateURL = "https://dummy"
	rateServ := rate.GetRateService(coingate)

	_, err := rateServ.GetRate()

	assert.Error(t, err)
}

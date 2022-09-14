package repository_test

import (
	"net/http"
	"rate-api/config"
	"rate-api/repository"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGetRate(t *testing.T) {
	config.LoadConfig()
	repo := repository.NewRateRepository()

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

	rate, _ := repo.GetRate()

	assert.Equal(t, expectedRate, rate.Price)
}

func TestGetRateMissing(t *testing.T) {
	config.LoadConfig()
	repo := repository.NewRateRepository()
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

	_, err := repo.GetRate()

	assert.Error(t, err)
	assert.ErrorIs(t, err, repository.ErrRateFieldMissed)
}

func TestGetRateIntegration(t *testing.T) {
	config.LoadConfig()
	repo := repository.NewRateRepository()

	rate, err := repo.GetRate()

	assert.NoError(t, err)
	assert.NotEmpty(t, rate.Price)
}

func TestGetRateFailedIntegration(t *testing.T) {
	config.Cfg.BtcURL = "https://dummy"
	repo := repository.NewRateRepository()

	_, err := repo.GetRate()

	assert.Error(t, err)
}

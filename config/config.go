package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

const ProjectName = "hw1-se-school_2022-code-review-AnnFedorenko"

type config struct {
	SMTPHost               string `env:"SMTP_HOST"        envDefault:"smtp.gmail.com1"`
	SMTPPort               int    `env:"SMTP_PORT"        envDefault:"587"`
	SMTPUsername           string `env:"EMAIL_ADDRESS"`
	SMTPPassword           string `env:"EMAIL_PASSWORD"`
	ServerHost             string `env:"SERVER_HOST"      envDefault:"localhost"`
	ServerPort             string `env:"SERVER_PORT"      envDefault:"8080"`
	BinanceURL             string `env:"BINANCE_URL"`
	CoinbaseURL            string `env:"COINBASE_URL"`
	CoingateURL            string `env:"COINGATE_URL"`
	EmailStorage           string `env:"EMAIL_STORAGE"`
	CryploCurrencyProvider string `env:"CRYPTO_CURRENCY_PROVIDER"`
}

var Cfg config

func LoadConfig() {
	if flag.Lookup("test.v") == nil {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("Error loading .env file")
		}
	} else {
		path, _ := os.Getwd()
		err := godotenv.Load(strings.Split(path, ProjectName)[0] + fmt.Sprintf("%s/", ProjectName) + ".env.test")
		if err != nil {
			log.Fatalf("Error loading .env.test file")
		}
	}

	Cfg = config{}

	err := env.Parse(&Cfg)
	if err != nil {
		log.Fatalf("unable to parse ennvironment variables: %e", err)
	}
}

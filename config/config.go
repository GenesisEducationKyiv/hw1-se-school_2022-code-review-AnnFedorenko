package config

import (
	"log"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type config struct {
	SMTPHost     string `env:"SMTP_HOST"        envDefault:"smtp.gmail.com1"`
	SMTPPort     int    `env:"SMTP_PORT"        envDefault:"587"`
	SMTPUsername string `env:"EMAIL_ADDRESS"`
	SMTPPassword string `env:"EMAIL_PASSWORD"`
	ServerHost   string `env:"SERVER_HOST"      envDefault:"localhost"`
	ServerPort   string `env:"SERVER_PORT"      envDefault:"8080"`
	BtcURL       string `env:"BTC_URL"`
}

var Cfg config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("unable to load .env file: %e", err)
	}

	Cfg = config{}

	err = env.Parse(&Cfg)
	if err != nil {
		log.Fatalf("unable to parse ennvironment variables: %e", err)
	}
}

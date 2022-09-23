package main

import (
	"rate-api/app"
	"rate-api/config"
)

func main() {
	config.LoadConfig()
	app.Run()
}

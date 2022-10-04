package main

import (
	"rate-api/config"
)

func main() {
	config.LoadConfig()
	Run()
}

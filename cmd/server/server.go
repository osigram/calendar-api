package main

import (
	"calendar-api/internal/app/server"
	"calendar-api/internal/config"
)

func main() {
	cfg := config.MustNewConfig()

	app := server.NewApp(cfg)

	app.Run()
}

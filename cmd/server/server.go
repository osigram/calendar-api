package main

import "calendar-api/internal/app/server"

func main() {
	app := server.NewApp()

	app.Run()
}

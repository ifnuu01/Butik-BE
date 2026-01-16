package main

import (
	"butik/internal/delivery/http"
	"butik/internal/infrastructure"

	Echo "github.com/labstack/echo/v4"
)

func main() {
	infrastructure.LoadEnv()
	db := infrastructure.SetupDB()
	e := Echo.New()

	http.RegisterRoutes(e, db)

	e.Logger.Fatal(e.Start(":" + infrastructure.GetEnv("PORT")))
}

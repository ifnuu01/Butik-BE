package main

import (
	"butik/internal/delivery/http"
	"butik/internal/infrastructure"
	"butik/pkg/utils"

	"github.com/go-playground/validator/v10"
	Echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	infrastructure.LoadEnv()
	db := infrastructure.SetupDB()
	e := Echo.New()
	e.Validator = &utils.CustomValidator{Validator: validator.New()}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{Echo.GET, Echo.PUT, Echo.POST, Echo.DELETE},
		AllowCredentials: true,
	}))
	http.RegisterRoutes(e, db)

	e.Logger.Fatal(e.Start(":" + infrastructure.GetEnv("PORT")))
}

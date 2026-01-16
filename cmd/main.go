package main

import (
	"butik/internal/delivery/http"
	"butik/internal/infrastructure"
	"butik/pkg/utils"

	"github.com/go-playground/validator/v10"
	Echo "github.com/labstack/echo/v4"
)

func main() {
	infrastructure.LoadEnv()
	db := infrastructure.SetupDB()
	e := Echo.New()
	e.Validator = &utils.CustomValidator{Validator: validator.New()}

	http.RegisterRoutes(e, db)

	e.Logger.Fatal(e.Start(":" + infrastructure.GetEnv("PORT")))
}

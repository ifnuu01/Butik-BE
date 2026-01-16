package http

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB) {
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "First Endpoint")
	})
}

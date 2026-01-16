package middlewares

import (
	"butik/internal/infrastructure"

	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware() echo.MiddlewareFunc {
	return echojwt.JWT(infrastructure.JWT_SECRET)
}

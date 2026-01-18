package http

import (
	"butik/internal/domain"
	"butik/internal/usecase"
	"butik/pkg/utils"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	loginPath        = "/login"
	refreshTokenPath = "/refresh-token"
)

type userHandler struct {
	Usecase usecase.UserUsecase
}

func RegisterUserRoutes(e *echo.Echo, userUsecase usecase.UserUsecase) {
	handler := &userHandler{Usecase: userUsecase}

	store := middleware.NewRateLimiterMemoryStoreWithConfig(middleware.RateLimiterMemoryStoreConfig{
		Rate:      5,
		Burst:     5,
		ExpiresIn: 1 * time.Minute,
	})
	loginLimiter := middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Store: store,
		IdentifierExtractor: func(c echo.Context) (string, error) {
			return c.RealIP(), nil
		},
	})

	e.POST(loginPath, handler.Login, loginLimiter)
	e.POST(refreshTokenPath, handler.RefreshToken)
}

func (h *userHandler) Login(c echo.Context) error {
	var req domain.LoginRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if err := c.Validate(&req); err != nil {
		return utils.ValidationErrorResponse(c, err)
	}

	res, err := h.Usecase.Login(req.Username, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}

func (h *userHandler) RefreshToken(c echo.Context) error {
	var req domain.RefreshTokenRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if err := c.Validate(&req); err != nil {
		return utils.ValidationErrorResponse(c, err)
	}

	res, err := h.Usecase.RefreshToken(req.RefreshToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}

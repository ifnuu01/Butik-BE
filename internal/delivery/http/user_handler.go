package http

import (
	"butik/internal/domain"
	"butik/internal/usecase"
	"butik/pkg/utils"
	"net/http"

	"github.com/labstack/echo/v4"
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
	e.POST(loginPath, handler.Login)
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

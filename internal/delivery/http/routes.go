package http

import (
	"butik/internal/repository"
	"butik/internal/usecase"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB) {
	userRepo := repository.NewUserRepo(db)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := &userHandler{Usecase: userUsecase}

	e.POST("/login", userHandler.Login)
	e.POST("/refresh-token", userHandler.RefreshToken)
}

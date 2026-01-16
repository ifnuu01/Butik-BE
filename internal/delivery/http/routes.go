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

	categoryRepo := repository.NewCategoryRepo(db)
	categoryUsecase := usecase.NewCategoryUsecase(categoryRepo)

	RegisterUserRoutes(e, userUsecase)
	RegisterCategoryRoutes(e, categoryUsecase)
}

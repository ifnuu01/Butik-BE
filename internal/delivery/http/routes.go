package http

import (
	"butik/internal/repository"
	"butik/internal/usecase"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB) {
	// User
	userRepo := repository.NewUserRepo(db)
	userUsecase := usecase.NewUserUsecase(userRepo)
	RegisterUserRoutes(e, userUsecase)

	// Category
	categoryRepo := repository.NewCategoryRepo(db)
	categoryUsecase := usecase.NewCategoryUsecase(categoryRepo)
	RegisterCategoryRoutes(e, categoryUsecase)

	// Product
	productRepo := repository.NewProductRepo(db)
	productUsecase := usecase.NewProductUsecase(productRepo, categoryRepo)
	RegisterProductRoutes(e, productUsecase)

	// Order
	orderRepo := repository.NewOrderRepo(db)
	orderUsecase := usecase.NewOrderUsecase(orderRepo, productRepo)
	RegisterOrderRoutes(e, orderUsecase)

	// static files (uploads)
	e.Static("/uploads", "uploads")
}

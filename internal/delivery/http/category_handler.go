package http

import (
	"butik/internal/delivery/http/middlewares"
	"butik/internal/domain"
	"butik/internal/usecase"
	"butik/pkg/utils"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

const (
	invalidCategoryIDMsg = "Invalid category ID"
	categoriesPath       = "/categories"
	categoryByIDPath     = "/categories/:id"
)

type categoryHandler struct {
	Usecase usecase.CategoryUsecase
}

func RegisterCategoryRoutes(e *echo.Echo, categoryUsecase usecase.CategoryUsecase) {
	handler := &categoryHandler{Usecase: categoryUsecase}
	e.GET(categoriesPath, handler.GetAllCategories)
	e.GET(categoryByIDPath, handler.GetCategoryByID)

	categoryGroup := e.Group(categoriesPath)
	categoryGroup.Use(middlewares.JWTMiddleware())
	categoryGroup.POST("", handler.CreateCategory)
	categoryGroup.PUT("/:id", handler.UpdateCategory)
	categoryGroup.DELETE("/:id", handler.DeleteCategory)
}

func (h *categoryHandler) CreateCategory(c echo.Context) error {
	var req domain.CreateCategoryRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if err := c.Validate(&req); err != nil {
		return utils.ValidationErrorResponse(c, err)
	}

	res, err := h.Usecase.CreateCategory(req.Name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, res)
}

func (h *categoryHandler) GetAllCategories(c echo.Context) error {
	pageStr := c.QueryParam("page")
	limitStr := c.QueryParam("limit")

	page := 1
	limit := 10

	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}

	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}

	offset := (page - 1) * limit
	categories, total, err := h.Usecase.GetAllCategories(offset, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	response := map[string]interface{}{
		"data":  categories,
		"page":  page,
		"limit": limit,
		"total": total,
	}
	return c.JSON(http.StatusOK, response)
}

func (h *categoryHandler) GetCategoryByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	if id <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": invalidCategoryIDMsg})
	}

	res, err := h.Usecase.GetCategoryByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}

func (h *categoryHandler) UpdateCategory(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var req domain.UpdateCategoryRequest

	if id <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": invalidCategoryIDMsg})
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if err := c.Validate(&req); err != nil {
		return utils.ValidationErrorResponse(c, err)
	}

	res, err := h.Usecase.UpdateCategory(uint(id), req.Name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}

func (h *categoryHandler) DeleteCategory(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	if id <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": invalidCategoryIDMsg})
	}

	res, err := h.Usecase.DeleteCategory(uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

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

type productHandler struct {
	Usecase usecase.ProductUsecase
}

func RegisterProductRoutes(e *echo.Echo, productUsecase usecase.ProductUsecase) {
	handler := &productHandler{Usecase: productUsecase}

	// Public
	e.GET("/products", handler.GetAllProducts)
	e.GET("/products/:id", handler.GetProductByID)

	// Protected
	productGroup := e.Group("/products", middlewares.JWTMiddleware())
	productGroup.POST("", handler.CreateProduct)
	productGroup.PUT("/:id", handler.UpdateProduct)
	productGroup.DELETE("/:id", handler.DeleteProduct)
}

func (h *productHandler) CreateProduct(c echo.Context) error {
	var req domain.CreateProductRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if err := c.Validate(&req); err != nil {
		return utils.ValidationErrorResponse(c, err)
	}

	// Handle file upload
	imageURL, err := utils.HandleFileUpload(c, "image", "uploads/products")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "image is required"})
	}

	res, err := h.Usecase.CreateProduct(req, imageURL)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, res)
}

func (h *productHandler) GetAllProducts(c echo.Context) error {
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
	products, total, err := h.Usecase.GetAllProducts(offset, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	response := map[string]interface{}{
		"data":  products,
		"total": total,
		"page":  page,
		"limit": limit,
	}
	return c.JSON(http.StatusOK, response)
}

func (h *productHandler) GetProductByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid product id"})
	}

	product, err := h.Usecase.GetProductByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, product)
}

func (h *productHandler) UpdateProduct(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid product id"})
	}

	var req domain.UpdateProductRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if err := c.Validate(&req); err != nil {
		return utils.ValidationErrorResponse(c, err)
	}

	// Handle file upload
	imageURL, _ := utils.HandleFileUpload(c, "image", "uploads/products")

	res, err := h.Usecase.UpdateProduct(uint(id), req, imageURL)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}

func (h *productHandler) DeleteProduct(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid product id"})
	}

	res, err := h.Usecase.DeleteProduct(uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}

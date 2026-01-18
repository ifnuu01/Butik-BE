package http

import (
	"butik/internal/delivery/http/middlewares"
	"butik/internal/domain"
	"butik/internal/usecase"
	"butik/pkg/utils"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type orderHandler struct {
	Usecase usecase.OrderUsecase
}

func RegisterOrderRoutes(e *echo.Echo, orderUsecase usecase.OrderUsecase) {
	handler := &orderHandler{Usecase: orderUsecase}

	// Public
	e.POST("/orders", handler.CreateOrder)
	e.GET("/orders/:id", handler.GetOrderByID)

	// Protected
	orderGroup := e.Group("/orders", middlewares.JWTMiddleware())
	orderGroup.GET("", handler.GetAllOrders)
	orderGroup.PUT("/:id/status", handler.UpdateOrderStatus)
	orderGroup.DELETE("/:id", handler.DeleteOrder)
}

func (h *orderHandler) CreateOrder(c echo.Context) error {
	var req domain.CreateOrderRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request format"})
	}

	itemsStr := c.FormValue("items")
	if itemsStr != "" {
		if err := json.Unmarshal([]byte(itemsStr), &req.Items); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "items must be a valid JSON array"})
		}
	}

	if err := c.Validate(&req); err != nil {
		return utils.ValidationErrorResponse(c, err)
	}

	// Handle file upload
	proofOfPayment, err := utils.HandleFileUpload(c, "proof_of_payment", "uploads/payments")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "proof of payment is required"})
	}

	res, err := h.Usecase.CreateOrder(req, proofOfPayment)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, res)
}

func (h *orderHandler) GetAllOrders(c echo.Context) error {
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
	orders, total, err := h.Usecase.GetAllOrders(offset, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	response := map[string]interface{}{
		"data":  orders,
		"total": total,
		"page":  page,
		"limit": limit,
	}
	return c.JSON(http.StatusOK, response)
}

func (h *orderHandler) GetOrderByID(c echo.Context) error {
	id := c.Param("id")

	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "order id is required"})
	}

	res, err := h.Usecase.GetOrderByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}

func (h *orderHandler) UpdateOrderStatus(c echo.Context) error {
	id := c.Param("id")

	var req domain.UpdateOrderStatusRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if err := c.Validate(&req); err != nil {
		return utils.ValidationErrorResponse(c, err)
	}

	res, err := h.Usecase.UpdateOrderStatus(id, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}

func (h *orderHandler) DeleteOrder(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "order id is required"})
	}
	if err := h.Usecase.DeleteOrder(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "order deleted successfully"})
}

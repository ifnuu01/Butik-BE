package usecase

import (
	"butik/internal/domain"
	"butik/internal/domain/dto"
	"butik/internal/repository"
	"errors"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

type OrderUsecase interface {
	CreateOrder(req domain.CreateOrderRequest, proofOfPayment string) (*domain.CreateOrderResponse, error)
	GetAllOrders(offset, limit int) ([]*domain.OrderResponse, int, error)
	GetOrderByID(id string) (*domain.OrderResponse, error)
	UpdateOrderStatus(id string, req domain.UpdateOrderStatusRequest) (*domain.UpdateOrderStatusResponse, error)
	DeleteOrder(id string) error
}

type orderUsecase struct {
	orderRepo   repository.OrderRepo
	productRepo repository.ProductRepo
}

func NewOrderUsecase(orderRepo repository.OrderRepo, productRepo repository.ProductRepo) OrderUsecase {
	return &orderUsecase{
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

func (u *orderUsecase) CreateOrder(req domain.CreateOrderRequest, proofOfPayment string) (*domain.CreateOrderResponse, error) {
	// Generate NanoID
	orderID, err := gonanoid.New()
	if err != nil {
		return nil, errors.New("failed to generate order ID")
	}

	var totalPrice float64
	var orderItems []domain.OrderItem
	var stockUpdates []struct {
		ProductID uint
		NewStock  int
	}

	// Validasi semua product dan stock
	for _, item := range req.Items {
		product, err := u.productRepo.GetProductByID(item.ProductID)
		if err != nil {
			return nil, errors.New("product not found")
		}

		if product.Stock < item.Quantity {
			return nil, errors.New("stock not enough for product: " + product.Name)
		}

		priceAtPurchase := product.Price * float64(item.Quantity)
		totalPrice += priceAtPurchase

		orderItems = append(orderItems, domain.OrderItem{
			OrderID:         orderID,
			ProductID:       item.ProductID,
			Product:         *product,
			Quantity:        item.Quantity,
			PriceAtPurchase: product.Price,
		})

		stockUpdates = append(stockUpdates, struct {
			ProductID uint
			NewStock  int
		}{
			ProductID: product.ID,
			NewStock:  product.Stock - item.Quantity,
		})
	}

	order := domain.Order{
		ID:             orderID,
		CustomerName:   req.CustomerName,
		Whatsapp:       req.Whatsapp,
		MapAddress:     req.MapAddress,
		Latitude:       req.Latitude,
		Longitude:      req.Longitude,
		AddressNote:    req.AddressNote,
		TotalPrice:     totalPrice,
		ProofOfPayment: proofOfPayment,
		Status:         domain.OrderStatusPending,
		OrderItems:     orderItems,
	}

	// Create order dengan transaction
	createdOrder, err := u.orderRepo.CreateOrderWithTransaction(order, stockUpdates)
	if err != nil {
		return nil, err
	}

	return &domain.CreateOrderResponse{
		Message: "Order created successfully",
		Order:   *dto.ToOrderResponse(createdOrder),
	}, nil
}

func (u *orderUsecase) GetAllOrders(offset, limit int) ([]*domain.OrderResponse, int, error) {
	orders, total, err := u.orderRepo.GetAllOrders(offset, limit)
	if err != nil {
		return nil, 0, err
	}
	orderResponses := dto.ToOrderResponses(orders)
	return orderResponses, total, nil
}

func (u *orderUsecase) GetOrderByID(id string) (*domain.OrderResponse, error) {
	order, err := u.orderRepo.GetOrderByID(id)
	if err != nil {
		return nil, err
	}
	return dto.ToOrderResponse(order), nil
}

func (u *orderUsecase) UpdateOrderStatus(id string, req domain.UpdateOrderStatusRequest) (*domain.UpdateOrderStatusResponse, error) {
	order, err := u.orderRepo.UpdateOrderStatus(id, req.Status)
	if err != nil {
		return nil, err
	}
	return &domain.UpdateOrderStatusResponse{
		Message: "Order status updated successfully",
		Order:   *dto.ToOrderResponse(order),
	}, nil
}

func (u *orderUsecase) DeleteOrder(id string) error {
	return u.orderRepo.DeleteOrder(id)
}

package dto

import (
	"butik/internal/domain"
	"time"
)

func ToOrderItemResponse(item *domain.OrderItem) domain.OrderItemResponse {
	return domain.OrderItemResponse{
		ID:              item.ID,
		ProductID:       item.ProductID,
		Product:         *ToProductResponse(&item.Product),
		Quantity:        item.Quantity,
		PriceAtPurchase: item.PriceAtPurchase,
	}
}

func ToOrderItemResponses(items []domain.OrderItem) []domain.OrderItemResponse {
	responses := make([]domain.OrderItemResponse, len(items))
	for i, item := range items {
		responses[i] = ToOrderItemResponse(&item)
	}
	return responses
}

func ToOrderResponse(order *domain.Order) *domain.OrderResponse {
	return &domain.OrderResponse{
		ID:             order.ID,
		CustomerName:   order.CustomerName,
		Whatsapp:       order.Whatsapp,
		MapAddress:     order.MapAddress,
		Latitude:       order.Latitude,
		Longitude:      order.Longitude,
		AddressNote:    order.AddressNote,
		TotalPrice:     order.TotalPrice,
		ProofOfPayment: order.ProofOfPayment,
		Status:         order.Status,
		OrderItems:     ToOrderItemResponses(order.OrderItems),
		CreatedAt:      order.CreatedAt.Format(time.RFC3339),
	}
}

func ToOrderResponses(orders []domain.Order) []*domain.OrderResponse {
	responses := make([]*domain.OrderResponse, len(orders))
	for i, order := range orders {
		responses[i] = ToOrderResponse(&order)
	}
	return responses
}

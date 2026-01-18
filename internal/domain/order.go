package domain

import "time"

type OrderStatus string

const (
	OrderStatusPending  OrderStatus = "pending"
	OrderStatusSuccess  OrderStatus = "success"
	OrderStatusRejected OrderStatus = "rejected"
)

type Order struct {
	ID             string      `gorm:"primaryKey" json:"id"`
	CustomerName   string      `json:"customer_name"`
	Whatsapp       string      `json:"whatsapp"`
	MapAddress     string      `json:"map_address"`
	Latitude       float64     `json:"latitude"`
	Longitude      float64     `json:"longitude"`
	AddressNote    string      `json:"address_note"`
	TotalPrice     float64     `json:"total_price"`
	ProofOfPayment string      `json:"proof_of_payment"`
	Status         OrderStatus `gorm:"default:pending" json:"status"`
	OrderItems     []OrderItem `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE;" json:"order_items"`
	CreatedAt      time.Time   `json:"created_at"`
}

type OrderItem struct {
	ID              uint    `gorm:"primaryKey" json:"id"`
	OrderID         string  `json:"order_id"`
	ProductID       uint    `json:"product_id"`
	Product         Product `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"product"`
	Quantity        int     `json:"quantity"`
	PriceAtPurchase float64 `json:"price_at_purchase"`
}

// Request DTOs
type OrderItemRequest struct {
	ProductID uint `json:"product_id" validate:"required,gt=0"`
	Quantity  int  `json:"quantity" validate:"required,gt=0,lte=100"`
}

type CreateOrderRequest struct {
	CustomerName string             `json:"customer_name" form:"customer_name" validate:"required,min=2,max=100"`
	Whatsapp     string             `json:"whatsapp" form:"whatsapp" validate:"required,min=10,max=15"`
	MapAddress   string             `json:"map_address" form:"map_address" validate:"max=500"`
	Latitude     float64            `json:"latitude" form:"latitude" validate:"gte=-90,lte=90"`
	Longitude    float64            `json:"longitude" form:"longitude" validate:"gte=-180,lte=180"`
	AddressNote  string             `json:"address_note" form:"address_note" validate:"max=500"`
	Items        []OrderItemRequest `json:"items" validate:"required,min=1,max=50,dive"`
}

type UpdateOrderStatusRequest struct {
	Status OrderStatus `json:"status" validate:"required,oneof=pending success rejected"`
}

// Response DTOs
type OrderItemResponse struct {
	ID              uint            `json:"id"`
	ProductID       uint            `json:"product_id"`
	Product         ProductResponse `json:"product"`
	Quantity        int             `json:"quantity"`
	PriceAtPurchase float64         `json:"price_at_purchase"`
}

type OrderResponse struct {
	ID             string              `json:"id"`
	CustomerName   string              `json:"customer_name"`
	Whatsapp       string              `json:"whatsapp"`
	MapAddress     string              `json:"map_address"`
	Latitude       float64             `json:"latitude"`
	Longitude      float64             `json:"longitude"`
	AddressNote    string              `json:"address_note"`
	TotalPrice     float64             `json:"total_price"`
	ProofOfPayment string              `json:"proof_of_payment"`
	Status         OrderStatus         `json:"status"`
	OrderItems     []OrderItemResponse `json:"order_items"`
	CreatedAt      string              `json:"created_at"`
}

type CreateOrderResponse struct {
	Message string        `json:"message"`
	Order   OrderResponse `json:"order"`
}

type GetOrderResponse struct {
	Order OrderResponse `json:"order"`
}

type GetOrdersResponse struct {
	Orders []*OrderResponse `json:"orders"`
}

type UpdateOrderStatusResponse struct {
	Message string        `json:"message"`
	Order   OrderResponse `json:"order"`
}

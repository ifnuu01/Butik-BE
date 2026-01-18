package repository

import (
	"butik/internal/domain"
	"errors"

	"gorm.io/gorm"
)

type OrderRepo interface {
	CreateOrderWithTransaction(order domain.Order, stockUpdates []struct {
		ProductID uint
		NewStock  int
	}) (*domain.Order, error)
	GetAllOrders(offset, limit int) ([]domain.Order, int, error)
	GetOrderByID(id string) (*domain.Order, error)
	UpdateOrderStatus(id string, status domain.OrderStatus) (*domain.Order, error)
	DeleteOrder(id string) error
}

type orderRepo struct {
	db *gorm.DB
}

func NewOrderRepo(db *gorm.DB) OrderRepo {
	return &orderRepo{db: db}
}

func (r *orderRepo) CreateOrderWithTransaction(order domain.Order, stockUpdates []struct {
	ProductID uint
	NewStock  int
}) (*domain.Order, error) {
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, errors.New("failed to start transaction")
	}

	// Update stock
	for _, su := range stockUpdates {
		result := tx.Model(&domain.Product{}).Where("id = ?", su.ProductID).Update("stock", su.NewStock)
		if result.Error != nil {
			tx.Rollback()
			return nil, errors.New("failed to update stock")
		}
	}
	// order transaction
	result := tx.Create(&order)
	if result.Error != nil {
		tx.Rollback()
		return nil, errors.New("failed to create order")
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, errors.New("failed to commit transaction")
	}

	return &order, nil
}

func (r *orderRepo) GetAllOrders(offset, limit int) ([]domain.Order, int, error) {
	var orders []domain.Order
	var total int64

	if err := r.db.Model(&domain.Order{}).Count(&total).Error; err != nil {
		return nil, 0, errors.New("failed to count orders")
	}

	if err := r.db.Preload("OrderItems.Product.Category").Order("created_at DESC").Offset(offset).Limit(limit).Find(&orders).Error; err != nil {
		return nil, 0, errors.New("failed to retrieve orders")
	}
	return orders, int(total), nil
}

func (r *orderRepo) GetOrderByID(id string) (*domain.Order, error) {
	order := &domain.Order{}
	result := r.db.Preload("OrderItems.Product.Category").First(order, "id = ?", id)
	if result.Error != nil {
		return nil, errors.New("order not found")
	}
	return order, nil
}

func (r *orderRepo) UpdateOrderStatus(id string, status domain.OrderStatus) (*domain.Order, error) {
	order, err := r.GetOrderByID(id)
	if err != nil {
		return nil, err
	}
	order.Status = status
	result := r.db.Save(order)
	if result.Error != nil {
		return nil, errors.New("failed to update order status")
	}
	return order, nil
}

func (r *orderRepo) DeleteOrder(id string) error {
	result := r.db.Delete(&domain.Order{}, "id = ?", id)
	if result.Error != nil {
		return errors.New("failed to delete order")
	}
	return nil
}

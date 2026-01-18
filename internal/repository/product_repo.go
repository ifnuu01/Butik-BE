package repository

import (
	"butik/internal/domain"
	"errors"

	"gorm.io/gorm"
)

type ProductRepo interface {
	CreateProduct(domain.Product) (*domain.Product, error)
	GetAllProducts(offset, limit int) ([]domain.Product, int, error)
	GetProductByID(id uint) (*domain.Product, error)
	UpdateProduct(id uint, product domain.Product) (*domain.Product, error)
	DeleteProduct(id uint) error
}

type productRepo struct {
	db *gorm.DB
}

func NewProductRepo(db *gorm.DB) ProductRepo {
	return &productRepo{db: db}
}

func (r *productRepo) CreateProduct(product domain.Product) (*domain.Product, error) {
	result := r.db.Create(&product)
	if result.Error != nil {
		return nil, errors.New("Failed to create product")
	}
	return &product, nil
}

func (r *productRepo) GetAllProducts(offset, limit int) ([]domain.Product, int, error) {
	var products []domain.Product
	var total int64

	if err := r.db.Model(&domain.Product{}).Count(&total).Error; err != nil {
		return nil, 0, errors.New("Failed to count products")
	}

	if err := r.db.Preload("Category").Order("created_at DESC").Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, 0, errors.New("Failed to retrieve products")
	}

	return products, int(total), nil
}

func (r *productRepo) GetProductByID(id uint) (*domain.Product, error) {
	product := &domain.Product{}
	result := r.db.Preload("Category").First(product, id)
	if result.Error != nil {
		return nil, errors.New("Product not found")
	}
	return product, nil
}

func (r *productRepo) UpdateProduct(id uint, updatedProduct domain.Product) (*domain.Product, error) {
	product, err := r.GetProductByID(id)
	if err != nil {
		return nil, errors.New("Product not found")
	}
	product.Name = updatedProduct.Name
	product.Price = updatedProduct.Price
	product.Description = updatedProduct.Description
	product.Stock = updatedProduct.Stock
	product.CategoryID = updatedProduct.CategoryID
	product.ImageURL = updatedProduct.ImageURL

	result := r.db.Save(product)
	if result.Error != nil {
		return nil, result.Error
	}
	return product, nil
}

func (r *productRepo) DeleteProduct(id uint) error {
	product, err := r.GetProductByID(id)
	if err != nil {
		return errors.New("Product not found")
	}
	result := r.db.Delete(product)
	if result.Error != nil {
		return errors.New("Failed to delete product")
	}
	return nil
}

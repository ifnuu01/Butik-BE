package domain

import "time"

type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	CategoryID  uint      `json:"category_id"`
	Category    Category  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"category"`
	ImageURL    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
}

type ProductResponse struct {
	ID          uint             `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Price       float64          `json:"price"`
	Stock       int              `json:"stock"`
	Category    CategoryResponse `json:"category"`
	ImageURL    string           `json:"image_url"`
	CreatedAt   string           `json:"created_at"`
}

type CreateProductRequest struct {
	Name        string  `json:"name" form:"name" validate:"required,min=2,max=200"`
	Description string  `json:"description" form:"description" validate:"required,min=10,max=2000"`
	Price       float64 `json:"price" form:"price" validate:"required,gt=0,lte=999999999"`
	Stock       int     `json:"stock" form:"stock" validate:"required,gte=0,lte=99999"`
	CategoryID  uint    `json:"category_id" form:"category_id" validate:"required,gt=0"`
}

type UpdateProductRequest struct {
	Name        string  `json:"name" form:"name" validate:"required,min=2,max=200"`
	Description string  `json:"description" form:"description" validate:"required,min=10,max=2000"`
	Price       float64 `json:"price" form:"price" validate:"required,gt=0,lte=999999999"`
	Stock       int     `json:"stock" form:"stock" validate:"required,gte=0,lte=99999"`
	CategoryID  uint    `json:"category_id" form:"category_id" validate:"required,gt=0"`
}

type CreateProductResponse struct {
	Message string          `json:"message"`
	Product ProductResponse `json:"product"`
}

type UpdateProductResponse struct {
	Message string          `json:"message"`
	Product ProductResponse `json:"product"`
}

type DeleteProductResponse struct {
	Message string `json:"message"`
}

package domain

import "time"

type Category struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"unique;not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type CategoryResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

type CreateCategoryRequest struct {
	Name string `json:"name" form:"name" validate:"required"`
}

type CreateCategoryResponse struct {
	Message  string           `json:"message"`
	Category CategoryResponse `json:"category"`
}

type UpdateCategoryRequest struct {
	Name string `json:"name" form:"name" validate:"required"`
}

type UpdateCategoryResponse struct {
	Message  string           `json:"message"`
	Category CategoryResponse `json:"category"`
}

type DeleteCategoryResponse struct {
	Message string `json:"message"`
}

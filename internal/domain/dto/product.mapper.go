package dto

import (
	"butik/internal/domain"
	"time"
)

func ToProductResponse(prod *domain.Product) *domain.ProductResponse {
	return &domain.ProductResponse{
		ID:          prod.ID,
		Name:        prod.Name,
		Description: prod.Description,
		Price:       prod.Price,
		Stock:       prod.Stock,
		Category:    *ToCategoryResponse(&prod.Category),
		ImageURL:    prod.ImageURL,
		CreatedAt:   prod.CreatedAt.Format(time.RFC3339),
	}
}

func ToProductResponses(prods []domain.Product) []*domain.ProductResponse {
	responses := make([]*domain.ProductResponse, len(prods))
	for i, prod := range prods {
		responses[i] = ToProductResponse(&prod)
	}
	return responses
}

package dto

import (
	"butik/internal/domain"
	"time"
)

func ToCategoryResponse(cat *domain.Category) *domain.CategoryResponse {
	return &domain.CategoryResponse{
		ID:        cat.ID,
		Name:      cat.Name,
		CreatedAt: cat.CreatedAt.Format(time.RFC3339),
	}
}

func ToCategoryResponses(cats []domain.Category) []*domain.CategoryResponse {
	responses := make([]*domain.CategoryResponse, len(cats))
	for i, cat := range cats {
		responses[i] = ToCategoryResponse(&cat)
	}
	return responses
}

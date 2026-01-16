package usecase

import (
	"butik/internal/domain"
	"butik/internal/domain/dto"
	"butik/internal/repository"
	"errors"
)

type CategoryUsecase interface {
	CreateCategory(name string) (*domain.CreateCategoryResponse, error)
	GetAllCategories() ([]*domain.CategoryResponse, error)
	GetCategoryByID(id uint) (*domain.CategoryResponse, error)
	UpdateCategory(id uint, name string) (*domain.UpdateCategoryResponse, error)
	DeleteCategory(id uint) (*domain.DeleteCategoryResponse, error)
}

type categoryUsecase struct {
	categoryRepo repository.CategoryRepo
}

func NewCategoryUsecase(categoryRepo repository.CategoryRepo) CategoryUsecase {
	return &categoryUsecase{categoryRepo: categoryRepo}
}

func (u *categoryUsecase) CreateCategory(name string) (*domain.CreateCategoryResponse, error) {
	category, err := u.categoryRepo.CreateCategory(name)
	if err != nil {
		return nil, errors.New("failed to create category")
	}

	catResp := dto.ToCategoryResponse(category)

	return &domain.CreateCategoryResponse{
		Message:  "Category created successfully",
		Category: *catResp,
	}, nil
}

func (u *categoryUsecase) GetAllCategories() ([]*domain.CategoryResponse, error) {
	categories, err := u.categoryRepo.GetAllCategories()
	if err != nil {
		return nil, errors.New("failed to get categories")
	}

	catResponses := dto.ToCategoryResponses(categories)
	return catResponses, nil
}

func (u *categoryUsecase) GetCategoryByID(id uint) (*domain.CategoryResponse, error) {
	category, err := u.categoryRepo.GetCategoryByID(id)
	if err != nil {
		return nil, errors.New("category not found")
	}
	catResp := dto.ToCategoryResponse(category)
	return catResp, nil
}

func (u *categoryUsecase) UpdateCategory(id uint, name string) (*domain.UpdateCategoryResponse, error) {
	category, err := u.categoryRepo.UpdateCategory(id, name)
	if err != nil {
		return nil, errors.New("failed to update category")
	}

	catResp := dto.ToCategoryResponse(category)

	return &domain.UpdateCategoryResponse{
		Message:  "Category updated successfully",
		Category: *catResp,
	}, nil
}

func (u *categoryUsecase) DeleteCategory(id uint) (*domain.DeleteCategoryResponse, error) {
	err := u.categoryRepo.DeleteCategory(id)
	if err != nil {
		return nil, errors.New("failed to delete category")
	}
	return &domain.DeleteCategoryResponse{
		Message: "Category deleted successfully",
	}, nil
}

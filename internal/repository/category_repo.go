package repository

import (
	"butik/internal/domain"
	"errors"

	"gorm.io/gorm"
)

type CategoryRepo interface {
	CreateCategory(name string) (*domain.Category, error)
	GetAllCategories(offset, limit int) ([]domain.Category, int, error)
	GetCategoryByID(id uint) (*domain.Category, error)
	UpdateCategory(id uint, name string) (*domain.Category, error)
	DeleteCategory(id uint) error
}

type categoryRepo struct {
	db *gorm.DB
}

func NewCategoryRepo(db *gorm.DB) CategoryRepo {
	return &categoryRepo{db: db}
}

func (r *categoryRepo) CreateCategory(name string) (*domain.Category, error) {
	category := &domain.Category{
		Name: name,
	}
	result := r.db.Create(category)
	if result.Error != nil {
		return nil, errors.New("failed to create category")
	}
	return category, nil
}

func (r *categoryRepo) GetAllCategories(offset, limit int) ([]domain.Category, int, error) {
	var categories []domain.Category
	var total int64

	if err := r.db.Model(&domain.Category{}).Count(&total).Error; err != nil {
		return nil, 0, errors.New("failed to count categories")
	}

	if err := r.db.Order("created_at DESC").Offset(offset).Limit(limit).Find(&categories).Error; err != nil {
		return nil, 0, errors.New("failed to retrieve categories")
	}

	return categories, int(total), nil
}

func (r *categoryRepo) GetCategoryByID(id uint) (*domain.Category, error) {
	category := &domain.Category{}
	result := r.db.First(category, id)
	if result.Error != nil {
		return nil, errors.New("category not found")
	}
	return category, nil
}

func (r *categoryRepo) UpdateCategory(id uint, name string) (*domain.Category, error) {
	category, err := r.GetCategoryByID(id)
	if err != nil {
		return nil, err
	}
	category.Name = name
	result := r.db.Save(category)
	if result.Error != nil {
		return nil, errors.New("failed to update category")
	}
	return category, nil
}

func (r *categoryRepo) DeleteCategory(id uint) error {
	category, err := r.GetCategoryByID(id)
	if err != nil {
		return err
	}
	result := r.db.Delete(category)
	if result.Error != nil {
		return errors.New("failed to delete category")
	}
	return nil
}

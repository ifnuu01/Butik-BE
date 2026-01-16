package repository

import (
	"butik/internal/domain"
	"errors"

	"gorm.io/gorm"
)

type UserRepo interface {
	GetByUsername(username string) (*domain.User, error)
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepo {
	return &userRepo{db: db}
}

func (r *userRepo) GetByUsername(username string) (*domain.User, error) {
	user := &domain.User{}
	result := r.db.Where("username = ?", username).First(user)
	if result.Error != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

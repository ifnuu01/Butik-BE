package usecase

import (
	"butik/internal/domain"
	"butik/internal/infrastructure"
	"butik/internal/repository"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecase interface {
	Login(username, password string) (*domain.LoginResponse, error)
	RefreshToken(refreshToken string) (*domain.RefreshTokenResponse, error)
}

type userUsecase struct {
	userRepo repository.UserRepo
}

func NewUserUsecase(userRepo repository.UserRepo) UserUsecase {
	return &userUsecase{userRepo: userRepo}
}

func (u *userUsecase) Login(username, password string) (*domain.LoginResponse, error) {
	user, err := u.userRepo.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	accessToken, refreshToken, err := infrastructure.CreateTokenPair(user.ID, user.Username)
	if err != nil {
		return nil, errors.New("failed to create tokens")
	}

	return &domain.LoginResponse{
		Message:      "Login successful",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (u *userUsecase) RefreshToken(refreshToken string) (*domain.RefreshTokenResponse, error) {
	accessToken, err := infrastructure.RefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}
	return &domain.RefreshTokenResponse{
		AccessToken: accessToken,
	}, nil
}

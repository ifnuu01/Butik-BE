package main

import (
	"butik/internal/domain"
	"butik/internal/infrastructure"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	infrastructure.LoadEnv()
	db := infrastructure.SetupDB()

	if err := SeedUser(db); err != nil {
		fmt.Println("Failed add seed")
	} else {
		fmt.Println("Seeding user completed")
	}
}

func SeedUser(db *gorm.DB) error {
	password := infrastructure.GetEnv("PASSWORD_ADMIN")
	username := infrastructure.GetEnv("USERNAME_ADMIN")
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := domain.User{
		Username: username,
		Password: string(hashed),
	}

	return db.FirstOrCreate(&user, domain.User{Username: username}).Error
}

package infrastructure

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDB() *gorm.DB {
	dbHost := GetEnv("DB_HOST")
	dbUser := GetEnv("DB_USER")
	dbPassword := GetEnv("DB_PASSWORD")
	dbName := GetEnv("DB_NAME")
	dbPort := GetEnv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Makassar",
		dbHost, dbUser, dbPassword, dbName, dbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connection established")
	return db
}

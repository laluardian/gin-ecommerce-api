package utils

import (
	"fmt"
	"log"

	"github.com/laluardian/gin-ecommerce-api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database")
	}

	db.AutoMigrate(&models.User{})
	fmt.Println("Connected to database")

	return db
}

package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/laluardian/gin-ecommerce-api/routes"
)

func main() {
	// load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	log.Fatal(routes.RunApi())
}

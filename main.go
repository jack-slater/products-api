package main

import (
	"os"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	loadEnvFiles()

	a := App{}
	a.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"))

	a.Run(":8080")
}

func loadEnvFiles() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
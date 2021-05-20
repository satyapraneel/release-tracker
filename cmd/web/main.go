package main

import (
	"github.com/joho/godotenv"
	"github.com/release-trackers/gin/routes"
	"log"
)
func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	routes.RouterGin()
}
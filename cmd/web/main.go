package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/release-trackers/gin/routes"
)

type Handler struct{}

func main() {
	// load .env environment variables
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	r := routes.RouterGin()
	r.Run(os.Getenv("PORT"))
}

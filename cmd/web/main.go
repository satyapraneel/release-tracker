package main

import (
	"github.com/joho/godotenv"
	"github.com/release-trackers/gin/routes"
	"os"
)

type Handler struct{}

func main() {
	// load .env environment variables
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	//bitbucket.GetPr()
	r := routes.RouterGin()
	r.Run(os.Getenv("PORT"))

}

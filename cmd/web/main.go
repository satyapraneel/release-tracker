package main

import (
	"github.com/joho/godotenv"
	"github.com/release-trackers/gin/cmd"
	"github.com/release-trackers/gin/database"
	"github.com/release-trackers/gin/routes"
	"os"
)


func main() {
	// load .env environment variables
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	//bitbucket.GetPr()
	app := &cmd.Application{
		Db:      database.InitConnection(),
		Name: "roopa",
	}
	r := routes.RouterGin(app)
	r.Run(os.Getenv("PORT"))

}

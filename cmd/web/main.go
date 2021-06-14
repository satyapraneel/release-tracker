package main

import (
	"flag"
	"os"

	"github.com/joho/godotenv"
	"github.com/release-trackers/gin/cmd"
	"github.com/release-trackers/gin/database"
	"github.com/release-trackers/gin/routes"
	"github.com/release-trackers/gin/schedular"
)

func main() {

	// defer workers.Wait()
	//sample is given below
	// job := workers.Job{
	// 	Action: workers.PrintPayload,
	// 	Payload: map[string]interface{}{
	// 		"time": 123,
	// 	},
	// }
	// job.Dispatch()
	// load .env environment variables
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	handle()

}

func handle() {
	flag.Parse()
	args := flag.Args()

	if len(args) >= 1 {
		switch args[0] {
		case "seed":
			database.DbSeed()
			os.Exit(0)
		case "migrate":
			database.DbMigrate()
			os.Exit(0)
		}
	}

	//bitbucket.GetPr()
	//bitbucket.CreatePr()
	//cmd.TriggerMail()
	app := &cmd.Application{
		Db:   database.InitConnection(),
		Name: "roopa",
	}
	schedular.Run(app)
	routes.RouterGin(app)
}

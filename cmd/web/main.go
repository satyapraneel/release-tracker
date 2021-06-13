package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/release-trackers/gin/cmd"
	"github.com/release-trackers/gin/database"
	"github.com/release-trackers/gin/routes"
	"github.com/release-trackers/gin/workers"
)

func main() {

	defer workers.Wait()
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
	//bitbucket.GetPr()
	//bitbucket.AuthCredentials(nil)
	//cmd.TriggerMail()
	app := &cmd.Application{
		Db:   database.InitConnection(),
		Name: "roopa",
	}
	r := routes.RouterGin(app)
	r.Run(os.Getenv("PORT"))

}

func main2() {
	req, err := http.NewRequest("GET", "https://bitbucket.org/site/oauth2/authorize?access_type=offline&client_id=YM4nYMmqPQJVy8JxWB&response_type=code&state=state", nil)
	if err != nil {
		panic(err)
	}
	client := new(http.Client)
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return errors.New("Redirect")
	}

	response, err := client.Do(req)
	fmt.Println(response.StatusCode)
	if err != nil {
		if response.StatusCode == 301 { //status code 302
			fmt.Println(response.Location())
		} else {
			panic(err)
		}
	}

}

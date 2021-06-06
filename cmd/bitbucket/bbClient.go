package bitbucket

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/ktrysmt/go-bitbucket"
)

type response struct{}

func GetPr() {
	c := bitbucket.NewBasicAuth("roopa1118", "pGNtuyFzHZjJkyz6yLbK")
	if c != nil {
		log.Println("connecting bitbucket")
	}
	log.Printf("connecting bitbucket %v", c)
	opt := &bitbucket.PullRequestsOptions{
		Owner:             "roopa1118",
		RepoSlug:          "test",
	}

	res, err := c.Repositories.PullRequests.Gets(opt)
	if err != nil {
		panic(err)
	}
	//log.Printf("prs collection %v", res)
	// convert map to json
	jsonString, _ := json.Marshal(res)
	fmt.Println(string(jsonString))

	// convert json to struct
	s := response{}
	json.Unmarshal(jsonString, &s)
	log.Println(s)

	//opt := &bitbucket.RepositoryOptions{
	//	Owner:    "roopa1118",
	//	RepoSlug: "golang",
	//}
	//
	//res, err := c.Repositories.Repository.Get(opt)
	//if err != nil {
	//	log.Print("The repository is not found.")
	//}
	//log.Println("repos ", res.Full_name)
	//if res.Full_name != "roopa1118/golang" {
	//	log.Print("Cannot catch repos full name.")
	//}
}

func CreatePr()  {
	c := bitbucket.NewBasicAuth("roopa1118", "pGNtuyFzHZjJkyz6yLbK")

	opt := &bitbucket.PullRequestsOptions{
		Owner:             "roopa1118",
		RepoSlug:          "test",
		SourceBranch:      "develop",
		DestinationBranch: "master",
		Title:             "fix bug. #9999",
		Description: "test from bitbucket api",
		//CloseSourceBranch: true,
	}

	res, err := c.Repositories.PullRequests.Create(opt)
	if err != nil {
		panic(err)
	}

	log.Printf("create pr %v ",res)
}
package bitbucket

import (
	"log"

	"github.com/ktrysmt/go-bitbucket"
)

func GetPr() {
	c := bitbucket.NewBasicAuth("roopa1118", "pGNtuyFzHZjJkyz6yLbK")
	if c != nil {
		log.Println("connecting bitbucket")
	}
	log.Printf("connecting bitbucket %v", c)
	opt := &bitbucket.PullRequestsOptions{
		Owner:             "roopa1118",
		RepoSlug:          "golang",
	}

	res, err := c.Repositories.PullRequests.Gets(opt)
	log.Printf("%v result", res)
	if err != nil {
		panic(err)
	}
	log.Println(res)
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
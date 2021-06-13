package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ktrysmt/go-bitbucket"
	"log"
	"net/http"
)

type response struct{}

type Branch struct{
	name string
	target struct{
		hash string
	}
}


func GetPr() {
	//c := bitbucket.NewBasicAuth("roopa1118", "Ewmc3mBtH4wMUHpxzZBP")  //pGNtuyFzHZjJkyz6yLbK
	//if c != nil {
	//	log.Println("connecting bitbucket")
	//}
	//log.Printf("connecting bitbucket %v", c)
	//opt := &bitbucket.PullRequestsOptions{
	//	Owner:             "roopa1118",
	//	RepoSlug:          "test",
	//}
	//
	//res, err := c.Repositories.PullRequests.Gets(opt)
	//if err != nil {
	//	panic(err)
	//}
	////log.Printf("prs collection %v", res)
	//// convert map to json
	//jsonString, _ := json.Marshal(res)
	//fmt.Println(string(jsonString))
	//
	//// convert json to struct
	//s := response{}
	//json.Unmarshal(jsonString, &s)
	//log.Println(s)

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

func CreateBasicAuth(){
	clientId := "YM4nYMmqPQJVy8JxWB"
	clientSecret := "vsaQdwxEdCNCdmpncehMcC5K68mQFh9R"
	c := bitbucket.NewOAuth(clientId, clientSecret)  //pGNtuyFzHZjJkyz6yLbK
	fmt.Printf("%+v", c.Auth)
	if c != nil {
		log.Println("connecting bitbucket")
	}
	log.Printf("connecting bitbucket %v", c)
}

func CreatePr(token string)  {
	//clientId := "YM4nYMmqPQJVy8JxWB"
	//clientSecret := "vsaQdwxEdCNCdmpncehMcC5K68mQFh9R"
	//cxt := bitbucket.NewOAuth(clientId, clientSecret)
	////fmt.Println("access token value : ", cxt.Auth)
	//fmt.Printf("%+v", cxt.Auth)

	//cxt = bitbucket.NewOAuthbearerToken("F8B9zc-iudSCQ-IhquodQysqp_62Z5-zrKnJuiFPTIv7AcrXziCOXRRRdNGNxkk5YNCYldub6DoisxJo7tb5X-WzP00dq5BTcGn9m1SKDNL32zCunLJrq525")
	//cxt, accessToken := bitbucket.NewOAuthWithCode(clientId, clientSecret, code)

	//if accessToken != ""{
	//	fmt.Println("access token value : ", accessToken)
	//}

	//bitbucket.BranchRestrictionsOptions

	//fmt.Println("request data", cxt)
	cxt := bitbucket.NewOAuthbearerToken(token)
	request := &Branch{name: "develop_roopa", target: struct{ hash string }{hash: "b678d75ac143df62c6371d2f1cd01ea6d3585d31"}}
	apiUrl := cxt.GetApiBaseURL()+"/repositories/roopajoshyam/test-pro/refs/branches"
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(request)
	fmt.Println("request data", request)
	res, err := cxt.HttpClient.Post(apiUrl, "application/json", payloadBuf)
	log.Printf("response : %+v", res)
	if err != nil {
		fmt.Println("error :",err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	createdBranch := new(Branch)
	if err = json.NewDecoder(res.Body).Decode(createdBranch); err != nil {
		log.Print(err)
	}

	log.Printf("branch response %v : ",createdBranch)
}
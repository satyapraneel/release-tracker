package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/providers/bitbucket"
	"log"
	"net/http"
	"net/url"
	"os"
)
type Sessions struct {
	*bitbucket.Session
}
type Payload struct {
	Name   string `json:"name"`
	Target Target `json:"target"`
}
type Target struct {
	Hash string `json:"hash"`
}

type BranchRestriction struct{
	Owner    string            `json:"owner"`
	RepoSlug string            `json:"repo_slug"`
	Kind     string            `json:"kind"`
	FullSlug string            `json:"full_slug"`
	Name     string            `json:"name"`
	Users    []Users          	`json:"users"`
	Pattern  string				`json:"pattern"`
}
type Users struct {
	Username string `json:"username"`
}

const callbackUrl = "http://localhost:4000/oauth/index"
const baseURL  =  "https://api.bitbucket.org/2.0"
const branchCreationUrl = "/repositories/roopajoshyam/test-pro/refs/branches"
const branchRestriction = "/repositories/roopajoshyam/test-pro/branch-restrictions"

func Authrorize() string {
	provider := bitbucketProvider()
	session, err := provider.BeginAuth("test")
	s := session.(*bitbucket.Session)
	if err != nil {
		log.Println(err)
	}
	return s.AuthURL
}

func GetAccessToken(token_code string) *bitbucket.Session {
	provider := bitbucketProvider()
	se := &bitbucket.Session{}
	urlParams := url.Values{}
	urlParams.Add("code", token_code)

	//Get AccessToken from Authorization code
	token, tokenErr := se.Authorize(provider, urlParams)
	if tokenErr != nil {
		fmt.Println(tokenErr)
		refreshToken, _ := provider.RefreshToken(se.RefreshToken)
		token = refreshToken.AccessToken
	}
	fmt.Printf("session storeed: %+v \n", se)
	fmt.Printf("session token: %+v \n", token)
	return se
}

func CreateBranch(c *gin.Context, branchType string, name string, reviewers []string) {
	session := sessions.Default(c)
	AccessToken :=  fmt.Sprintf("%v", session.Get("access_token"))
	log.Printf("acess ---- %+v : ", AccessToken)
	branch := branchType+"/"+name
	request := Payload{Name: branch, Target: Target{Hash: "b678d75ac143df62c6371d2f1cd01ea6d3585d31"}}
	apiUrl := baseURL+branchCreationUrl
	payloadBytes, _ := json.Marshal(request)
	body := bytes.NewReader(payloadBytes)
	resp := PostRequest(apiUrl, body, AccessToken)
	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		fmt.Errorf("unknown error, status code: %d", resp.StatusCode)
	}
	createdBranch := new(Payload)
	errs := json.NewDecoder(resp.Body).Decode(createdBranch)
	if errs != nil {
		log.Print(errs)
	}
	log.Printf("branch Name %+v : ", createdBranch)
	branchRestrictions(AccessToken, branch, reviewers)
}

func branchRestrictions(token string, branchName string, ReviewerList []string)  {
	apiUrl := baseURL+branchRestriction
	var arrayOfUsers  []Users
	for _, reviewer := range ReviewerList {
		user := Users{Username: reviewer}
		arrayOfUsers = append(arrayOfUsers, user)
	}
	request := BranchRestriction{
		Kind: "restrict_merges",
		Owner: "roopajoshyam",
		RepoSlug: "test-pro",
		Pattern: "*"+branchName+"*",
		Users: arrayOfUsers,
	}
	payloadBytes, _ := json.Marshal(request)
	body := bytes.NewReader(payloadBytes)
	res := PostRequest(apiUrl, body, token)
	defer res.Body.Close()
	fmt.Printf("branch restrcition %+v", res)
	if res.StatusCode != 201 {
		fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}
}

func PostRequest(apiUrl string, body *bytes.Reader, token string) *http.Response {
	fmt.Printf("**Access token: %+v \n", token)
	req, err := http.NewRequest("POST", apiUrl, body)
	if err != nil {
		fmt.Println("error :", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("error 2:", err)
	}
	return resp
}


func bitbucketProvider() *bitbucket.Provider {
	return bitbucket.New(os.Getenv("BITBUCKET_KEY"), os.Getenv("BITBUCKET_SECRET"), callbackUrl, "account:write")
}
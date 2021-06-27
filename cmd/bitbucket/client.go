package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/providers/bitbucket"
	"github.com/release-trackers/gin/models"
)

//type Sessions struct {
//	*bitbucket.Session
//}
type Payload struct {
	Name   string `json:"name"`
	Target Target `json:"target"`
}
type Target struct {
	Hash string `json:"hash"`
}

type GetBody struct{}

type BBAccessToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type BranchRestriction struct {
	Owner    string  `json:"owner"`
	RepoSlug string  `json:"repo_slug"`
	Kind     string  `json:"kind"`
	FullSlug string  `json:"full_slug"`
	Name     string  `json:"name"`
	Users    []Users `json:"users"`
	Pattern  string  `json:"pattern"`
}
type Users struct {
	Username string `json:"username"`
}

const callbackUrl = "http://localhost:4000/oauth/index"
const baseURL = "https://api.bitbucket.org/2.0"

func Authrorize() string {
	provider := bitbucketProvider()
	session, err := provider.BeginAuth("test")
	s := session.(*bitbucket.Session)
	if err != nil {
		log.Println(err)
	}
	return s.AuthURL
}

func TestGetAccessToken(token_code string) *bitbucket.Session {
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

func GetAccessToken(c *gin.Context)  {
	session := sessions.Default(c)
	accessToken := session.Get("access_token")
	if accessToken == nil {
		params := url.Values{}
		params.Add("grant_type", `client_credentials`)
		body := strings.NewReader(params.Encode())

		req, err := http.NewRequest("POST", "https://bitbucket.org/site/oauth2/access_token", body)
		if err != nil {
			// handle err
		}
		req.SetBasicAuth(os.Getenv("BITBUCKET_KEY"), os.Getenv("BITBUCKET_SECRET"))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			// handle err
		}
		defer resp.Body.Close()
		access := new(BBAccessToken)
		errs := json.NewDecoder(resp.Body).Decode(access)
		if errs != nil {
			log.Print(errs)
		}
		log.Printf("access token %+v : ", access.AccessToken)
		//newT := currentTime.Add(time.Second * time.Duration(access.ExpiresIn))
		//gob.Register(time.Time{})
		sess := sessions.Default(c)
		sess.Set("access_token", access.AccessToken)
		if err := sess.Save(); err != nil {
			log.Print(err)
			return
		}
	}
	fmt.Printf("access in\n: %s\n", session.Get("access_token"))

}

func CreateBranch(c *gin.Context, branchType string, name string, reviewers []string, projectRepoName string) {
	session := sessions.Default(c)
	AccessToken :=  fmt.Sprintf("%v", session.Get("access_token"))
	branch := branchType+"/"+name
	request := Payload{Name: branch, Target: Target{Hash: "master"}}
	branchCreationUrl := "/repositories/"+os.Getenv("BITBUCKET_OWNER")+"/"+projectRepoName+"/refs/branches"
	apiUrl := baseURL+branchCreationUrl
	log.Printf("API URL ---- %+v : ", apiUrl)
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
	branchRestrictions(AccessToken, branch, reviewers, projectRepoName)

}

func branchRestrictions(token string, branchName string, ReviewerList []string, projectRepoName string)  {
	branchRestriction := "/repositories/"+os.Getenv("BITBUCKET_OWNER")+"/"+projectRepoName+"/branch-restrictions"
	apiUrl := baseURL+branchRestriction
	var arrayOfUsers  []Users
	for _, reviewer := range ReviewerList {
		user := Users{Username: reviewer}
		arrayOfUsers = append(arrayOfUsers, user)
	}
	request := BranchRestriction{
		Kind: "restrict_merges",
		Owner: os.Getenv("BITBUCKET_OWNER"),
		RepoSlug: projectRepoName,
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

func PostRequest(apiUrl string, method string, body *bytes.Reader, token string) *http.Response {
	fmt.Printf("**Access token: %+v \n", token)
	req, err := http.NewRequest(method, apiUrl, body)
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

func GetReleseIssuesIds(c *gin.Context, release models.Release, project models.Project) []string {
	var releaseIds []string
	session := sessions.Default(c)
	AccessToken := fmt.Sprintf("%v", session.Get("access_token"))
	getCommits := "/repositories/" + os.Getenv("BITBUCKET_OWNER") + "/" + project.Name + "/commits/" + release.Type + "/" + release.Name + "?exclude=master"
	apiUrl := baseURL + getCommits
	request := GetBody{}
	payloadBytes, _ := json.Marshal(request)
	body := bytes.NewReader(payloadBytes)
	resp := PostRequest(apiUrl, "GET", body, AccessToken)
	defer resp.Body.Close()
	responseCoomits, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(responseCoomits))
	return releaseIds
}

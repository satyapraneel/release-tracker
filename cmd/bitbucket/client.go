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
	"regexp"
	"strings"

	"github.com/markbates/goth/providers/bitbucket"
	"github.com/release-trackers/gin/models"
	"gorm.io/gorm"
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

type BBAccessToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type GetBody struct{}

type CommitResonse struct {
	PageLen int       `json:"pagelen"`
	Commits []Commits `json:"values"`
	Next    string    `json:"next"`
}

type Commits struct {
	Message string `json:"message"`
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

type UpdateRestriction struct {
	Values []*BranchRestriction
}

type Users struct {
	Username string `json:"username"`
}

type BranchRestrictionResponse struct {
	Id int
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

func GetAccessToken() string {
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

	return access.AccessToken
}

func CreateBranch(db *gorm.DB, release models.Release, reviewers []string, projectRepoName string) {
	AccessToken := GetAccessToken()
	branch := release.Type + "/" + release.Name
	request := Payload{Name: branch, Target: Target{Hash: "master"}}
	branchCreationUrl := "/repositories/" + os.Getenv("BITBUCKET_OWNER") + "/" + projectRepoName + "/refs/branches"
	apiUrl := baseURL + branchCreationUrl
	log.Printf("API URL ---- %+v : ", apiUrl)
	log.Printf("Branch URL ---- %+v : ", branch)
	payloadBytes, _ := json.Marshal(request)
	body := bytes.NewReader(payloadBytes)
	resp := PostRequest(apiUrl, body, AccessToken, "POST")
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
	restrictionPushId := branchRestrictions(AccessToken, branch, reviewers, projectRepoName, "push")
	restrictionMergeId := branchRestrictions(AccessToken, branch, reviewers, projectRepoName, "restrict_merges")
	log.Printf("restriction push id: %v ", restrictionPushId)
	releaseRestriction := db.Model(&release).Updates(map[string]interface{}{"restriction_push_id": restrictionPushId, "restriction_merge_id": restrictionMergeId})
	log.Printf("restriction error: %v ", releaseRestriction.Error)
	if releaseRestriction.Error != nil {
		log.Print(releaseRestriction.Error)
	}
}

func branchRestrictions(token string, branchName string, ReviewerList []string, projectRepoName string, restrictionKind string) int {
	branchRestriction := "/repositories/" + os.Getenv("BITBUCKET_OWNER") + "/" + projectRepoName + "/branch-restrictions"
	apiUrl := baseURL + branchRestriction
	var arrayOfUsers []Users
	for _, reviewer := range ReviewerList {
		user := Users{Username: reviewer}
		arrayOfUsers = append(arrayOfUsers, user)
	}
	request := BranchRestriction{
		Kind:     restrictionKind,
		Owner:    os.Getenv("BITBUCKET_OWNER"),
		RepoSlug: projectRepoName,
		Pattern:  "*" + branchName + "*",
		Users:    arrayOfUsers,
	}
	payloadBytes, _ := json.Marshal(request)
	body := bytes.NewReader(payloadBytes)
	res := PostRequest(apiUrl, body, token, "POST")
	defer res.Body.Close()
	fmt.Printf("branch restrcition %+v", res)
	if res.StatusCode != 201 {
		fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}
	restriction := new(BranchRestrictionResponse)
	errs := json.NewDecoder(res.Body).Decode(restriction)
	if errs != nil {
		log.Print(errs)
	}

	return restriction.Id
}

func UpdateBranchRestriction(projectRepoName string, restrictionId string, branchName string, restrictionKind string) {
	AccessToken := GetAccessToken()
	branchRestriction := "/repositories/" + os.Getenv("BITBUCKET_OWNER") + "/" + projectRepoName + "/branch-restrictions/" + restrictionId
	apiUrl := baseURL + branchRestriction
	var arrayOfUsers []Users
	request := &BranchRestriction{
		Kind:     restrictionKind,
		Owner:    os.Getenv("BITBUCKET_OWNER"),
		RepoSlug: projectRepoName,
		Pattern:  "*" + branchName + "*",
		Users:    arrayOfUsers,
	}
	payloadBytes, _ := json.Marshal(request)
	body := bytes.NewReader(payloadBytes)
	res := PostRequest(apiUrl, body, AccessToken, "PUT")
	defer res.Body.Close()
	fmt.Printf("branch update restrcition %+v", res)
	if res.StatusCode != 200 {
		fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}
}

func PostRequest(apiUrl string, body *bytes.Reader, token string, reqType string) *http.Response {
	fmt.Printf("**Access token: %+v \n", token)
	req, err := http.NewRequest(reqType, apiUrl, body)
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

func GetReleseIssuesIds(release models.Release, project models.Project) []string {
	AccessToken := GetAccessToken()
	getCommits := "/repositories/" + os.Getenv("BITBUCKET_OWNER") + "/" + project.RepoName + "/commits/" + release.Type + "/" + release.Name + "?exclude=master"
	apiUrl := baseURL + getCommits
	request := GetBody{}
	payloadBytes, _ := json.Marshal(request)
	body := bytes.NewReader(payloadBytes)
	resp := PostRequest(apiUrl, body, AccessToken, "GET")
	defer resp.Body.Close()
	responseCommits, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var commitResonse CommitResonse
	json.Unmarshal(responseCommits, &commitResonse)
	if commitResonse.Next != "" {
		fmt.Print("not nil")
	} else {
		fmt.Print(" nil")
	}
	return getRelatedIdsFromCommit(commitResonse, project.RelatedCodes)
}

func getRelatedIdsFromCommit(commitResonse CommitResonse, relatedCodes string) []string {
	var releaseIds []string
	codes := strings.Split(relatedCodes, ",")
	for _, commit := range commitResonse.Commits {
		for _, code := range codes {
			r := regexp.MustCompile(code + "-[0-9]+")
			releaseIds = append(releaseIds, r.FindAllString(commit.Message, -1)...)
		}
	}
	return removeDuplicateValues(releaseIds)
}

func removeDuplicateValues(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

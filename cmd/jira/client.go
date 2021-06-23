package jira

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"log"
	"os"
	"strings"
	"time"
)

func setUpClient() (*jira.Client, error) {
	base := os.Getenv("JIRA_BASE_URL")
	tp := jira.BasicAuthTransport{
		Username: os.Getenv("JIRA_USERNAME"),  //"jroopanov11@gmail.com",
		Password: os.Getenv("JIRA_SECRET"),   //"MIjRCG1iztxW6yU8Xk754F98",
	}

	jiraClient, err := jira.NewClient(tp.Client(), base)
	return jiraClient,err
}

func GetIssueByName(){
	jiraClient, err := setUpClient()
	req, _ := jiraClient.NewRequest("GET", "rest/api/2/issue/LOYAL-4643", nil)

	issue := new(jira.Issue)
	_, err = jiraClient.Do(req, issue)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s: %+v\n", issue.Key, issue.Fields.Summary)
	fmt.Printf("Type: %s\n", issue.Fields.Type.Name)
	fmt.Printf("Priority: %s\n", issue.Fields.Priority.Name)
}

func GetIssuesByLabel() {
	jiraClient, err := setUpClient()
	var issues []jira.Issue

	// appendFunc will append jira issues to []jira.Issue
	appendFunc := func(i jira.Issue) (err error) {
		issues = append(issues, i)
		return err
	}

	err = jiraClient.Issue.SearchPages(fmt.Sprintf(`labels IN (%s, %s)`, strings.TrimSpace("release01"),strings.TrimSpace("example02")), nil, appendFunc)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d issues found.\n", len(issues))

	for _, i := range issues {
		t := time.Time(i.Fields.Created) // convert go-jira.Time to time.Time for manipulation
		date := t.Format("2006-01-02")
		clock := t.Format("15:04")
		fmt.Printf("Creation Date: %s\nCreation Time: %s\nIssue Key: %s\nIssue Summary: %s\n\n", date, clock, i.Key, i.Fields.Summary)
	}
}


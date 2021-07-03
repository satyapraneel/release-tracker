package jira

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"log"
	"os"
	"strings"
	"time"
)

type JiraTickets struct {
	Id string
	Summary string
	CreationDate string
	CreationTime string
	Type string
	Project string
	Priority string
	Status string
}

func setUpClient() (*jira.Client, error) {
	base := os.Getenv("JIRA_BASE_URL")
	tp := jira.BasicAuthTransport{
		Username: os.Getenv("JIRA_USERNAME"),  //"jroopanov11@gmail.com",
		Password: os.Getenv("JIRA_SECRET"),   //"MIjRCG1iztxW6yU8Xk754F98",
	}

	jiraClient, err := jira.NewClient(tp.Client(), base)
	return jiraClient,err
}

func GetIssueDetails(issueKey string) *JiraTickets {
	jiraClient, err := setUpClient()
	req, _ := jiraClient.NewRequest("GET", "rest/api/2/issue/"+issueKey, nil)

	issue := new(jira.Issue)
	_, err = jiraClient.Do(req, issue)
	if err != nil {
		panic(err)
	}


	t := time.Time(issue.Fields.Created) // convert go-jira.Time to time.Time for manipulation
	date := t.Format("2006-01-02")
	clock := t.Format("15:04")
	jiraInfo := &JiraTickets{Id: issue.Key, Summary: issue.Fields.Summary, CreationDate: date, CreationTime: clock,
		Type: issue.Fields.Type.Name, Project: issue.Fields.Project.Name, Priority: issue.Fields.Priority.Name,
		Status: issue.Fields.Status.Name,
	}

	return jiraInfo
}

func GetIssuesByLabel(releaseName string) []*JiraTickets {
	jiraClient, err := setUpClient()
	var issues []jira.Issue

	// appendFunc will append jira issues to []jira.Issue
	appendFunc := func(i jira.Issue) (err error) {
		issues = append(issues, i)
		return err
	}

	err = jiraClient.Issue.SearchPages(fmt.Sprintf(`labels IN (%s)`, strings.TrimSpace(releaseName)), nil, appendFunc)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d issues found.\n", len(issues))
	jiraArr := []*JiraTickets{}
	for _, i := range issues {
		t := time.Time(i.Fields.Created) // convert go-jira.Time to time.Time for manipulation
		date := t.Format("2006-01-02")
		clock := t.Format("15:04")
		jiraInfo := &JiraTickets{Id: i.Key, Summary: i.Fields.Summary, CreationDate: date, CreationTime: clock,
			Type: i.Fields.Type.Name, Project: i.Fields.Project.Name, Priority: i.Fields.Priority.Name,
			Status: i.Fields.Status.Name,
		}

		jiraArr = append(jiraArr, jiraInfo)
		fmt.Printf("info of jira ticker %+v", i)
		fmt.Printf("Creation Date: %s\nCreation Time: %s\nIssue Key: %s\nIssue Summary: %s\n\n", date, clock, i.Key, i.Fields.Summary)
	}

	return jiraArr
}


package mails

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/release-trackers/gin/models"
)

type ReleaseCreateData struct {
	ProjectName string
	ReleaseName string
}

type ReleaseNotesDate struct {
	ReleaseTickets []*models.ReleaseTickets
	Release        *models.Release
	ReleaseDate    string
	JiraUrl      string
}

func SendReleaseCreatedMail(release *models.Release, project *models.Project) {
	// tagetDate := &release.TargetDate
	subject := "Jira Label for " + project.Name + "!"
	reviews := strings.Split(project.ReviewerList, ",")
	mail := NewMail(reviews, subject, "", "")
	releaseData := &ReleaseCreateData{ProjectName: project.Name, ReleaseName: release.Name}
	errs := mail.ParseTemplate("/ui/html/mails/release_create.html", releaseData)
	if errs != nil {
		log.Printf("template parse : %v", errs)
	}
	ok, err := mail.SendEmail()
	fmt.Println("err in sending mail")
	fmt.Println(err)
	fmt.Println(ok)
}

func SendReleaseNotes(release *models.Release, releaseTickets []*models.ReleaseTickets) (bool, error) {
	subject := "Release Notes for " + release.Name
	reviews := strings.Split(release.Owner, ",")
	mail := NewMail(reviews, subject, "", "")
	jiraBrowseUrl := os.Getenv("JIRA_BASE_URL") + "browse"
	tickets := &ReleaseNotesDate{ReleaseTickets: releaseTickets, Release: release, ReleaseDate: release.TargetDate.Format("2006-01-02"), JiraUrl: jiraBrowseUrl}
	errs := mail.ParseTemplate("/ui/html/mails/release_notes.html", tickets)
	if errs != nil {
		log.Printf("template parse : %v", errs)
	}
	ok, err := mail.SendEmail()
	fmt.Println(ok)
	return ok, err
}

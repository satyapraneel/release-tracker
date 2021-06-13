package mails

import (
	"fmt"
	"log"
	"strings"

	"github.com/release-trackers/gin/models"
)

type ReleaseCreateData struct {
	ProjectName string
	ReleaseName string
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
	ok, _ := mail.SendEmail()
	fmt.Println(ok)
}

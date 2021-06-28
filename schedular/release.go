package schedular

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/release-trackers/gin/cmd/bitbucket"
	"github.com/release-trackers/gin/cmd/jira"
	"github.com/release-trackers/gin/models"
	"github.com/release-trackers/gin/notifications/mails"
	"github.com/release-trackers/gin/repositories"
)

func (app *Application) ReleaseDateReminder() {
	//
	release := repositories.NewReleaseHandler(app.Application)
	releases, err := release.GetLatestReleases()
	if err != nil {
		log.Println(err.Error())
		return
	}
	//if current date is beta release date send reminder
	for _, releaseRecord := range releases {
		projects, _, err := release.GetReleaseProjects(releaseRecord)
		if err != nil {
			log.Println(err.Error())
			return
		}
		for _, project := range projects {
			app.TriggerReminderMail(releaseRecord, project)

		}
	}
}

func (app *Application) TriggerReminderMail(releaseRecord models.Release, project *models.Project) {
	app.getDate(project.BetaReleaseDate, project, &releaseRecord, "beta")
	app.getDate(project.RegressionSignorDate, project, &releaseRecord, "regression")
	app.getDate(project.CodeFreezeDate, project, &releaseRecord, "code_freeze")
	app.getDate(project.DevCompletionDate, project, &releaseRecord, "dev_completion")
	app.getDate("1", project, &releaseRecord, "release_date")
}

func (app *Application) getDate(days string, project *models.Project, releaseRecord *models.Release, typeOfRelease string) {
	daysToSubtract, err := strconv.Atoi(days)
	if err != nil {
		log.Println(err)
	}
	releaseDate := releaseRecord.TargetDate.AddDate(0, 0, -daysToSubtract).Truncate(24 * time.Hour)
	today := time.Now().Truncate(24 * time.Hour)
	if today.Equal(releaseDate) {
		app.TriggerMailIfDate(typeOfRelease, project, releaseRecord)
	}
}

func (app *Application) TriggerMailIfDate(typeOfRelease string, project *models.Project, releaseRecord *models.Release) {
	dlsRepository := repositories.NewReleaseHandler(app.Application)
	var dlType string
	var releaseType string
	switch typeOfRelease {
	case "beta":
		dlType = models.DEV
		releaseType = "Beta"
	case "regression":
		dlType = models.DEV
		releaseType = "Regression"
		log.Println("Regression")
		//get jira tickets
		jiraArr := jira.GetIssuesByLabel(releaseRecord.Name)
		for _, i := range jiraArr {
			fmt.Printf("%+v", i)
		}
	case "code_freeze":
		dlType = models.QA
		releaseType = "Code Freeze"
		log.Println("its freeze")
		token := bitbucket.GetAccessToken()
		branchName := releaseRecord.Type + "/" + releaseRecord.Name
		var reviewers []string
		bitbucket.BranchRestrictions(token, branchName, reviewers, project.RepoName)
	case "dev_completion":
		dlType = models.DEV
		releaseType = "Dev Completion"
		log.Println("its devcompletion")
		//get the project
		//get the dls list for project depending on type of release
	case "release_date":
		dlType = models.DEV
		releaseType = "Release Date"
	}
	dlsList, _ := dlsRepository.GetDLsByProject(project.ID, dlType)
	//if current date is beta release date send reminder
	for _, dls := range dlsList {
		mails.SendReminderMail(project, dls, releaseType)
	}
}

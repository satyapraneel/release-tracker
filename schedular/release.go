package schedular

import (
	"fmt"
	"log"
	"os"
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
	release.CloseRelease()
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
	app.getDate(project.BetaReleaseDate, project, &releaseRecord, "beta", 1)
	app.getDate(project.RegressionSignorDate, project, &releaseRecord, "regression", 0)
	app.getDate(project.CodeFreezeDate, project, &releaseRecord, "code_freeze", 0)
	app.getDate(project.CodeFreezeDate, project, &releaseRecord, "code_freeze_restrict_branch", -1) //restrict branch next day
	app.getDate(project.DevCompletionDate, project, &releaseRecord, "dev_completion", 0)
	app.getDate("1", project, &releaseRecord, "release_date", 0)
}

func (app *Application) getDate(days string, project *models.Project, releaseRecord *models.Release, typeOfRelease string, subtractor int) {
	daysToSubtract, err := strconv.Atoi(days)
	if err != nil {
		log.Println(err)
		return
	}
	//if we want to send a reminder before release days Ex: send reminder before beta release
	daysToSubtract = daysToSubtract - subtractor
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
	var template string
	dataTosend := new(mails.MailData)
	switch typeOfRelease {
	//tested
	case "beta":
		dlType = models.DEV
		releaseType = "Beta"
		log.Println("its beta")
		subject := releaseRecord.Name + ": Reminder mail for " + releaseType
		template = "/ui/html/mails/release_beta.html"
		dataTosend = &mails.MailData{
			ProjectName:  project.Name,
			ReminderType: releaseType,
			Subject:      subject,
			DLType:       dlType,
		}
	//tested
	case "regression":
		dlType = models.DEV
		releaseType = "Regression"
		log.Println("its regression")
		subject := releaseRecord.Name + ": Regression Sign-off reminder for " + releaseRecord.Name
		template = "/ui/html/mails/regression.html"
		jiraList := app.getRegressionData(releaseRecord, dlsRepository)
		jiraBrowseUrl := os.Getenv("JIRA_BASE_URL") + "browse"
		dataTosend = &mails.MailData{
			JiraTickets:  jiraList,
			ProjectName:  project.Name,
			ReminderType: releaseType,
			Subject:      subject,
			Release:      releaseRecord,
			DLType:       dlType,
			JiraUrl:      jiraBrowseUrl,
		}
	//tested
	case "code_freeze":
		dlType = models.DEV
		releaseType = "Code Freeze"
		log.Println("its code freeze")
		subject := releaseRecord.Name + ": Reminder mail for " + releaseType
		template = "/ui/html/mails/code_freeze.html"
		dataTosend = &mails.MailData{
			ProjectName:  project.Name,
			ReminderType: releaseType,
			Subject:      subject,
			DLType:       dlType,
		}

	case "code_freeze_restrict_branch":
		log.Println("its branch restriction")
		restrictionId := fmt.Sprint(releaseRecord.RestrictionPushId)
		restrictionMergeId := fmt.Sprint(releaseRecord.RestrictionMergeId)
		bitbucket.UpdateBranchRestriction(project.RepoName, restrictionId, releaseRecord.Name, "push")
		bitbucket.UpdateBranchRestriction(project.RepoName, restrictionMergeId, releaseRecord.Name, "restrict_merges")
		dlType = models.DEV
		releaseType = "Code Freeze"
		subject := releaseRecord.Name + ": Reminder mail for " + releaseType
		template = "/ui/html/mails/code_freeze_restrict_branch.html"
		dataTosend = &mails.MailData{
			ProjectName:  project.Name,
			ReminderType: releaseType,
			Subject:      subject,
			DLType:       dlType,
			Release:      releaseRecord,
		}
	//tested
	case "dev_completion":
		dlType = models.DEV
		releaseType = "Dev Completion"
		log.Println("its devcompletion")
		subject := releaseRecord.Name + ": Reminder mail for " + releaseType
		template = "/ui/html/mails/dev_completion.html"
		dataTosend = &mails.MailData{
			ProjectName:  project.Name,
			ReminderType: releaseType,
			Subject:      subject,
			DLType:       dlType,
		}
	//tested
	case "release_date":
		log.Println("its release date")
		dlType = models.DEV
		releaseType = "Release Date"
		subject := releaseRecord.Name + ": Reminder mail for " + releaseType
		template = "/ui/html/mails/release_reminder.html"
		dataTosend = &mails.MailData{
			ProjectName:  project.Name,
			ReminderType: releaseType,
			Subject:      subject,
			DLType:       dlType,
		}
	}
	dlsList, _ := dlsRepository.GetDLsByProject(project.ID)
	//if current date is beta release date send reminder
	go mails.SendReminderMail(dlsList, dataTosend, template)
}

func (app *Application) getRegressionData(releaseRecord *models.Release, repository *repositories.App) []*jira.JiraTickets {
	jirsList := jira.GetIssuesByLabel(releaseRecord.Name)
	log.Printf("regression Data %+v", jirsList)
	for _, tick := range jirsList {
		log.Printf("regression Data 2 %+v", tick.Id)

	}
	go repository.UpdateJiraTicketsToDB(jirsList, releaseRecord.ID)

	return jirsList
}

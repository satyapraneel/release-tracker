package schedular

import (
	"log"
	"strconv"
	"time"

	"github.com/release-trackers/gin/models"
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
	app.getDate(project.BetaReleaseDate, &releaseRecord, "beta")
	app.getDate(project.RegressionSignorDate, &releaseRecord, "regression")
	app.getDate(project.CodeFreezeDate, &releaseRecord, "code_freeze")
	app.getDate(project.DevCompletionDate, &releaseRecord, "dev_completion")
}

func (app *Application) getDate(days string, releaseRecord *models.Release, typeOfRelease string) {
	daysToSubtract, err := strconv.Atoi(days)
	if err != nil {
		log.Println(err)
	}
	releaseDate := releaseRecord.TargetDate.AddDate(0, 0, -daysToSubtract).Truncate(24 * time.Hour)
	today := time.Now().Truncate(24 * time.Hour)
	if today.Equal(releaseDate) {
		app.TriggerMailIfDate(typeOfRelease, releaseRecord)
	}
}

func (app *Application) TriggerMailIfDate(typeOfRelease string, releaseRecord *models.Release) {
	switch typeOfRelease {
	case "beta":
		log.Println("its working")
	case "regression":
		log.Println("Regression")
	case "code_freeze":
		log.Println("its freeze")
	case "dev_completion":
		log.Println("its devcompletion")
	}
}

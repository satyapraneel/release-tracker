package schedular

import (
	"log"
	"strconv"
	"time"

	"github.com/release-trackers/gin/repositories"
)

func (app *Application) ReleaseDateReminder() {
	//
	release := repositories.NewReleaseHandler(app.Application)
	releases, projects, _, err := release.GetLatestReleases()
	if err != nil {
		log.Println(err.Error())
		return
	}
	//if current date is beta release date send reminder

	for _, project := range projects {
		log.Printf("project in loop : %v", project.Name)
		beta, err := strconv.Atoi(project.BetaReleaseDate)
		if err != nil {
			log.Println(err.Error())
			return
		}
		betaReleaseDate := releases.TargetDate.AddDate(0, 0, beta).Truncate(24 * time.Hour)
		today := time.Now().Truncate(24 * time.Hour)
		if today.Equal(betaReleaseDate) {

		}
	}
}

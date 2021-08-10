package seed

import (
	"github.com/jinzhu/gorm"
	"github.com/release-trackers/gin/models"
)

func CreateProject(db *gorm.DB, name string,
	reponame string, reviewers string,
	betaReleaseDate string, regressionDate string,
	codeFreezedate string, devComplDate string, status string,
) error {
	return db.Create(&models.
		Project{
		Name:                 name,
		RepoName:             reponame,
		ReviewerList:         reviewers,
		BetaReleaseDate:      betaReleaseDate,
		RegressionSignorDate: regressionDate,
		CodeFreezeDate:       codeFreezedate,
		DevCompletionDate:    devComplDate,
		Status:               status,
	}).Error
}

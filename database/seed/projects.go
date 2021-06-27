package seed

import (
	"github.com/release-trackers/gin/models"
	"gorm.io/gorm"
)

func CreateProject(db *gorm.DB, name string,
	reponame string, reviewers string,
	betaReleaseDate string, regressionDate string,
	codeFreezedate string, devComplDate string, relatedCodes string, status string,
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
		RelatedCodes:         relatedCodes,
		Status:               status,
	}).Error
}

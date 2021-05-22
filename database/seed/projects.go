package seed

import (
	"github.com/release-trackers/gin/models"
	"gorm.io/gorm"
)

func CreateProject(db *gorm.DB, name string,
	bitbucketUrl string, reviewers string,
	betaReleaseDate string, regressionDate string,
	codeFreezedate string, devComplDate string,
	)  error {
	return db.Create(&models.
		Project{
			Name: name,
			BitbucketUrl: bitbucketUrl,
			ReviewerList: reviewers,
			BetaReleaseDate: betaReleaseDate,
			RegressionSignorDate: regressionDate,
			CodeFreezeDate: codeFreezedate,
			DevCompletionDate: devComplDate,
		}).Error
}

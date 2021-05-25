package repositories

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/release-trackers/gin/database"
	"github.com/release-trackers/gin/models"
	"log"
)

var (
	errInvalidBody     = errors.New("Invalid request body")
	errNotExist			= errors.New("No records found")
	db = database.InitConnection()
)
func CreateRelease(c *gin.Context, release models.Release) (uint, error){
	// Get DB from Mysql Config
	//release := models.Release{}
	err := c.Bind(&release)
	if err != nil {
		log.Print(err)
	}
	log.Print("in create user method 2")
	//release.ID=4
	//release.Name="may_19_release"
	//release.Type="new"
	//release.Owner="roopa@gmail.com"
	//release.TargetDate=time.Now()
	log.Print("release TargetDate", release.TargetDate)

	createdRelease := db.Debug().Create(&release)
	var errMessage = createdRelease.Error
	log.Print("error release", errMessage)

	if createdRelease.Error != nil {
		log.Print(errMessage)

	}
	 db.Model(&models.ReleaseProject{}).Create([]map[string]interface{}{
		{"ReleaseId": release.ID, "ProjectId": 1},
		{"ReleaseId": release.ID, "ProjectId": 2},
	})

	return release.ID, nil
	//if releaseProjectData.Error != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": errMessage})
	//	return
	//}
	//c.JSON(http.StatusOK, gin.H{"status": "success", "release created successful": &release})
}

func GetAllReleases (c *gin.Context)  ([]*models.Release, error) {
	var release []models.Release
	records := db.Debug().Find(&release)
	if records.Error != nil {
		log.Fatalln(records.Error)
	}
	log.Printf("%d rows found.", records.RowsAffected)
	rows, err := records.Rows()
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()

	releaseArr := []*models.Release{}

	for rows.Next() {
		release:= &models.Release{}
		err := db.Debug().ScanRows(rows, &release)
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("%+v\n", release)
		releaseArr=append(releaseArr, release)
	}
	return releaseArr, nil
}

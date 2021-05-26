package repositories

import (
	"github.com/gin-gonic/gin"
	"github.com/release-trackers/gin/controllers"
	"github.com/release-trackers/gin/models"
	"log"
)



func CreateRelease(c *gin.Context, release models.Release, app controllers.App) (uint, error){
	err := c.Bind(&release)
	if err != nil {
		log.Print(err)
	}
	createdRelease := app.Db.Debug().Create(&release)
	var errMessage = createdRelease.Error
	log.Print("error release", errMessage)

	if createdRelease.Error != nil {
		log.Print(errMessage)

	}
	 app.Db.Model(&models.ReleaseProject{}).Create([]map[string]interface{}{
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

func  GetAllReleases (c *gin.Context, app controllers.App)  ([]*models.Release, error) {
	var release []models.Release
	records := app.Db.Debug().Find(&release)
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
		err := app.Db.Debug().ScanRows(rows, &release)
		if err != nil {
			log.Fatalln(err)
		}

		//log.Printf("%+v\n", release)
		releaseArr=append(releaseArr, release)
	}
	return releaseArr, nil
}

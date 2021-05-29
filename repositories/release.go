package repositories

import (
	"github.com/gin-gonic/gin"
	"github.com/release-trackers/gin/database"
	"github.com/release-trackers/gin/models"
	"log"
)

var (
	db = database.InitConnection()
)

func CreateRelease(c *gin.Context, release models.Release, projectIds []int) (uint, error){
	err := c.Bind(&release)
	if err != nil {
		log.Print(err)
	}
	createdRelease := db.Debug().Create(&release)
	var errMessage = createdRelease.Error
	log.Print("error release", errMessage)

	if createdRelease.Error != nil {
		log.Print(errMessage)

	}
	// db.Model(&models.ReleaseProject{}).Create([]map[string]interface{}{
	//	{"ReleaseId": release.ID, "ProjectId": projectIds},
	//})

	for _, projectId := range projectIds {
		db.Create(&models.ReleaseProject{
			ReleaseId: release.ID, ProjectId: uint(projectId),
		})
	}
	return release.ID, errMessage
}

func  GetAllReleases (c *gin.Context)  ([]*models.Release, error) {
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
	//projectArr := []*models.Project{}

	for rows.Next() {
		release:= &models.Release{}
		err := db.Debug().ScanRows(rows, &release)
		if err != nil {
			log.Fatalln(err)
		}
		//projectArr = getReleaseProjects(release, err)
		releaseArr = append(releaseArr, release)
	}
	return releaseArr, err
}

func getReleaseProjects(release *models.Release, err error)  ([]*models.Project) {
	projects := []models.ReleaseProject{}
	log.Printf("release Id : %+v", release.ID)
	projectRecords := db.Debug().Where("release_id = ?", release.ID).Find(&projects)
	projrows, err := projectRecords.Rows()
	log.Printf("project %+v\n", projrows)
	projectArr := []*models.Project{}

	for projrows.Next() {
		project:= &models.Project{}
		err := db.Debug().ScanRows(projrows, &project)
		if err != nil {
			log.Fatalln(err)
		}
		projectArr = append(projectArr, project)
	}

	return projectArr
}

func GetProjects(c *gin.Context) ([]*models.Project, error) {
	var project []models.Project
	records := db.Debug().Find(&project)
	if records.Error != nil {
		log.Fatalln(records.Error)
	}
	//log.Printf("%d project rows found.", records.RowsAffected)
	rows, err := records.Rows()
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()

	projectArr := []*models.Project{}
	for rows.Next() {
		project:= &models.Project{}
		err := db.Debug().ScanRows(rows, &project)
		if err != nil {
			log.Fatalln(err)
		}
		//log.Printf("%+v\n", project)
		projectArr=append(projectArr, project)
	}
	return projectArr, err
}

package repositories

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/release-trackers/gin/models"
)

func (app *App) GetAllReviewers(dt models.DataTableValues) models.ReviewerResult {
	table := "reviewers"
	db := app.Db
	var total, filtered int64
	var reviewers []models.Reviewers
	query := db.Table(table)
	query = query.Offset(dt.Offset)
	query = query.Limit(dt.Limit)
	query = query.Scopes(dt.Search)
	query = query.Order("id " + dt.Order)

	if err := query.Find(&reviewers).Error; err != nil {
		return models.ReviewerResult{
			Total:    0,
			Filtered: 0,
			Data:     reviewers,
		}
	}

	// Filtered data count
	query.Table(table).Count(&filtered)

	// Total data count
	db.Table(table).Count(&total)

	return models.ReviewerResult{
		Total:    total,
		Filtered: filtered,
		Data:     reviewers,
	}
}

func (app *App) GetReviewer(c *gin.Context) (models.Reviewers, error) {
	id := c.Param("id")
	reviewer := models.Reviewers{}
	result := app.Db.First(&reviewer, id)
	return reviewer, result.Error
}

func (app *App) UpdateReviewer(c *gin.Context, reviewerData models.Reviewers) (uint, error) {
	reviewer, err := app.GetReviewer(c)
	if err != nil {
		return 0, err
	}
	updatedProject := app.Db.Model(&reviewer).Updates(&reviewerData)
	var errMessage = updatedProject.Error
	if updatedProject.Error != nil {
		log.Print(errMessage)

	}
	return reviewer.ID, errMessage
}

func (app *App) CreateReviewer(c *gin.Context, reviewer models.Reviewers) (uint, error) {
	err := c.Bind(&reviewer)
	if err != nil {
		log.Print(err)
	}
	createdProject := app.Db.Create(&reviewer)
	var errMessage = createdProject.Error

	if createdProject.Error != nil {
		log.Print(errMessage)

	}
	return reviewer.ID, errMessage
}

func (app *App) DeleteReviewer(c *gin.Context) (uint, error) {
	reviewer, err := app.GetReviewer(c)
	if err != nil {
		return 0, err
	}
	app.Db.Unscoped().Delete(&reviewer)
	return 1, nil
}

func (app *App) GetAllReviewersList(c *gin.Context) ([]*models.Reviewers, error) {
	db := app.Db
	var reviewersList []models.Reviewers
	records := db.Find(&reviewersList)
	if records.Error != nil {
		log.Fatalln(records.Error)
	}
	//log.Printf("%d project rows found.", records.RowsAffected)
	rows, err := records.Rows()
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()

	reviewerArr := []*models.Reviewers{}
	for rows.Next() {
		reviewer := &models.Reviewers{}
		err := db.ScanRows(rows, &reviewer)
		if err != nil {
			log.Fatalln(err)
		}
		reviewerArr = append(reviewerArr, reviewer)
	}
	return reviewerArr, err
}

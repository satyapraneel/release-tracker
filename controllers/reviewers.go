package controllers

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/release-trackers/gin/models"
	"github.com/release-trackers/gin/repositories"
	"github.com/release-trackers/gin/validation"
)

func (app *App) GetReviewers(c *gin.Context) {
	c.HTML(http.StatusOK, "reviewers/view", gin.H{})
}

func (app *App) GetListOfReviewers(c *gin.Context) {
	var columnOrder string
	columnOrder = "desc"
	c.Request.ParseForm()
	order := c.PostFormMap("order")
	for _, value := range order { // Order not specified
		if value == "asc" || value == "desc" {
			columnOrder = value
		}
	}
	dtValues := models.DataTableValues{
		Offset: QueryOffset(c),
		Limit:  QueryLimit(c),
		Search: SearchScope(c),
		Order:  columnOrder,
	}
	reviewerRepsitoryHandler := repositories.NewRepositoryHandler(app.Application)
	reviewers := reviewerRepsitoryHandler.GetAllReviewers(dtValues)
	c.JSON(http.StatusOK, reviewers)
}

func (app *App) CreateReviewersForm(c *gin.Context) {
	c.HTML(http.StatusOK, "reviewers/create", gin.H{
		"title": "Create Reviewer",
	})
}

func (app *App) CreateReviewer(c *gin.Context) {
	validation.ValidateReviewerCreate(c)
	err := c.Request.ParseForm()
	if err != nil {
		http.Error(c.Writer, "Bad Request", http.StatusBadRequest)
		return
	}

	reviewer := models.Reviewers{
		Name:     c.Request.PostForm.Get("name"),
		Email:    c.Request.PostForm.Get("email"),
		UserName: c.Request.PostForm.Get("user_name"),
	}
	repsitoryHandler := repositories.NewRepositoryHandler(app.Application)
	reviewerCreated, err := repsitoryHandler.CreateReviewer(c, reviewer)
	session := sessions.Default(c)
	if err != nil {
		log.Print(err)
		session.AddFlash(err, "error")
	}
	if reviewerCreated != 0 {
		session.AddFlash("Project created successfully", "success")
	}
	session.Save()
	c.Redirect(http.StatusFound, "/reviewers")
}

func (app *App) ViewReviewerForm(c *gin.Context) {
	releaseRepsitoryHandler := repositories.NewReleaseHandler(app.Application)
	reviewer, err := releaseRepsitoryHandler.GetReviewer(c)
	if err != nil {
		c.Redirect(http.StatusFound, "reviewers")
	}
	c.HTML(http.StatusOK, "reviewers/edit", gin.H{
		"values": reviewer,
	})
}

func (app *App) UpdateReviewer(c *gin.Context) {
	validation.ValidateReviewerCreate(c)
	err := c.Request.ParseForm()
	if err != nil {
		http.Error(c.Writer, "Bad Request", http.StatusBadRequest)
		return
	}

	reviewer := models.Reviewers{
		Name:     c.Request.PostForm.Get("name"),
		Email:    c.Request.PostForm.Get("email"),
		UserName: c.Request.PostForm.Get("user_name"),
	}
	repsitoryHandler := repositories.NewRepositoryHandler(app.Application)
	updatesReviewer, err := repsitoryHandler.UpdateReviewer(c, reviewer)
	session := sessions.Default(c)
	if err != nil {
		log.Print(err)
		session.AddFlash(err, "error")
	}
	if updatesReviewer != 0 {
		session.AddFlash("Reviewer updated successfully", "success")
	}
	session.Save()
	c.Redirect(http.StatusFound, "/reviewers")
}

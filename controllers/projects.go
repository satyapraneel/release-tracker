package controllers

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/release-trackers/gin/models"
	"github.com/release-trackers/gin/repositories"
	"github.com/release-trackers/gin/validation"
)

func (app *App) GetProjects(c *gin.Context) {
	c.HTML(http.StatusOK, "projects/view", gin.H{})
}

func (app *App) GetListOfProjects(c *gin.Context) {
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
	projectRepsitoryHandler := repositories.NewRepositoryHandler(app.Application)
	projects := projectRepsitoryHandler.GetAllProjects(dtValues)
	c.JSON(http.StatusOK, projects)
}

func (app *App) CreateProjectForm(c *gin.Context) {
	releaseRepsitoryHandler := repositories.NewReleaseHandler(app.Application)
	reviewers, _ := releaseRepsitoryHandler.GetAllReviewersList(c)
	c.HTML(http.StatusOK, "projects/create", gin.H{
		"title":     "Create Project",
		"reviewers": reviewers,
	})
}

func (app *App) CreateProject(c *gin.Context) {
	validation.ValidateProjectCreate(c)
	err := c.Request.ParseForm()
	if err != nil {
		http.Error(c.Writer, "Bad Request", http.StatusBadRequest)
		return
	}

	project := models.Project{
		Name:                 c.Request.PostForm.Get("name"),
		RepoName:             c.Request.PostForm.Get("repo_name"),
		ReviewerList:         c.Request.PostForm.Get("reviewer_list"),
		BetaReleaseDate:      c.Request.PostForm.Get("beta_release_date"),
		RegressionSignorDate: c.Request.PostForm.Get("regression_signor_date"),
		CodeFreezeDate:       c.Request.PostForm.Get("code_freeze_date"),
		DevCompletionDate:    c.Request.PostForm.Get("dev_completion_date"),
		RelatedCodes:         c.Request.PostForm.Get("related_codes"),
		Status:               c.Request.PostForm.Get("status"),
	}
	reviewersList := c.PostFormArray("reviewers")
	if len(reviewersList) > 0 {
		project.ReviewerList = strings.Join(reviewersList, ",")
	}
	repsitoryHandler := repositories.NewRepositoryHandler(app.Application)
	createProject, err := repsitoryHandler.CreateProject(c, project)
	session := sessions.Default(c)
	if err != nil {
		log.Print(err)
		session.AddFlash(err, "error")
	}
	if createProject != 0 {
		session.AddFlash("Project created successfully", "success")
	}
	session.Save()
	c.Redirect(http.StatusFound, "/projects")
}

func (app *App) ViewProjectForm(c *gin.Context) {
	releaseRepsitoryHandler := repositories.NewReleaseHandler(app.Application)
	reviewers, _ := releaseRepsitoryHandler.GetAllReviewersList(c)
	projects, err := releaseRepsitoryHandler.GetProject(c)
	if err != nil {
		c.Redirect(http.StatusFound, "projects")
	}
	c.HTML(http.StatusOK, "projects/edit", gin.H{
		"values":             projects,
		"reviewers":          reviewers,
		"selected_reviewers": projects.ReviewerList,
	})
}

func (app *App) UpdateProject(c *gin.Context) {
	validation.ValidateProjectCreate(c)
	err := c.Request.ParseForm()
	if err != nil {
		http.Error(c.Writer, "Bad Request", http.StatusBadRequest)
		return
	}

	project := models.Project{
		Name:                 c.Request.PostForm.Get("name"),
		RepoName:             c.Request.PostForm.Get("repo_name"),
		ReviewerList:         c.Request.PostForm.Get("reviewer_list"),
		BetaReleaseDate:      c.Request.PostForm.Get("beta_release_date"),
		RegressionSignorDate: c.Request.PostForm.Get("regression_signor_date"),
		CodeFreezeDate:       c.Request.PostForm.Get("code_freeze_date"),
		DevCompletionDate:    c.Request.PostForm.Get("dev_completion_date"),
		RelatedCodes:         c.Request.PostForm.Get("related_codes"),
		Status:               c.Request.PostForm.Get("status"),
	}
	reviewersList := c.PostFormArray("reviewers")
	if len(reviewersList) > 0 {
		project.ReviewerList = strings.Join(reviewersList, ",")
	} else {
		project.ReviewerList = ""
	}
	repsitoryHandler := repositories.NewRepositoryHandler(app.Application)
	updateProject, err := repsitoryHandler.UpdateProject(c, project)
	session := sessions.Default(c)
	if err != nil {
		log.Print(err)
		session.AddFlash(err, "error")
	}
	if updateProject != 0 {
		session.AddFlash("Project updated successfully", "success")
	}
	session.Save()
	c.Redirect(http.StatusFound, "/projects")
}

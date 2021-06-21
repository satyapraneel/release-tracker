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
	c.HTML(http.StatusOK, "projects/create", gin.H{
		"title": "Create Project",
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
		BitbucketUrl:         c.Request.PostForm.Get("bitbucket_url"),
		ReviewerList:         c.Request.PostForm.Get("reviewer_list"),
		BetaReleaseDate:      c.Request.PostForm.Get("beta_release_date"),
		RegressionSignorDate: c.Request.PostForm.Get("regression_signor_date"),
		CodeFreezeDate:       c.Request.PostForm.Get("code_freeze_date"),
		DevCompletionDate:    c.Request.PostForm.Get("dev_completion_date"),
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

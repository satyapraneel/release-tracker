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

func (app *App) GetDLs(c *gin.Context) {
	c.HTML(http.StatusOK, "dls/view", gin.H{})
}

func (app *App) GetListOfDLs(c *gin.Context) {
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
	repsitoryHandler := repositories.NewRepositoryHandler(app.Application)
	values := repsitoryHandler.GetAllDls(dtValues)
	c.JSON(http.StatusOK, values)
}

func (app *App) CreateDLsForm(c *gin.Context) {
	dlTypes := map[string]string{models.QA: "QA", models.DEV: "Developers", models.PM: "PMs", models.DEVOPS: "DevOps"}
	repsitoryHandler := repositories.NewRepositoryHandler(app.Application)
	projects, err := repsitoryHandler.GetAllProjectsList(c)
	if err != nil {
		emptyProject := []*models.Project{}
		projects = emptyProject
	}
	c.HTML(http.StatusOK, "dls/create", gin.H{
		"title":    "Create DLs",
		"DlTypes":  dlTypes,
		"Projects": projects,
	})
}

func (app *App) CreateDL(c *gin.Context) {
	validation.ValidateDLs(c)
	err := c.Request.ParseForm()
	if err != nil {
		http.Error(c.Writer, "Bad Request", http.StatusBadRequest)
		return
	}

	dl := models.DLS{
		Email:  c.Request.PostForm.Get("email"),
		DlType: c.Request.PostForm.Get("dl_type"),
	}
	dlProjectId := c.Request.PostForm.Get("project_id")
	repsitoryHandler := repositories.NewRepositoryHandler(app.Application)
	created, err := repsitoryHandler.CreateDL(c, dl, dlProjectId)
	session := sessions.Default(c)
	if err != nil {
		log.Print(err)
		session.AddFlash(err, "error")
	}
	if created != 0 {
		session.AddFlash("DL created successfully", "success")
	}
	session.Save()
	c.Redirect(http.StatusFound, "/dls")
}

func (app *App) ViewDLForm(c *gin.Context) {
	dlTypes := map[string]string{models.QA: "QA", models.DEV: "Developers", models.PM: "PMs", models.DEVOPS: "DevOps"}
	repsitoryHandler := repositories.NewReleaseHandler(app.Application)
	dl, err := repsitoryHandler.GetDL(c.Param("id"))

	if err != nil {
		c.Redirect(http.StatusFound, "dls")
	}

	selectedProject, project_error := repsitoryHandler.GetDLProject(c.Param("id"))
	if project_error != nil {
		c.Redirect(http.StatusFound, "dls")
	}
	projects, _ := repsitoryHandler.GetAllProjectsList(c)
	c.HTML(http.StatusOK, "dls/edit", gin.H{
		"values":          dl,
		"DlTypes":         dlTypes,
		"SelectedDL":      dl.DlType,
		"Projects":        projects,
		"SelectedProject": selectedProject.ProjectId,
	})
}

func (app *App) UpdateDL(c *gin.Context) {
	validation.ValidateDLs(c)
	err := c.Request.ParseForm()
	if err != nil {
		http.Error(c.Writer, "Bad Request", http.StatusBadRequest)
		return
	}

	dl := models.DLS{
		Email:  c.Request.PostForm.Get("email"),
		DlType: c.Request.PostForm.Get("dl_type"),
	}
	dlProjectId := c.Request.PostForm.Get("project_id")
	repsitoryHandler := repositories.NewRepositoryHandler(app.Application)
	updated, err := repsitoryHandler.UpdateDL(c, dl, dlProjectId)
	session := sessions.Default(c)
	if err != nil {
		log.Print(err)
		session.AddFlash(err, "error")
	}
	if updated != 0 {
		session.AddFlash("DL updated successfully", "success")
	}
	session.Save()
	c.Redirect(http.StatusFound, "/dls")
}

func (app *App) DeleteDL(c *gin.Context) {
	repsitoryHandler := repositories.NewRepositoryHandler(app.Application)
	deleteRecord, err := repsitoryHandler.DeleteDL(c.Param("id"))
	session := sessions.Default(c)
	if err != nil {
		log.Print(err)
		session.AddFlash(err, "error")
	}
	if deleteRecord != 0 {
		session.AddFlash("DL deleted successfully", "success")
	}
	session.Save()
	c.Redirect(http.StatusFound, "/dls")
}

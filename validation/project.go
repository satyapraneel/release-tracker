package validation

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ProjectFields struct {
	Name                 string `form:"name" json:"name" binding:"required"`
	Reponame             string `form:"repo_name" json:"repo_name" binding:"required"`
	BetaReleaseDate      string `form:"beta_release_date" json:"beta_release_date" binding:"required,number"`
	RegressionSignorDate string `form:"regression_signor_date" json:"regression_signor_date" binding:"required,number"`
	CodeFreezeDate       string `form:"code_freeze_date" json:"code_freeze_date" binding:"required,number"`
	DevCompletionDate    string `form:"dev_completion_date" json:"dev_completion_date" binding:"required,number"`
	Status               string `form:"status" json:"status" binding:"required,number"`
}

func ValidateProjectCreate(c *gin.Context) {
	var form ProjectFields
	// This will infer what binder to use depending on the content-type header.
	if err := c.ShouldBind(&form); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			// c.JSON(http.StatusBadRequest, gin.H{"errors": Descriptive(verr)})

			// print(Descriptive(verr))
			// print(err.Error())

			c.HTML(http.StatusUnauthorized, "projects/create", gin.H{"errors": DescriptiveError(verr), "values": form})
			c.Abort()
		}
	}
}

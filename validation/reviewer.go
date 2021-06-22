package validation

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ReviewerFields struct {
	Name     string `form:"name" json:"name" binding:"required"`
	Email    string `form:"email" json:"email" binding:"required"`
	UserName string `form:"user_name" json:"user_name" binding:"required"`
}

func ValidateReviewerCreate(c *gin.Context) {
	var form ReviewerFields
	// This will infer what binder to use depending on the content-type header.
	if err := c.ShouldBind(&form); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			// c.JSON(http.StatusBadRequest, gin.H{"errors": Descriptive(verr)})

			// print(Descriptive(verr))
			// print(err.Error())

			c.HTML(http.StatusUnauthorized, "reviewers/create", gin.H{"errors": DescriptiveError(verr), "values": form})
			c.Abort()
		}
	}
}

package validation

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type DLsFields struct {
	Email     string `form:"email" json:"email" binding:"required"`
	DlType    string `form:"dl_type" json:"dl_type" binding:"required"`
	ProjectID string `form:"project_id" json:"project_id" binding:"required"`
}

func ValidateDLs(c *gin.Context) {
	var form DLsFields
	// This will infer what binder to use depending on the content-type header.
	if err := c.ShouldBind(&form); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			c.HTML(http.StatusUnauthorized, "dls/create", gin.H{"errors": DescriptiveError(verr), "values": form})
			c.Abort()
		}
	}
}

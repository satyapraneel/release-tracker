package validation

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type LoginFields struct {
	Email    string `form:"email" json:"email" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required"`
}

func ValidateLoginForm(c *gin.Context) {
	var form LoginFields
	// This will infer what binder to use depending on the content-type header.
	if err := c.ShouldBind(&form); err != nil {
		log.Print(form)
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			// c.JSON(http.StatusBadRequest, gin.H{"errors": Descriptive(verr)})

			// print(Descriptive(verr))
			// print(err.Error())

			c.HTML(http.StatusUnauthorized, "auth/login", gin.H{"errors": DescriptiveError(verr), "values": form})
			c.Abort()
		}
	}
}

package helpers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"h-two/internal/errors"
	"net/http"
	"strings"
)

func ParseRequestBody(c *gin.Context, req interface{}) any {
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		if validationErrs, ok := bindErr.(validator.ValidationErrors); ok {
			// Handle validation errors
			var res []errors.FieldError
			for _, e := range validationErrs {
				// Extract the field name and the error message
				fieldName := strings.Split(e.Namespace(), ".")[1]
				errorMessage := e.ActualTag()
				// Translate each error one at a time
				res = append(res, errors.FieldError{Field: fieldName, Message: errorMessage})
			}
			c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": res})

		} else {
			// Handle other errors (like invalid JSON)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		}
		return bindErr
	}
	return nil
}

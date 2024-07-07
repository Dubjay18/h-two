package helpers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"h-two/internal/errors"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func ParseRequestBody(c *gin.Context, req interface{}) error {
	if c.Request.ContentLength == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "request body must not be empty"})
		return nil
	}

	// Read the request body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't read request body"})
		return nil
	}

	// Unmarshal the JSON body
	err = json.Unmarshal(body, req)
	if err != nil {
		// Check if the error is a JSON syntax error
		if syntaxErr, ok := err.(*json.SyntaxError); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON syntax at " + strconv.FormatInt(syntaxErr.Offset, 10)})
			return syntaxErr
		}
		if typeErr, ok := err.(*json.UnmarshalTypeError); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid value for field " + typeErr.Field + " at offset " + strconv.FormatInt(typeErr.Offset, 10)})
			return typeErr
		}

		// Handle other errors
		errs := err.(validator.ValidationErrors)
		var res []errors.FieldError
		for _, e := range errs {
			// Extract the field name and the error message
			fieldName := strings.Split(e.Namespace(), ".")[1]
			errorMessage := e.ActualTag()
			// Translate each error one at a time
			res = append(res, errors.FieldError{Field: fieldName, Message: errorMessage})
		}
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": res})
		return err
	}

	return nil
}

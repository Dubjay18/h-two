package server

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"h-two/internal/dto"
	"h-two/internal/errors"
	"net/http"
	"strings"
)

func (s *Server) RegisterHandler(c *gin.Context) {
	var req *dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errs := err.(validator.ValidationErrors)
		var res []errors.FieldError
		for _, e := range errs {
			// Extract the field name and the error message
			fieldName := strings.Split(e.Namespace(), ".")[1]
			errorMessage := e.ActualTag()
			// Translate each error one at a time
			res = append(res, errors.FieldError{Field: fieldName, Message: errorMessage})
		}
		c.JSON(http.StatusBadRequest, gin.H{"errors": res})
		return
	}
	resp, err := s.AuthService.CreateUser(c, req)
	if err != nil {
		c.JSON(err.StatusCode, err)
		return

	}
	response := make(map[string]interface{})
	response["user"] = resp
	c.JSON(http.StatusCreated, dto.ApiSuccessResponse{
		Status:  "success",
		Message: "Registration successful",
		Data:    response,
	})

}

func (s *Server) LoginHandler(c *gin.Context) {
	var req *dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errs := err.(validator.ValidationErrors)
		var res []errors.FieldError
		for _, e := range errs {
			// Extract the field name and the error message
			fieldName := strings.Split(e.Namespace(), ".")[1]
			errorMessage := e.ActualTag()
			// Translate each error one at a time
			res = append(res, errors.FieldError{Field: fieldName, Message: errorMessage})
		}
		c.JSON(http.StatusBadRequest, gin.H{"errors": res})
		return
	}
	resp, err := s.AuthService.Login(c, req)
	if err != nil {
		c.JSON(err.StatusCode, err)
		return
	}

	c.JSON(http.StatusOK, dto.ApiSuccessResponse{
		Status:  "success",
		Message: "Login successful",
		Data:    resp,
	})
}

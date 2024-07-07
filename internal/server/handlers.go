package server

import (
	"github.com/gin-gonic/gin"
	"h-two/internal/dto"
	"h-two/internal/helpers"
	"log"
	"net/http"
)

func (s *Server) RegisterHandler(c *gin.Context) {
	var req *dto.CreateUserRequest
	perr := helpers.ParseRequestBody(c, &req)
	if perr != nil {
		return
	}
	resp, err := s.AuthService.CreateUserAndOrganization(c, req)
	if err != nil {
		c.JSON(err.StatusCode, err)
		return

	}

	c.JSON(http.StatusCreated, dto.ApiSuccessResponse{
		Status:  "success",
		Message: "Registration successful",
		Data:    resp,
	})

}

func (s *Server) LoginHandler(c *gin.Context) {
	var req *dto.LoginRequest
	perr := helpers.ParseRequestBody(c, &req)
	if perr != nil {
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

func (s *Server) GetUserDetailsHandler(c *gin.Context) {

	// Get the user ID from the context
	userID := c.Params.ByName("id")
	user, err := s.UserService.GetUserDetails(c, userID)
	if err != nil {
		c.JSON(err.StatusCode, err)
		return
	}
	c.JSON(http.StatusOK, dto.ApiSuccessResponse{
		Status:  "success",
		Message: "User details retrieved successfully",
		Data:    user,
	})
}
func (s *Server) GetOrganizationsHandler(c *gin.Context) {
	// Get the user ID from the context
	userID := c.GetString("userId")
	orgs, err := s.OrganizationService.GetUserOrganizations(userID)
	if err != nil {
		c.JSON(err.StatusCode, err)
		return
	}
	c.JSON(http.StatusOK, dto.ApiSuccessResponse{
		Status:  "success",
		Message: "Organizations retrieved successfully",
		Data: gin.H{
			"organisations": orgs,
		},
	})
}

func (s *Server) GetOrganizationHandler(c *gin.Context) {
	// Get the user ID from the context
	userID := c.GetString("userId")
	// Get the organization ID from the URL parameters
	orgId := c.Param("orgId")
	org, err := s.OrganizationService.GetOrganizationById(userID, orgId)
	if err != nil {
		c.JSON(err.StatusCode, err)
		return
	}
	c.JSON(http.StatusOK, dto.ApiSuccessResponse{
		Status:  "success",
		Message: "Organization retrieved successfully",
		Data:    org,
	})
}

func (s *Server) CreateOrganizationHandler(c *gin.Context) {
	userID := c.GetString("userId")
	var req dto.CreateOrganizationRequest
	perr := helpers.ParseRequestBody(c, &req)
	if perr != nil {
		log.Println(perr)
		return
	}

	org, err := s.OrganizationService.CreateOrganization(userID, &req)
	if err != nil {
		c.JSON(err.StatusCode, err)
		return
	}

	c.JSON(http.StatusCreated, dto.ApiSuccessResponse{
		Status:  "success",
		Message: "Organization created successfully",
		Data:    org,
	})

}

func (s *Server) AddUserToOrganizationHandler(c *gin.Context) {
	orgID := c.Param("orgId")
	var req dto.AddUserToOrganizationRequest
	perr := helpers.ParseRequestBody(c, &req)
	if perr != nil {
		log.Println(perr)
		return
	}
	err := s.OrganizationService.AddUserToOrganization(req.UserId, orgID)
	if err != nil {
		c.JSON(err.StatusCode, err)
		return
	}

	c.JSON(http.StatusOK, dto.ApiSuccessResponse{
		Status:  "success",
		Message: "User added to organization successfully",
	})
}

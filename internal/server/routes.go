package server

import (
	"h-two/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.GET("/", s.HelloWorldHandler)
	authGroup := r.Group("/auth")
	apiGroup := r.Group("/api")
	{
		authGroup.POST("/register", s.RegisterHandler)
		authGroup.POST("/login", s.LoginHandler)
		apiGroup.GET("/users/:id", middleware.AuthMiddleware, s.GetUserDetailsHandler)
		apiGroup.GET("/organisations", middleware.AuthMiddleware, s.GetOrganizationsHandler)
		apiGroup.GET("/organisations/:orgId", middleware.AuthMiddleware, s.GetOrganizationHandler)
		apiGroup.POST("/organisations", middleware.AuthMiddleware, s.CreateOrganizationHandler)
		apiGroup.POST("/organisations/:orgId/users", middleware.AuthMiddleware, s.AddUserToOrganizationHandler)
	}

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

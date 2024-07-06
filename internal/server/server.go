package server

import (
	"fmt"
	"h-two/internal/repository"
	"h-two/internal/services"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"h-two/internal/database"
)

type Server struct {
	port                int
	AuthService         services.AuthService
	OrganizationService services.OrganizationService
	Db                  *database.DbService
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	dbInstance := database.New()

	organizationRep := repository.NewOrganizationRepository(dbInstance.Db)
	organizationService := services.NewOrganizationService(organizationRep)
	userRepo := repository.NewUserRepository(dbInstance.Db)               // Pass the dbInstance to the UserRepository
	authService := services.NewAuthService(userRepo, organizationService) // Pass the UserRepository to the AuthService

	NewServer := &Server{
		port:                port,
		AuthService:         authService,
		OrganizationService: organizationService,
		Db:                  database.New(),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

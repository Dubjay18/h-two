package services

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"h-two/internal/dto"
	"h-two/internal/errors"
	"h-two/internal/models"
	"h-two/internal/repository"
	"net/http"
)

type AuthService interface {
	CreateUser(c *gin.Context, user *dto.CreateUserRequest) (*dto.CreateUserResponse, *errors.ApiError)
}

type DefaultAuthService struct {
	repo *repository.DefaultUserRepository
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func verifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (a *DefaultAuthService) CreateUser(c *gin.Context, user *dto.CreateUserRequest) (*dto.CreateUserResponse, *errors.ApiError) {
	// Hash the user's password
	hash, err := hashPassword(user.Password)
	if err != nil {
		return nil, &errors.ApiError{
			Status:     errors.InternalServerError,
			Message:    "Registration unsuccessful",
			StatusCode: http.StatusInternalServerError,
		}
	}
	user.Password = hash

	// Save the user to the database
	userResponse, dbErr := a.repo.CreateUser(&models.User{FirstName: user.FirstName,
		Email:    user.Email,
		Password: user.Password,
		LastName: user.LastName,
	})
	if dbErr != nil {
		return nil, &errors.ApiError{
			Status:     errors.InternalServerError,
			Message:    "Registration unsuccessful",
			StatusCode: http.StatusInternalServerError,
		}
	}
	return userResponse, nil
}

func NewAuthService(repo *repository.DefaultUserRepository) AuthService {
	return &DefaultAuthService{repo: repo}
}

package services

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"h-two/internal/dto"
	"h-two/internal/errors"
	"h-two/internal/models"
	"h-two/internal/repository"
	"net/http"
	"os"
	"time"
)

type AuthService interface {
	CreateUser(c *gin.Context, user *dto.CreateUserRequest) (*dto.CreateUserResponse, *errors.ApiError)
	Login(c *gin.Context, user *dto.LoginRequest) (*dto.LoginResponse, *errors.ApiError)
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

func GenerateJWT(userId string) (string, error) {
	secretKey := os.Getenv("JWT_SECRET")
	expirationTime := time.Now().Add(1 * time.Hour).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"exp":    expirationTime,
	})
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
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

func (a *DefaultAuthService) Login(c *gin.Context, user *dto.LoginRequest) (*dto.LoginResponse, *errors.ApiError) {

	// Get the user from the database
	u, err := a.repo.GetUserByEmail(user.Email)
	if err != nil {
		return nil, &errors.ApiError{
			Status:     errors.ValidationError,
			Message:    "Authentication Failed",
			StatusCode: http.StatusInternalServerError,
		}
	}
	// Verify the user's password
	if !verifyPassword(user.Password, u.Password) {
		return nil, &errors.ApiError{
			Status:     errors.ValidationError,
			Message:    "Authentication Failed",
			StatusCode: http.StatusUnauthorized,
		}
	}
	// Generate a JWT token
	token, err := GenerateJWT(u.UserId)
	if err != nil {
		return nil, &errors.ApiError{
			Status:     errors.ValidationError,
			Message:    "Authentication Failed",
			StatusCode: http.StatusUnauthorized,
		}
	}
	return &dto.LoginResponse{
		AccessToken: token,
		User: struct {
			UserId    string `json:"userId"`
			FirstName string `json:"firstName"`
			LastName  string `json:"lastName"`
			Email     string `json:"email"`
			Phone     string `json:"phone"`
		}{
			UserId:    u.UserId,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Email:     u.Email,
			Phone:     u.Phone,
		},
	}, nil
}
func NewAuthService(repo *repository.DefaultUserRepository) AuthService {
	return &DefaultAuthService{repo: repo}
}

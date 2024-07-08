package services

import (
	"github.com/gin-gonic/gin"
	"h-two/internal/dto"
	"h-two/internal/errors"
	"h-two/internal/repository"
	"net/http"
)

type UserService interface {
	GetUserDetails(c *gin.Context, userId string) (*dto.UserResponse, *errors.ApiError)
}

type DefaultUserService struct {
	repo repository.UserRepository
}

func (s *DefaultUserService) GetUserDetails(c *gin.Context, userId string) (*dto.UserResponse, *errors.ApiError) {
	if c.GetString("userId") != userId {
		return nil, &errors.ApiError{
			Message:    "Unauthorized",
			StatusCode: http.StatusUnauthorized,
			Status:     errors.UnAuthorized,
		}
	}
	user, err := s.repo.GetUserById(userId)
	if err != nil {
		return nil, &errors.ApiError{
			Message:    "User not found",
			StatusCode: http.StatusNotFound,
			Status:     errors.UserNotFound,
		}
	}
	return &dto.UserResponse{
		UserId:    user.UserId,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Phone:     user.Phone,
	}, nil
}

func NewUserService(repo repository.UserRepository) *DefaultUserService {
	return &DefaultUserService{repo: repo}
}

package services

import (
	"github.com/gin-gonic/gin"
	"h-two/internal/dto"
	"h-two/internal/errors"
	"h-two/internal/repository"
	"log"
	"net/http"
)

type UserService interface {
	GetUserDetails(c *gin.Context, userId string) (*dto.UserResponse, *errors.ApiError)
}

type DefaultUserService struct {
	repo repository.UserRepository
}

//func (s *DefaultUserService) GetUserDetails(c *gin.Context, userId string) (*dto.UserResponse, *errors.ApiError) {
//	if c.GetString("userId") != userId {
//		return nil, &errors.ApiError{
//			Message:    "Unauthorized",
//			StatusCode: http.StatusUnauthorized,
//			Status:     errors.UnAuthorized,
//		}
//	}
//	user, err := s.repo.GetUserById(userId)
//	if err != nil {
//		return nil, &errors.ApiError{
//			Message:    "User not found",
//			StatusCode: http.StatusNotFound,
//			Status:     errors.UserNotFound,
//		}
//	}
//	return &dto.UserResponse{
//		UserId:    user.UserId,
//		FirstName: user.FirstName,
//		LastName:  user.LastName,
//		Email:     user.Email,
//		Phone:     user.Phone,
//	}, nil
//}

func (s *DefaultUserService) GetUserDetails(c *gin.Context, userId string) (*dto.UserResponse, *errors.ApiError) {
	requestingUserId := c.GetString("userId")
	log.Println("Requesting user ID: ", requestingUserId)
	if requestingUserId == userId {
		// The user is requesting their own details
		user, err := s.repo.GetUserById(userId)
		if err != nil {
			return nil, &errors.ApiError{
				Message:    "Unauthorized",
				StatusCode: http.StatusUnauthorized,
				Status:     errors.UnAuthorized,
			}
		}
		return &dto.UserResponse{
			UserId:    user.UserId,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Phone:     user.Phone,
		}, nil
	} else {
		// The user is requesting details of another user\
		// Check if the requesting user is in the same organization as the user

		if same, _ := s.repo.AreUsersInSameOrganization(requestingUserId, userId); same {
			log.Println(same)
			// get requesting user details
			requestingUser, err := s.repo.GetUserById(userId)
			if err != nil {
				return nil, &errors.ApiError{
					Message:    "Unauthorized",
					StatusCode: http.StatusUnauthorized,
					Status:     errors.UnAuthorized,
				}
			}
			return &dto.UserResponse{
				UserId:    requestingUser.UserId,
				FirstName: requestingUser.FirstName,
				LastName:  requestingUser.LastName,
				Email:     requestingUser.Email,
				Phone:     requestingUser.Phone,
			}, nil

		}

		// The user is not in the same organization as the user
		return nil, &errors.ApiError{
			Message:    "Unauthorized",
			StatusCode: http.StatusUnauthorized,
			Status:     errors.UnAuthorized,
		}
	}
}

func NewUserService(repo repository.UserRepository) *DefaultUserService {
	return &DefaultUserService{repo: repo}
}

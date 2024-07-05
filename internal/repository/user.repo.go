package repository

import (
	"gorm.io/gorm"
	"h-two/internal/dto"
	"h-two/internal/models"
)

type UserRepository interface {
	CreateUser(user *models.User) (*dto.CreateUserResponse, error)
}

type DefaultUserRepository struct {
	db *gorm.DB
}

func (r *DefaultUserRepository) CreateUser(user *models.User) (*dto.CreateUserResponse, error) {
	err := r.db.Create(&user).Error
	if err != nil {
		return &dto.CreateUserResponse{}, err
	}
	return &dto.CreateUserResponse{
		UserId:    user.UserId,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Phone:     user.Phone,
	}, nil
}

func NewUserRepository(db *gorm.DB) *DefaultUserRepository {
	return &DefaultUserRepository{db: db}
}

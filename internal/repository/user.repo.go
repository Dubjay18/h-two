package repository

import (
	"gorm.io/gorm"
	"h-two/internal/dto"
	"h-two/internal/models"
)

type UserRepository interface {
	CreateUser(user *models.User) (*dto.CreateUserResponse, error)
	GetUserByEmail(email string) (*models.User, error)
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

func (r *DefaultUserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil

}

func NewUserRepository(db *gorm.DB) *DefaultUserRepository {
	return &DefaultUserRepository{db: db}
}

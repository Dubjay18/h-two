package repository

import (
	"gorm.io/gorm"
	"h-two/internal/dto"
	"h-two/internal/models"
)

type UserRepository interface {
	CreateUser(user *models.User) (*dto.UserResponse, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserById(userId string) (*models.User, error)
	Begin() *gorm.DB
}

type DefaultUserRepository struct {
	db *gorm.DB
}

func (r *DefaultUserRepository) CreateUser(user *models.User) (*dto.UserResponse, error) {
	if u := r.db.Where("email = ?", user.Email).First(&models.User{}); u.RowsAffected > 0 {
		return &dto.UserResponse{}, gorm.ErrRecordNotFound
	}
	err := r.db.Create(&user).Error
	if err != nil {
		return &dto.UserResponse{}, err
	}
	return &dto.UserResponse{
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

func (r *DefaultUserRepository) GetUserById(userId string) (*models.User, error) {
	var user models.User
	err := r.db.Where("user_id = ?", userId).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil

}

func (r *DefaultUserRepository) Begin() *gorm.DB {
	return r.db.Begin()
}

func NewUserRepository(db *gorm.DB) *DefaultUserRepository {
	return &DefaultUserRepository{db: db}
}

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
	GetUserOrganization(id string) (*models.User, error)
	AreUsersInSameOrganization(userId1 string, userId2 string) (bool, error)
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

func (r *DefaultUserRepository) GetUserOrganization(id string) (*models.User, error) {
	var user models.User
	err := r.db.Where("user_id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *DefaultUserRepository) AreUsersInSameOrganization(userId1 string, userId2 string) (bool, error) {
	var userOrgs1 []models.UserOrganization
	var userOrgs2 []models.UserOrganization
	err := r.db.Where("user_id = ?", userId1).Find(&userOrgs1).Error
	if err != nil {
		return false, err
	}
	err = r.db.Where("user_id = ?", userId2).Find(&userOrgs2).Error
	if err != nil {
		return false, err
	}

	// Create a map for faster lookup
	orgs1 := make(map[string]bool)
	for _, org := range userOrgs1 {
		orgs1[org.OrgId] = true
	}

	// Check if any organization of user2 is also in user1's organizations
	for _, org := range userOrgs2 {
		if _, ok := orgs1[org.OrgId]; ok {
			return true, nil
		}
	}

	return false, nil
}
func (r *DefaultUserRepository) Begin() *gorm.DB {
	return r.db.Begin()
}

func NewUserRepository(db *gorm.DB) *DefaultUserRepository {
	return &DefaultUserRepository{db: db}
}

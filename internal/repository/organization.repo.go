package repository

import (
	"fmt"
	"gorm.io/gorm"
	"h-two/internal/models"
)

type OrganizationRepository interface {
	CreateOrganization(org *models.Organization) error
	GetOrganizationsByUser(userId string) ([]*models.Organization, error)
	GetOrganizationById(userId string, orgId string) (*models.Organization, error)
	AddUserToOrganization(orgId string, userId string) error
	IsUserInOrganization(userId string, orgId string) (bool, error)
	AreUsersInSameOrganization(userId1 string, userId2 string) (bool, error)
	Begin() *gorm.DB
}

type DefaultOrganizationRepository struct {
	db *gorm.DB
}

func (r *DefaultOrganizationRepository) IsUserInOrganization(userId string, orgId string) (bool, error) {
	var userOrg models.UserOrganization
	if err := r.db.Where("org_id = ? AND user_id = ?", orgId, userId).First(&userOrg).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil

}

func (r *DefaultOrganizationRepository) CreateOrganization(org *models.Organization) error {

	// Start a new transaction
	tx := r.db.Begin()
	// Check for errors starting the transaction
	if tx.Error != nil {
		return tx.Error
	}
	// Create the organization
	err := r.db.Create(&org).Error
	if err != nil {
		return err
	}
	// Add the user to the organization
	userOrg := &models.UserOrganization{
		OrgId:  org.OrgId,
		UserId: org.Owner,
	}
	err = tx.Create(&userOrg).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	// Commit the transaction
	tx.Commit()
	return nil
}
func (r *DefaultOrganizationRepository) GetOrganizationsByUser(userId string) ([]*models.Organization, error) {
	var orgs []*models.Organization
	err := r.db.Table("organizations").
		Joins("JOIN user_organizations ON organizations.org_id = user_organizations.org_id").
		Where("user_organizations.user_id = ?", userId).
		Find(&orgs).Error
	if err != nil {
		return nil, err
	}
	return orgs, nil
}

func (r *DefaultOrganizationRepository) GetOrganizationById(userId string, orgId string) (*models.Organization, error) {

	var org models.Organization
	err := r.db.Table("organizations").
		Joins("JOIN user_organizations ON organizations.org_id = user_organizations.org_id").
		Where("user_organizations.user_id = ? AND organizations.org_id = ?", userId, orgId).
		First(&org).Error
	if err != nil {
		return nil, err
	}
	return &org, nil
}

func (r *DefaultOrganizationRepository) AddUserToOrganization(orgId string, userId string) error {
	// Check if the user already belongs to the organization
	var userOrg models.UserOrganization
	if err := r.db.Where("org_id = ? AND user_id = ?", orgId, userId).First(&userOrg).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			// An error occurred while trying to fetch the record
			return err
		}
	} else {
		// The user already belongs to the organization
		return fmt.Errorf("user %s already belongs to organization %s", userId, orgId)
	}

	// The user does not belong to the organization, so add them
	userOrg = models.UserOrganization{
		OrgId:  orgId,
		UserId: userId,
	}
	if err := r.db.Create(&userOrg).Error; err != nil {
		return err
	}

	return nil
}

func (r *DefaultOrganizationRepository) AreUsersInSameOrganization(userId1 string, userId2 string) (bool, error) {
	var userOrg1, userOrg2 models.UserOrganization
	if err := r.db.Where("user_id = ?", userId1).First(&userOrg1).Error; err != nil {
		return false, err
	}
	if err := r.db.Where("user_id = ?", userId2).First(&userOrg2).Error; err != nil {
		return false, err
	}
	return userOrg1.OrgId == userOrg2.OrgId, nil
}

func (r *DefaultOrganizationRepository) Begin() *gorm.DB {
	return r.db.Begin()
}

func NewOrganizationRepository(db *gorm.DB) *DefaultOrganizationRepository {
	return &DefaultOrganizationRepository{db: db}
}

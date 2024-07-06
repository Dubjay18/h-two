package repository

import (
	"gorm.io/gorm"
	"h-two/internal/models"
)

type OrganizationRepository interface {
	CreateOrganization(org *models.Organization) error
	Begin() *gorm.DB
}

type DefaultOrganizationRepository struct {
	db *gorm.DB
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
func (r *DefaultOrganizationRepository) Begin() *gorm.DB {
	return r.db.Begin()
}

func NewOrganizationRepository(db *gorm.DB) *DefaultOrganizationRepository {
	return &DefaultOrganizationRepository{db: db}
}
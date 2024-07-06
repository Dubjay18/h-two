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
	err := r.db.Create(&org).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *DefaultOrganizationRepository) Begin() *gorm.DB {
	return r.db.Begin()
}

func NewOrganizationRepository(db *gorm.DB) *DefaultOrganizationRepository {
	return &DefaultOrganizationRepository{db: db}
}

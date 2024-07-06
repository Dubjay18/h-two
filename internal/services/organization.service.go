package services

import (
	"fmt"
	"h-two/internal/errors"
	"h-two/internal/models"
	"h-two/internal/repository"
)

type OrganizationService interface {
	CreateOrganizationByFirstName(name string) *errors.ApiError
}

type DefaultOrganizationService struct {
	repo *repository.DefaultOrganizationRepository
}

func (s *DefaultOrganizationService) CreateOrganizationByFirstName(name string) *errors.ApiError {
	org := &models.Organization{
		Name: fmt.Sprintf("%s's Organization", name),
	}
	err := s.repo.CreateOrganization(org)
	if err != nil {
		return &errors.ApiError{
			Message:    "Failed to create organization",
			StatusCode: 500,
			Status:     errors.InternalServerError,
		}
	}
	return nil
}

func NewOrganizationService(repo *repository.DefaultOrganizationRepository) *DefaultOrganizationService {
	return &DefaultOrganizationService{repo: repo}
}

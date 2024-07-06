package services

import (
	"fmt"
	"h-two/internal/models"
	"h-two/internal/repository"
)

type OrganizationService interface {
	CreateOrganizationByFirstName(name string) error
}

type DefaultOrganizationService struct {
	repo *repository.DefaultOrganizationRepository
}

func (s *DefaultOrganizationService) CreateOrganizationByFirstName(name string) error {
	org := &models.Organization{
		Name: fmt.Sprintf("%s's Organization", name),
	}
	err := s.repo.CreateOrganization(org)
	if err != nil {
		return err
	}
	return nil
}

func NewOrganizationService(repo *repository.DefaultOrganizationRepository) *DefaultOrganizationService {
	return &DefaultOrganizationService{repo: repo}
}
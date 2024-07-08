package services

import (
	"fmt"
	"h-two/internal/dto"
	"h-two/internal/errors"
	"h-two/internal/models"
	"h-two/internal/repository"
	"net/http"
)

type OrganizationService interface {
	CreateOrganizationByFirstName(name string, userId string) *errors.ApiError
	GetUserOrganizations(userId string) ([]*dto.GetOrganizationResponse, *errors.ApiError)
	GetOrganizationById(userId string, orgId string) (*dto.GetOrganizationResponse, *errors.ApiError)
	CreateOrganization(userId string, req *dto.CreateOrganizationRequest) (*dto.GetOrganizationResponse, *errors.ApiError)
	AddUserToOrganization(orgId string, userId string) *errors.ApiError
	IsUserInOrganization(userId string, orgId string) (bool, *errors.ApiError)
}

type DefaultOrganizationService struct {
	repo repository.OrganizationRepository
}

func (s *DefaultOrganizationService) IsUserInOrganization(userId string, orgId string) (bool, *errors.ApiError) {
	inOrg, err := s.repo.IsUserInOrganization(userId, orgId)
	if err != nil {
		return false, &errors.ApiError{
			Message:    "Client Error",
			StatusCode: http.StatusBadRequest,
			Status:     errors.ValidationError,
		}
	}
	return inOrg, nil

}
func (s *DefaultOrganizationService) CreateOrganizationByFirstName(name string, userId string) *errors.ApiError {

	org := &models.Organization{
		Name:  fmt.Sprintf("%s's Organization", name),
		Owner: userId,
	}

	err := s.repo.CreateOrganization(org)
	if err != nil {
		return &errors.ApiError{
			Message:    "Client Error",
			StatusCode: http.StatusBadRequest,
			Status:     errors.ValidationError,
		}
	}
	return nil
}

func (s *DefaultOrganizationService) GetUserOrganizations(userId string) ([]*dto.GetOrganizationResponse, *errors.ApiError) {
	orgs, err := s.repo.GetOrganizationsByUser(userId)
	if err != nil {
		return nil, &errors.ApiError{
			Message:    "Failed to get organizations",
			StatusCode: 500,
			Status:     errors.InternalServerError,
		}
	}
	var response []*dto.GetOrganizationResponse
	for _, org := range orgs {
		response = append(response, &dto.GetOrganizationResponse{
			OrgId:       org.OrgId,
			Name:        org.Name,
			Description: org.Description,
		})
	}
	return response, nil
}
func (s *DefaultOrganizationService) GetOrganizationById(userId string, orgId string) (*dto.GetOrganizationResponse, *errors.ApiError) {
	org, err := s.repo.GetOrganizationById(userId, orgId)
	if err != nil {
		return nil, &errors.ApiError{
			Message:    "Organization not found",
			StatusCode: 404,
			Status:     "Not Found",
		}
	}
	return &dto.GetOrganizationResponse{
		OrgId:       org.OrgId,
		Name:        org.Name,
		Description: org.Description,
	}, nil
}

func (s *DefaultOrganizationService) CreateOrganization(userId string, req *dto.CreateOrganizationRequest) (*dto.GetOrganizationResponse, *errors.ApiError) {
	org := &models.Organization{
		Name:        req.Name,
		Description: req.Description,
		Owner:       userId,
	}
	err := s.repo.CreateOrganization(org)
	if err != nil {
		return nil, &errors.ApiError{
			Message:    "Failed to create organization",
			StatusCode: 500,
			Status:     errors.InternalServerError,
		}
	}
	return &dto.GetOrganizationResponse{
		OrgId:       org.OrgId,
		Name:        org.Name,
		Description: org.Description,
	}, nil
}
func (s *DefaultOrganizationService) AddUserToOrganization(orgId string, userId string) *errors.ApiError {
	err := s.repo.AddUserToOrganization(orgId, userId)
	if err != nil {
		return &errors.ApiError{
			Message:    "Client error",
			StatusCode: http.StatusBadRequest,
			Status:     errors.ValidationError,
		}
	}
	return nil
}

func NewOrganizationService(repo repository.OrganizationRepository) *DefaultOrganizationService {
	return &DefaultOrganizationService{repo: repo}
}

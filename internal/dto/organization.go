package dto

type GetOrganizationResponse struct {
	OrgId       string `json:"orgId"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateOrganizationRequest struct {
	Name        string `json:"name"binding:"required"`
	Description string `json:"description"`
}

type AddUserToOrganizationRequest struct {
	UserId string `json:"userId" binding:"required"`
}

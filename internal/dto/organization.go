package dto

type GetOrganizationResponse struct {
	OrgId       string `json:"orgId"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

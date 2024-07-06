package models

type Organization struct {
	OrgId       string `json:"orgId"gorm:"type:uuid;default:uuid_generate_v4();primarykey""`
	Name        string `json:"name"gorm:"type:varchar(100);not null"`
	Description string `json:"description"gorm:"type:varchar(100);not null"`
	Owner       string `json:"owner"gorm:"type:uuid;not null"`
}

type UserOrganization struct {
	OrgId  string `json:"orgId,omitempty"gorm:"type:uuid;references:Organization"`
	UserId string `json:"userId"gorm:"type:uuid;references:User"`
	Id     string `json:"id"gorm:"type:uuid;default:uuid_generate_v4();primarykey"`
}

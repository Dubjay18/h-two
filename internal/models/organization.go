package models

type Organization struct {
	orgId       string `json:"orgId"gorm:"type:uuid;default:uuid_generate_v4();primarykey""`
	Name        string `json:"name"gorm:"type:varchar(100);not null"`
	Description string `json:"description"gorm:"type:varchar(100);not null"`
}

package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserId    string `json:"userId"gorm:"type:uuid;default:uuid_generate_v4();primarykey""`
	Email     string `json:"email"gorm:"type:varchar(100);unique;not null"`
	Password  string `json:"password"gorm:"type:varchar(100);not null"`
	FirstName string `json:"firstName"gorm:"type:varchar(100);not null"`
	LastName  string `json:"lastName"gorm:"type:varchar(100);not null"`
	Phone     string `json:"phone"gorm:"type:varchar(100);not null"`
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &Organization{})
}

package model

import (
	"time"

	"gorm.io/gorm"
)

func (User) TableName() string {
	return "users"
}

type User struct {
	gorm.Model
	FirstName     string     `json:"first_name" gorm:"column:first_name"`
	LastName      string     `json:"last_name" gorm:"column:last_name"`
	Email         string     `json:"email" gorm:"column:email"`
	Password      string     `json:"password" gorm:"column:password"`
	Code          string     `json:"code" gorm:"column:code"`
	Mobile        string     `json:"mobile" gorm:"column:mobile"`
	IsActive      bool       `json:"is_active" gorm:"column:is_active"`
	EmailVerified bool       `json:"email_verified" gorm:"column:email_verified"`
	SignUpType    string     `json:"sign_up_type" gorm:"column:sign_up_type"`
	CreatedAt     time.Time  `json:"created_at" gorm:"column:created_at"`
	DeletedAt     *time.Time `json:"deleted_at" gorm:"column:deleted_at;default:null"`
	UpdatedAt     time.Time  `json:"updated_at" gorm:"column:updated_at"`
}

type UserLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

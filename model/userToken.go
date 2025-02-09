package model

import (
	"gorm.io/gorm"
	"time"
)

type Tabler interface {
	TableName() string
}

// TableName overrides the table name used by User to `profiles`
func (UserToken) TableName() string {
	return "user_token"
}

type UserToken struct {
	gorm.Model
	UserId    int       `gorm:"column:user_id;index:user_id_token_is_active"`
	Token     string    `gorm:"column:token;size:300;index:user_id_token_is_active;index:token_is_active"`
	IsActive  bool      `gorm:"column:is_active;default:true;index:user_id_token_is_active;index:token_is_active"`
	ExpiresAt time.Time `gorm:"column:expires_at"`
}

type UserDeviceToken struct {
	gorm.Model
	UserID      int    `json:"user_id"`
	DeviceToken string `json:"device_token"`
	EndpointARN string `json:"endpoint_arn"`
}

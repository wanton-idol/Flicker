package model

import (
	"gorm.io/gorm"
)

// TableName overrides the table name used by UserMedia to `profile_media`
func (UserMedia) TableName() string {
	return "profile_media"
}

type UserMedia struct {
	gorm.Model
	UserId        int     `json:"user_id" gorm:"column:user_id"`
	UserProfileId int     `json:"user_profile_id" gorm:"column:user_profile_id"`
	URL           string  `json:"url" gorm:"column:url"`
	OrderId       int     `json:"order_id" gorm:"column:order_id"`
	ImageText     string  `json:"image_text" gorm:"column:image_text"`
	Latitude      float64 `json:"latitude" gorm:"column:latitude"`
	Longitude     float64 `json:"longitude" gorm:"column:longitude"`
	City          string  `json:"city" gorm:"column:city"`
}

type MediaValues struct {
	Text      string  `json:"text"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type MediaOrderId struct {
	MediaID int `json:"media_id"`
	OrderID int `json:"order_id"`
}

package model

import (
	"gorm.io/gorm"
	"time"
)

type Event struct {
	gorm.Model
	UserId      int       `gorm:"user_id"`
	EventTime   time.Time `gorm:"event_time"`
	Type        string    `gorm:"type"`
	Description string    `gorm:"description"`
	Attendees   string    `gorm:"attendees"`
	ExpiresAt   time.Time `gorm:"expires_at"`
	Address1    string    `gorm:"address1"`
	Address2    string    `gorm:"address2"`
	City        string    `gorm:"city"`
	State       string    `gorm:"state"`
	Country     string    `gorm:"country"`
	Pincode     string    `gorm:"pincode"`
	Latitude    float32   `gorm:"latitude"`
	Longitude   float32   `gorm:"longitude"`
}

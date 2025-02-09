package model

import (
	"gorm.io/gorm"
	"time"
)

type UserVerificationOTP struct {
	gorm.Model
	PhoneNumber string    `json:"phoneNumber" gorm:"column:phone_number"`
	OTP         string    `json:"otp" gorm:"column:otp"`
	ExpiresAt   time.Time `gorm:"column:expires_at"`
}

type EmailVerification struct {
	gorm.Model
	UserId           int       `json:"userId" gorm:"column:user_id"`
	VerificationCode string    `json:"verificationCode" gorm:"column:verification_code"`
	IsVerified       bool      `json:"isVerified" gorm:"column:is_verified"`
	ExpiresAt        time.Time `json:"expires_at" gorm:"column:expires_at"`
}

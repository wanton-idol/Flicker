package dao

import (
	"errors"
	"fmt"
	"github.com/SuperMatch/model"
	"github.com/SuperMatch/zapLogger"
	"gorm.io/gorm"
)

//go:generate mockgen -package mocks -destination mocks/user_verification_otp_dao_mock.go github.com/SuperMatch/pkg/db/dao UserVerificationOTPRepository
type UserVerificationOTPRepository interface {
	Insert(userOtp model.UserVerificationOTP) error
	FindByPhoneAndOTP(phoneNumber string, otp string) (model.UserVerificationOTP, error)
}

type UserVerificationOTPDao struct {
	Connection gorm.DB
}

func (o *UserVerificationOTPDao) Insert(userOtp model.UserVerificationOTP) error {
	tx := o.Connection.Table("user_verification_otp").Create(&userOtp)

	if tx.Error != nil {
		zapLogger.Logger.Error("error inserting userOtp in user_verification_otp table")
		return tx.Error
	}

	return nil
}

func (o *UserVerificationOTPDao) FindByPhoneAndOTP(phoneNumber string, otp string) (model.UserVerificationOTP, error) {
	var userOtpDetails model.UserVerificationOTP
	err := o.Connection.Table("user_verification_otp").Where("phone_number = ? and otp = ?", phoneNumber, otp).Find(&userOtpDetails)
	if errors.Is(err.Error, gorm.ErrRecordNotFound) {
		zapLogger.Logger.Error(fmt.Sprintf("No record found for otp: %s", otp))
		return userOtpDetails, err.Error
	} else if err.Error != nil {
		zapLogger.Logger.Error(fmt.Sprintf("error in getting otp details for phone number: %s", phoneNumber))
		return userOtpDetails, err.Error
	}

	return userOtpDetails, nil
}

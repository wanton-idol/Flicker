package dao

import (
	"errors"
	"fmt"
	"github.com/SuperMatch/model"
	"github.com/SuperMatch/zapLogger"
	"gorm.io/gorm"
)

//go:generate mockgen -package mocks -destination mocks/emailVerification_dao_mock.go github.com/SuperMatch/pkg/db/dao EmailVerificationRepository
type EmailVerificationRepository interface {
	Insert(verificationDetails model.EmailVerification) error
	FindByVerificationCode(verificationCode string) (model.EmailVerification, error)
	UpdateEmailVerificationDetails(verificationDetails model.EmailVerification) (model.EmailVerification, error)
}

type EmailVerification struct {
	Connection gorm.DB
}

//func NewEmailVerificationRepository() EmailVerificationRepository {
//	return &EmailVerification{Connection: *db.GlobalOrm}
//}

func (e *EmailVerification) Insert(verificationDetails model.EmailVerification) error {
	err := e.Connection.Table("email_verification").Create(&verificationDetails)
	if err.Error != nil {
		zapLogger.Logger.Error("error in inserting verificationDetails in email_verification table")
		return err.Error
	}

	return nil
}

func (e *EmailVerification) FindByVerificationCode(verificationCode string) (model.EmailVerification, error) {
	var verificationDetails model.EmailVerification
	err := e.Connection.Table("email_verification").Where("verification_code = ? and is_verified = ?", verificationCode, false).Find(&verificationDetails)
	if errors.Is(err.Error, gorm.ErrRecordNotFound) {
		zapLogger.Logger.Error(fmt.Sprintf("No record found for verification_code: %s", verificationCode))
		return verificationDetails, err.Error
	} else if err.Error != nil {
		zapLogger.Logger.Error("error retrieving email_verification details")
		return verificationDetails, err.Error
	}

	return verificationDetails, nil

}

func (e *EmailVerification) UpdateEmailVerificationDetails(verificationDetails model.EmailVerification) (model.EmailVerification, error) {
	tx := e.Connection.Debug().Table("email_verification").Where("user_id = ?", verificationDetails.UserId).UpdateColumns(&verificationDetails)
	if tx.Error != nil {
		zapLogger.Logger.Error("error updating email_verification details")
		return verificationDetails, tx.Error
	}

	return verificationDetails, nil
}

package dao

import (
	"github.com/SuperMatch/model"
	"github.com/SuperMatch/pkg/db"
	"github.com/SuperMatch/zapLogger"
	"gorm.io/gorm"
)

type UserDeviceTokenDao interface {
	Insert(deviceToken model.UserDeviceToken) error
	GetDeviceTokensByUserId(userId int) ([]model.UserDeviceToken, error)
}

type UserDeviceTokenDaoImpl struct {
	Connection gorm.DB
}

func NewUserDeviceTokenDaoImpl() *UserDeviceTokenDaoImpl {
	return &UserDeviceTokenDaoImpl{Connection: *db.GlobalOrm}
}

func (d *UserDeviceTokenDaoImpl) Insert(deviceToken model.UserDeviceToken) error {
	err := d.Connection.Table("user_device_tokens").Create(&deviceToken)
	if err.Error != nil {
		zapLogger.Logger.Error("error inserting user device token in DB")
		return err.Error
	}

	return nil
}

func (d *UserDeviceTokenDaoImpl) GetDeviceTokensByUserId(userId int) ([]model.UserDeviceToken, error) {
	var deviceTokens []model.UserDeviceToken
	err := d.Connection.Table("user_device_tokens").Where("user_id = ?", userId).Find(&deviceTokens)
	if err.Error != nil {
		zapLogger.Logger.Error("error getting user device token in DB")
		return nil, err.Error
	}

	return deviceTokens, nil
}

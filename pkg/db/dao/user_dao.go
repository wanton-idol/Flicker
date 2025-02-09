package dao

import (
	"context"
	"github.com/SuperMatch/model"
	"github.com/SuperMatch/zapLogger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

//go:generate mockgen -package mocks -destination mocks/user_dao_mock.go github.com/SuperMatch/pkg/db/dao UserRepository
type UserRepository interface {
	Insert(ctx context.Context, user model.User) (model.User, error)
	FindById(ctx context.Context, id int64) (model.User, error)
	FindByEmail(ctx context.Context, email string) (model.User, error)
	FindByMobile(mobile string) (model.User, error)
	FindByMobileEmail(mobile, email string) ([]model.User, error)
	UpdateUser(user model.User) (model.User, error)
	DeleteUser(userID int, email string) error
}

// type UserDao struct {
// 	Connection db.Interface
// }

type UserDao struct {
	Connection gorm.DB
}

func (u *UserDao) Insert(user model.User) (model.User, error) {

	tx := u.Connection.Create(&user)

	if tx.Error != nil {
		return user, tx.Error
	}
	return user, nil
}

func (u *UserDao) FindById(externalId int) (model.User, error) {
	var user model.User
	err := u.Connection.Where("is_active = ?", true).Where("ID = ?", externalId).First(&user)

	if err != nil {
		return user, err.Error
	}
	return user, nil
}

func (u *UserDao) FindByEmail(email string) (model.User, error) {
	var user model.User
	err := u.Connection.Where("is_active = ?", true).Where("email = ?", email).First(&user)
	return user, err.Error
}

func (u *UserDao) FindByMobile(mobile string) (model.User, error) {
	var user model.User
	err := u.Connection.Where("is_active = ?", true).Where("mobile = ? ", mobile).First(&user)
	return user, err.Error
}

func (u *UserDao) FindByMobileEmail(mobile, email string) ([]model.User, error) {
	var user []model.User
	err := u.Connection.Where("is_active = ?", true).Where("mobile = ?  or email = ?", mobile, email).Find(&user)
	return user, err.Error
}

func (u *UserDao) UpdateUser(user model.User) (model.User, error) {
	tx := u.Connection.Table("users").Where("is_active = ?", true).Where("ID = ?", user.ID).UpdateColumns(&user)
	if tx.Error != nil {
		zapLogger.Logger.Error("error updating user", zap.Error(tx.Error))
		return user, tx.Error
	}

	return user, nil
}

func (u *UserDao) DeleteUser(userID int, email string) error {
	tx := u.Connection.Debug().Table("users").Where("ID = ?", userID).Update("email", email)
	if tx.Error != nil {
		zapLogger.Logger.Error("error deleting user", zap.Error(tx.Error))
		return tx.Error
	}

	return nil
}

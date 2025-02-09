package dao

import (
	"context"
	"time"

	"github.com/SuperMatch/model"
	"gorm.io/gorm"
)

//go:generate mockgen -package mocks -destination mocks/user_token_dao_mock.go github.com/SuperMatch/pkg/db/dao UserTokenRepository
type UserTokenRepository interface {
	Insert(ctx context.Context, userToken model.UserToken) (model.UserToken, error)
	FindByUserId(ctx context.Context, userID int64) (model.UserToken, error)
	FindByToken(ctx context.Context, token string) (model.UserToken, error)
	RemoveSession(ctx context.Context, userId int64) error
	FindByTokenAndUserId(ctx context.Context, token string, userId int) (model.UserToken, error)
}

type UserTokenDao struct {
	Connection gorm.DB
}

func (d *UserTokenDao) Insert(ctx context.Context, userToken model.UserToken) (model.UserToken, error) {

	tx := d.Connection.Create(&userToken)
	return userToken, tx.Error
}

func (d *UserTokenDao) FindByUserId(ctx context.Context, userId int) ([]model.UserToken, error) {
	var userToken []model.UserToken
	err := d.Connection.Where("is_active = ?", true).Where("user_id = ?", userId).First(&userToken)
	return userToken, err.Error
}

func (d *UserTokenDao) FindByToken(ctx context.Context, token string) (model.UserToken, error) {
	var userToken model.UserToken
	err := d.Connection.Where("is_active = ?", true).Where("token = ?", token).Where("expired_at > ?", time.Now()).First(&userToken)
	return userToken, err.Error
}

func (d *UserTokenDao) RemoveSession(ctx context.Context, userId int64) error {
	err := d.Connection.Where("user_id = ?", userId).Where("is_active = ?", true).UpdateColumn("is_active", false)
	return err.Error
}

func (d *UserTokenDao) FindByTokenAndUserId(ctx context.Context, token string, userId int) (model.UserToken, error) {
	var userToken model.UserToken
	err := d.Connection.Debug().Where("is_active = ?", true).Where("token = ?", token).Where("user_id = ?", userId).Where("expires_at > ?", time.Now()).First(&userToken)
	return userToken, err.Error
}

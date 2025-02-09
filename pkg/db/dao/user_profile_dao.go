package dao

import (
	"context"
	"errors"
	"fmt"
	"github.com/SuperMatch/model"
	"github.com/SuperMatch/pkg/db"
	"github.com/SuperMatch/zapLogger"
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

//go:generate mockgen -package mocks -destination mocks/user_profile_dao_mock.go github.com/SuperMatch/pkg/db/dao UserProfileRepository
type UserProfileRepository interface {
	CreateUserProfile(ctx context.Context, userProfile model.UserProfile) (model.UserProfile, error)
	UpdateUserProfile(ctx context.Context, userProfile model.UserProfile) (model.UserProfile, error)
	FindById(ctx context.Context, profileId int) (model.UserProfile, error)
	FindByUserId(ctx context.Context, userId int) (model.UserProfile, error)
	UpdateProfileByMap(ctx context.Context, userProfileMap map[string]interface{}) (map[string]interface{}, error)
}

type UserProfile struct {
	Connection gorm.DB
}

func NewUserProfileRepository() UserProfileRepository {
	return &UserProfile{Connection: *db.GlobalOrm}
}

func (p *UserProfile) CreateUserProfile(ctx context.Context, userProfile model.UserProfile) (model.UserProfile, error) {
	zapLogger.Logger.Info("CreateUserProfile is started...")
	tx := p.Connection.Debug().Create(&userProfile)
	if tx.Error != nil {
		zapLogger.Logger.Error("error while creating user profile", zap.Error(tx.Error))
		return userProfile, tx.Error
	}

	return userProfile, tx.Error
}

func (p *UserProfile) UpdateUserProfile(ctx context.Context, userProfile model.UserProfile) (model.UserProfile, error) {
	zapLogger.Logger.Info("UpdateUserProfile is started...")

	//update user profile
	tx := p.Connection.Debug().Where("user_id = ?", userProfile.UserId).UpdateColumns(&userProfile)
	if tx.Error != nil {
		sentry.CaptureException(tx.Error)
		zapLogger.Logger.Error("error while updating user profile", zap.Error(tx.Error))
		return userProfile, tx.Error
	}

	zapLogger.Logger.Info(fmt.Sprintf("Updated user profile for user id: %d", userProfile.UserId))
	return userProfile, tx.Error
}

func (p *UserProfile) FindById(ctx context.Context, profileId int) (model.UserProfile, error) {
	zapLogger.Logger.Info("FindById in userProfile Dao is started...")

	var userProfile model.UserProfile
	tx := p.Connection.Where("deleted_at IS NULL").Where("ID = ?", profileId).First(&userProfile)

	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return userProfile, tx.Error
	}

	if tx.Error != nil {
		zapLogger.Logger.Error("error while getting user profile: ", zap.Error(tx.Error))
		return userProfile, tx.Error
	}
	return userProfile, nil
}

func (p *UserProfile) FindByUserId(ctx context.Context, userId int) (model.UserProfile, error) {
	zapLogger.Logger.Info("FindByUserId in userProfile Dao is started...")

	var userProfile model.UserProfile
	tx := p.Connection.Where("deleted_at IS NULL").Where("user_id = ?", userId).First(&userProfile)

	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return userProfile, tx.Error
	}

	if tx.Error != nil {
		zapLogger.Logger.Error("error while getting user profile: ", zap.Error(tx.Error))
		return userProfile, tx.Error
	}
	return userProfile, nil
}

func (p *UserProfile) UpdateProfileByMap(ctx context.Context, userProfileMap map[string]interface{}) (map[string]interface{}, error) {
	zapLogger.Logger.Info("UpdateProfileByMap in userProfile Dao is started...")

	//update user profile sql raw query

	tx := p.Connection.Debug().Table("user_profile").Where("ID", userProfileMap["ID"]).UpdateColumns(userProfileMap)
	if tx.Error != nil {
		zapLogger.Logger.Error("error while updating user profile: ", zap.Error(tx.Error))
		return userProfileMap, tx.Error
	}
	return userProfileMap, nil
}

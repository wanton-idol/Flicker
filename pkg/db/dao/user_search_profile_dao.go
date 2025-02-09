package dao

import (
	"context"
	"errors"
	"github.com/SuperMatch/model"
	"github.com/SuperMatch/pkg/db"
	"github.com/SuperMatch/zapLogger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

//go:generate mockgen -package mocks -destination mocks/user_search_profile_dao_mock.go github.com/SuperMatch/pkg/db/dao UserSearchProfileRepository
type UserSearchProfileRepository interface {
	CreateUserSearchProfile(ctx context.Context, userProfile map[string]interface{}) (model.UserSearchProfile, error)
	UpdateUserSearchProfile(ctx context.Context, userProfile map[string]interface{}) (model.UserSearchProfile, error)
	FindByProfileId(ctx context.Context, profileId int) (model.UserSearchProfile, error)
	FindByUserId(ctx context.Context, userId int) (model.UserSearchProfile, error)
}

type UserSearchProfile struct {
	Connection            gorm.DB
	UserProfileRepository UserProfileRepository
}

func NewUserSearchProfile() UserSearchProfileRepository {
	return &UserSearchProfile{
		Connection:            *db.GlobalOrm,
		UserProfileRepository: NewUserProfileRepository(),
	}
}

func (u *UserSearchProfile) CreateUserSearchProfile(ctx context.Context, userSearchProfile map[string]interface{}) (model.UserSearchProfile, error) {
	searchProfile := model.UserSearchProfile{}

	tx := u.Connection.Debug().Table("user_search_profile").Create(&userSearchProfile)
	if tx.Error != nil {
		zapLogger.Logger.Error("error in creating user search profile in db.", zap.Error(tx.Error))
		return searchProfile, tx.Error
	}

	return searchProfile, nil
}

func (u *UserSearchProfile) UpdateUserSearchProfile(ctx context.Context, userSearchProfile map[string]interface{}) (model.UserSearchProfile, error) {
	searchProfile := model.UserSearchProfile{}
	conditions := make(map[string]interface{})
	conditions["user_profile_id"] = userSearchProfile["user_profile_id"]

	tx := u.Connection.Model(searchProfile).Where(conditions).Updates(userSearchProfile)
	if tx.Error != nil {
		zapLogger.Logger.Error("error in updating user search profile in db.", zap.Error(tx.Error))
		return searchProfile, tx.Error
	}

	return searchProfile, nil
}

func (u *UserSearchProfile) FindByProfileId(ctx context.Context, userProfileId int) (model.UserSearchProfile, error) {
	zapLogger.Logger.Info("FindByProfileId in userSearchProfile Dao is started...")

	//mp := make(map[string]interface{})
	var userSearchProfile model.UserSearchProfile
	tx := u.Connection.Table("user_search_profile").Where("deleted_at IS NULL").Where("user_profile_id = ?", userProfileId).First(&userSearchProfile)

	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		zapLogger.Logger.Info("user_search_profile not found for user_profile_id", zap.Any("user_profile_id", userProfileId))
		return userSearchProfile, tx.Error
	}

	if tx.Error != nil {
		zapLogger.Logger.Error("error while getting user profile:", zap.Error(tx.Error))
		return userSearchProfile, tx.Error
	}

	return userSearchProfile, nil
}

func (u *UserSearchProfile) FindByUserId(ctx context.Context, userId int) (model.UserSearchProfile, error) {
	zapLogger.Logger.Info("FindByUserId in userSearchProfile Dao is started...")

	var userSearchProfile model.UserSearchProfile
	tx := u.Connection.Where("deleted_at IS NULL").Where("user_id = ?", userId).First(&userSearchProfile)

	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return userSearchProfile, tx.Error
	}

	if tx.Error != nil {
		zapLogger.Logger.Error("error while getting user profile: ", zap.Error(tx.Error))
		return userSearchProfile, tx.Error
	}
	return userSearchProfile, nil
}

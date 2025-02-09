package dao

import (
	"github.com/SuperMatch/model"
	"github.com/SuperMatch/pkg/db"
	"github.com/SuperMatch/zapLogger"
	"gorm.io/gorm"
)

//go:generate mockgen -package mocks -destination mocks/interests_dao_mock.go github.com/SuperMatch/pkg/db/dao InterestsDao
type InterestsDao interface {
	GetInterestsList() ([]model.Interests, error)
	CreateUserInterests(userInterests []model.UserInterests, userID int) ([]model.UserInterests, error)
	GetUserInterests(userID int) ([]model.UserInterests, error)
	UpdateUserInterests(userInterests []model.UserInterests, userID int) ([]model.UserInterests, error)
}

type InterestsDaoImpl struct {
	Connection gorm.DB
}

func NewInterestsDaoImpl() *InterestsDaoImpl {
	return &InterestsDaoImpl{Connection: *db.GlobalOrm}
}

func (i *InterestsDaoImpl) GetInterestsList() ([]model.Interests, error) {
	var interestDB []model.Interests

	err := i.Connection.Table("interests_details").Find(&interestDB)
	if err.Error != nil {
		zapLogger.Logger.Error("[GetInterestsList] error in retrieving data from interest_details table")
		return interestDB, err.Error
	}

	return interestDB, nil
}

func (i *InterestsDaoImpl) CreateUserInterests(userInterests []model.UserInterests, userID int) ([]model.UserInterests, error) {
	err := i.Connection.Table("user_interests").Create(&userInterests)
	if err.Error != nil {
		zapLogger.Logger.Error("[CreateUserInterests] error in creating user_interest in db")
		return userInterests, err.Error
	}

	return userInterests, nil
}

func (i *InterestsDaoImpl) GetUserInterests(userID int) ([]model.UserInterests, error) {
	var userInterests []model.UserInterests
	err := i.Connection.Table("user_interests").Where("user_id = ?", userID).Find(&userInterests)
	if err.RowsAffected == 0 {
		zapLogger.Logger.Error("[GetUserInterests] no user_interests found")
		return userInterests, gorm.ErrRecordNotFound
	}
	if err.Error != nil {
		zapLogger.Logger.Error("[GetUserInterests] error in getting user_interests from database")
		return userInterests, err.Error
	}

	return userInterests, nil
}

func (i *InterestsDaoImpl) UpdateUserInterests(userInterests []model.UserInterests, userID int) ([]model.UserInterests, error) {
	err := i.Connection.Table("user_interests").Where("user_id = ?", userID).Save(&userInterests)
	if err.Error != nil {
		zapLogger.Logger.Error("[UpdateUserInterests] error in updating user_interests in db")
		return userInterests, err.Error
	}

	return userInterests, nil
}

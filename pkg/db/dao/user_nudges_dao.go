package dao

import (
	"errors"
	"fmt"
	"github.com/SuperMatch/model"
	"github.com/SuperMatch/pkg/db"
	"github.com/SuperMatch/zapLogger"
	"gorm.io/gorm"
)

//go:generate mockgen -package mocks -destination mocks/user_nudges_dao_mock.go github.com/SuperMatch/pkg/db/dao UserNudgesDao
type UserNudgesDao interface {
	GetNudgesDB() ([]model.Nudge, error)
	CreateUserNudgesDB(userNudges model.UserNudge) (model.UserNudge, error)
	GetUserNudgesDB(userID int) ([]model.UserNudge, error)
	GetUserNudgeById(id int) (model.UserNudge, error)
	UpdateUserNudge(userNudge model.UserNudge, id int) (model.UserNudge, error)
	DeleteUserNudge(id int) error
}

type UserNudgesDaoImpl struct {
	Connection gorm.DB
}

func NewUserNudgesDaoImpl() *UserNudgesDaoImpl {
	return &UserNudgesDaoImpl{Connection: *db.GlobalOrm}
}

func (un *UserNudgesDaoImpl) GetNudgesDB() ([]model.Nudge, error) {
	var nudges []model.Nudge
	err := un.Connection.Table("nudges").Find(&nudges)
	if err.Error != nil {
		zapLogger.Logger.Error("[GetNudges] error in getting nudges list")
		return nudges, err.Error
	}

	return nudges, nil
}

func (un *UserNudgesDaoImpl) CreateUserNudgesDB(userNudges model.UserNudge) (model.UserNudge, error) {
	err := un.Connection.Table("user_nudges").Create(&userNudges)
	if err.Error != nil {
		zapLogger.Logger.Error("[CreateUserNudges] error in creating user_nudges in db")
		return userNudges, err.Error
	}

	return userNudges, nil
}

func (un *UserNudgesDaoImpl) GetUserNudgesDB(userID int) ([]model.UserNudge, error) {
	var userNudges []model.UserNudge
	err := un.Connection.Table("user_nudges").Where("user_id = ?", userID).Find(&userNudges)
	if errors.Is(err.Error, gorm.ErrRecordNotFound) {
		zapLogger.Logger.Error(fmt.Sprintf("no record found for userID: %d in DB", userID))
		return userNudges, err.Error
	} else if err.Error != nil {
		zapLogger.Logger.Error("error in retrieving userNudges from DB")
		return userNudges, err.Error
	}

	return userNudges, nil
}

func (un *UserNudgesDaoImpl) GetUserNudgeById(id int) (model.UserNudge, error) {
	var userNudge model.UserNudge
	err := un.Connection.Table("user_nudges").Where("id = ?", id).First(&userNudge)
	if errors.Is(err.Error, gorm.ErrRecordNotFound) {
		zapLogger.Logger.Error(fmt.Sprintf("no record found for id: %d in DB", id))
		return userNudge, err.Error
	} else if err.Error != nil {
		zapLogger.Logger.Error("error in retrieving userNudges from DB")
		return userNudge, err.Error
	}

	return userNudge, nil
}

func (un *UserNudgesDaoImpl) UpdateUserNudge(userNudge model.UserNudge, id int) (model.UserNudge, error) {
	err := un.Connection.Table("user_nudges").Select("*").Where("id = ?", id).Updates(&userNudge)
	if err.Error != nil {
		zapLogger.Logger.Error(fmt.Sprintf("error in updating userNudge with id: %d", id))
		return userNudge, err.Error
	}

	return userNudge, nil
}

func (un *UserNudgesDaoImpl) DeleteUserNudge(id int) error {
	err := un.Connection.Table("user_nudges").Where("id = ?", id).Delete(&model.UserNudge{})
	if err.Error != nil {
		zapLogger.Logger.Error(fmt.Sprintf("error in deleting userNudge with id: %d", id))
		return err.Error
	}

	return nil
}

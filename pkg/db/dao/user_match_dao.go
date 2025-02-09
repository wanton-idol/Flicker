package dao

import (
	"context"

	"github.com/SuperMatch/model"
	"github.com/SuperMatch/pkg/db"
	"gorm.io/gorm"
)

//go:generate mockgen -package mocks -destination mocks/user_match_dao_mock.go github.com/SuperMatch/pkg/db/dao UserMatchDao
type UserMatchDao interface {
	Insert(ctx context.Context, userMatch model.UserMatch) (model.UserMatch, error)
	InsertMany(ctx context.Context, userMatch []model.UserMatch) ([]model.UserMatch, error)
	FindByUserId(ctx context.Context, userId int) ([]model.UserMatch, error)
	FindByUserIdMatchId(ctx context.Context, userId, matchId int) (model.UserMatch, error)
	DeleteByUserID(ctx context.Context, userId int, matchId int) error
}

type UserMatchDaoImpl struct {
	Connection gorm.DB
}

func NewUserMatchDaoImpl() *UserMatchDaoImpl {
	return &UserMatchDaoImpl{Connection: *db.GlobalOrm}
}

func (UserMatchDao *UserMatchDaoImpl) Insert(ctx context.Context, userMatch model.UserMatch) (model.UserMatch, error) {

	tx := UserMatchDao.Connection.Create(&userMatch)

	if tx.Error != nil {
		return userMatch, tx.Error
	}
	return userMatch, nil
}

func (UserMatchDao *UserMatchDaoImpl) InsertMany(ctx context.Context, userMatch []model.UserMatch) ([]model.UserMatch, error) {
	tx := UserMatchDao.Connection.Create(&userMatch)
	if tx.Error != nil {
		return userMatch, tx.Error
	}
	return userMatch, nil
}

func (UserMatchDao *UserMatchDaoImpl) FindByUserId(ctx context.Context, userId int) ([]model.UserMatch, error) {
	var userMatches []model.UserMatch
	tx := UserMatchDao.Connection.Where("user_id = ?", userId).Find(&userMatches)

	if tx.Error != nil {
		return userMatches, tx.Error
	}
	return userMatches, nil
}

func (UserMatchDao *UserMatchDaoImpl) FindByUserIdMatchId(ctx context.Context, userId, matchId int) (model.UserMatch, error) {
	var userMatch model.UserMatch
	tx := UserMatchDao.Connection.Where("user_id = ? and match_id = ?", userId, matchId).First(&userMatch)

	if tx.Error != nil {
		return userMatch, tx.Error
	}
	return userMatch, nil
}

func (UserMatchDao *UserMatchDaoImpl) DeleteByUserID(ctx context.Context, userId int, matchId int) error {
	tx := UserMatchDao.Connection.Where("user_id = ? and match_id = ?", userId, matchId).Or("user_id = ? and match_id = ?", matchId, userId).Delete(&model.UserMatch{})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

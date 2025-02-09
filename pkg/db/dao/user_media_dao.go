package dao

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/SuperMatch/zapLogger"
	"go.uber.org/zap"

	"github.com/SuperMatch/model"
	"github.com/SuperMatch/pkg/db"
	"gorm.io/gorm"
)

//go:generate mockgen -package mocks -destination mocks/user_media_dao_mock.go github.com/SuperMatch/pkg/db/dao UserMediaRepository
type UserMediaRepository interface {
	Insert(ctx context.Context, userMedia model.UserMedia) (model.UserMedia, error)
	FindById(ctx context.Context, id int) (model.UserMedia, error)
	FindByIdAndUserId(ctx context.Context, id, userId int) (*model.UserMedia, error)
	FindByUserId(ctx context.Context, userId int) ([]model.UserMedia, error)
	DeleteById(ctx context.Context, id int, userId int) error
	FindFirstGroupByUserID(ctc context.Context, userIDs []int) ([]UserMatchUserMediaDTO, error)
	FindByUserIDOrderID(userID, orderID int) (model.UserMedia, error)
	UpdateProfileMedia(ctx context.Context, user model.UserMedia, mediaID int) (model.UserMedia, error)
	FindByUserIDs(userIDs []int) ([]model.UserMedia, error)
}

type UserMedia struct {
	Connection gorm.DB
}

func NewUserMediaRepository() UserMediaRepository {
	return &UserMedia{Connection: *db.GlobalOrm}
}

func (u *UserMedia) Insert(ctx context.Context, userMedia model.UserMedia) (model.UserMedia, error) {

	tx := u.Connection.Create(&userMedia)
	if tx.Error != nil {
		return userMedia, tx.Error
	}
	return userMedia, nil
}

func (u *UserMedia) FindById(ctx context.Context, id int) (model.UserMedia, error) {
	var usermedia model.UserMedia
	tx := u.Connection.Debug().Where("deleted_at IS NULL").Where("ID = ?", id).Find(&usermedia)

	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return usermedia, nil
	}

	if tx.Error != nil {
		zapLogger.Logger.Error("error while getting user media: ", zap.Error(tx.Error))
		return usermedia, tx.Error
	}
	return usermedia, nil
}

func (u *UserMedia) FindByIdAndUserId(ctx context.Context, id, userId int) (*model.UserMedia, error) {
	var usermedia model.UserMedia

	tx := u.Connection.Where("deleted_at IS NULL").Where("ID = ? and user_id = ?", id, userId).First(&usermedia)

	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		zapLogger.Logger.Info(fmt.Sprintf("no media record found in DB for userID: %d and mediaID: %d\n", userId, id))
		return &usermedia, tx.Error
	}

	if tx.Error != nil {
		zapLogger.Logger.Error("error while getting user media: ", zap.Error(tx.Error))
		return &usermedia, tx.Error
	}
	return &usermedia, nil
}

func (u *UserMedia) FindByUserId(ctx context.Context, userId int) ([]model.UserMedia, error) {

	var usermedia []model.UserMedia
	tx := u.Connection.Debug().Where("deleted_at IS NULL").Where("user_id = ?", userId).Find(&usermedia)

	if tx.RowsAffected == 0 {
		zapLogger.Logger.Error(fmt.Sprintf("no media record found in DB for userID: %d", userId), zap.Error(tx.Error))
		return usermedia, gorm.ErrRecordNotFound
	}

	if tx.Error != nil {
		zapLogger.Logger.Error("error while getting user media: ", zap.Error(tx.Error))
		return usermedia, tx.Error
	}

	return usermedia, nil
}

func (u *UserMedia) DeleteById(ctx context.Context, id int, userId int) error {

	tx := u.Connection.Where("ID = ? and user_id = ?", id, userId).Delete(&model.UserMedia{})

	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil
	}
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

type UserMatchUserMediaDTO struct {
	ID        int        `json:"ID"`
	UserId    int        `json:"user_id"`
	MatchId   int        `json:"match_id"`
	CreatedAt time.Time  `json:"created_at"`
	OrderId   int        `json:"order_id"`
	MediaId   int        `json:"media_id"`
	URL       string     `json:"url"`
	ChatId    *string    `json:"chat_id"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func (u *UserMedia) FindFirstGroupByUserID(ctx context.Context, userIDs []int) ([]UserMatchUserMediaDTO, error) {

	var userMedia []UserMatchUserMediaDTO

	query := "SELECT um.ID,um.user_id,um.match_id,um.created_at,pf.order_id,pf.ID as media_id,pf.url,pf.deleted_at,um.chat_id from user_match as um LEFT JOIN profile_media as pf ON um.match_id = pf.user_id where pf.order_id = 1 and um.match_id in (?) order by um.ID;"

	tx := u.Connection.Debug().Raw(query, userIDs).Find(&userMedia)

	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return userMedia, nil
	}

	if tx.Error != nil {
		zapLogger.Logger.Error("error while getting user media: ", zap.Error(tx.Error))
		return userMedia, tx.Error
	}

	return userMedia, nil
}

func (u *UserMedia) UpdateProfileMedia(ctx context.Context, media model.UserMedia, mediaID int) (model.UserMedia, error) {
	tx := u.Connection.Debug().Table("profile_media").Where("ID=? and user_id=?", mediaID, media.UserId).UpdateColumns(&media)
	if tx.Error != nil {
		zapLogger.Logger.Error("error in updating profile media")
		return media, tx.Error
	}

	return media, nil
}

func (u *UserMedia) FindByUserIDOrderID(userID, orderID int) (model.UserMedia, error) {

	var userMedia model.UserMedia
	tx := u.Connection.Table("profile_media").Where("deleted_at IS NULL").Where("user_id = ? and order_id = ?", userID, orderID).Find(&userMedia)

	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return userMedia, nil
	}

	if tx.Error != nil {
		zapLogger.Logger.Error("error while getting user media by userID and orderID: ", zap.Error(tx.Error))
		return userMedia, tx.Error
	}

	return userMedia, nil
}

func (u *UserMedia) FindByUserIDs(userIDs []int) ([]model.UserMedia, error) {
	var userMedia []model.UserMedia
	tx := u.Connection.Table("profile_media").Where("deleted_at IS NULL").Where("user_id in (?) and order_id = 1", userIDs).Find(&userMedia)
	if tx.Error != nil {
		zapLogger.Logger.Error("error while getting user medias by userIDs: ", zap.Error(tx.Error))
		return userMedia, tx.Error
	}

	return userMedia, nil
}

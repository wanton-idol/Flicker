package dao

import (
	"context"
	"errors"
	"fmt"
	"github.com/SuperMatch/model"
	"github.com/SuperMatch/pkg/db"
	"github.com/SuperMatch/zapLogger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

//go:generate mockgen -package mocks -destination mocks/advanced_filter_dao_mock.go github.com/SuperMatch/pkg/db/dao AdvancedFilterRepository
type AdvancedFilterRepository interface {
	CreateAdvancedFilter(context context.Context, advancedFilter model.AdvancedFilter) error
	UpdateAdvancedFilter(advancedFilter model.AdvancedFilter) (model.AdvancedFilter, error)
	FindByUserID(userId int) (model.AdvancedFilter, error)
}

type AdvancedFilter struct {
	Connection gorm.DB
}

func NewAdvancedFilterRepository() AdvancedFilterRepository {
	return &AdvancedFilter{Connection: *db.GlobalOrm}
}

func (a *AdvancedFilter) CreateAdvancedFilter(context context.Context, advancedFilter model.AdvancedFilter) error {
	zapLogger.Logger.Info("CreateAdvancedFilter is started...")
	tx := a.Connection.Table("advanced_filters").Create(&advancedFilter)
	if tx != nil {
		zapLogger.Logger.Error("error while creating advanced filters:", zap.Error(tx.Error))
		return tx.Error
	}

	zapLogger.Logger.Info(fmt.Sprintf("Inserted ID: %d", advancedFilter.ID))

	return tx.Error
}

func (a *AdvancedFilter) UpdateAdvancedFilter(advancedFilter model.AdvancedFilter) (model.AdvancedFilter, error) {
	zapLogger.Logger.Info("UpdateAdvancedFilter is started...")
	tx := a.Connection.Debug().Where("user_id = ?", advancedFilter.UserId).Where("user_profile_id = ?", advancedFilter.UserProfileId).UpdateColumns(&advancedFilter)
	if tx.Error != nil {
		zapLogger.Logger.Error("error while updating advanced filters:", zap.Error(tx.Error))
		return advancedFilter, tx.Error
	}

	zapLogger.Logger.Info(fmt.Sprintf("updated advanced filter for id: %d", advancedFilter.ID))

	return advancedFilter, tx.Error
}

func (a *AdvancedFilter) FindByUserID(userId int) (model.AdvancedFilter, error) {
	zapLogger.Logger.Info("FindByUserId in AdvancedFilter Dao is started...")

	var advancedFilter model.AdvancedFilter
	tx := a.Connection.Where("deleted_at IS NULL").Where("user_id = ?", userId).First(&advancedFilter)

	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return advancedFilter, tx.Error
	}

	if tx.Error != nil {
		zapLogger.Logger.Error("error while getting advanced filters: ", zap.Error(tx.Error))
		return advancedFilter, tx.Error
	}
	return advancedFilter, nil
}

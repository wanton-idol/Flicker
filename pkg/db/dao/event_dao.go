package dao

import (
	"github.com/SuperMatch/model"
	"github.com/SuperMatch/pkg/db"
	"github.com/SuperMatch/zapLogger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type EventRepository interface {
	InsertEvent(event model.Event) (model.Event, error)
	GetEventById(id int) (model.Event, error)
	GetEventByUserId(userId int) ([]model.Event, error)
	GetEventsByUserIdAndEventID(userId int, eventId []int) ([]model.Event, error)
	UpdateEvent(event model.Event) (model.Event, error)
	DeleteEvent(userID, eventID int) error
}

type EventRepositoryImpl struct {
	Connection gorm.DB
}

func NewEventRepositoryImpl() *EventRepositoryImpl {
	return &EventRepositoryImpl{Connection: *db.GlobalOrm}
}

func (repo *EventRepositoryImpl) InsertEvent(event model.Event) (model.Event, error) {
	tx := repo.Connection.Create(&event)

	if tx.Error != nil {
		zapLogger.Logger.Error("EventRepository insert event failed with error %", zap.Error(tx.Error))
		return event, tx.Error
	}
	return event, nil
}

func (repo *EventRepositoryImpl) GetEventById(id int) (model.Event, error) {

	var event model.Event
	tx := repo.Connection.Where("deleted_at IS NULL").Where("ID = ?", id).Find(&event)

	if tx.Error != nil {
		zapLogger.Logger.Error("EventRepository get event failed with error %", zap.Error(tx.Error))
		return event, tx.Error
	}

	return event, nil
}

func (repo *EventRepositoryImpl) GetEventByUserId(userId int) ([]model.Event, error) {
	var events []model.Event
	tx := repo.Connection.Where("deleted_at IS NULL").Where("user_id = ?", userId).Find(&events)

	if tx.Error != nil {
		zapLogger.Logger.Error("EventRepository get events failed with error %", zap.Error(tx.Error))
		return events, tx.Error
	}
	return events, nil
}

func (repo *EventRepositoryImpl) GetEventsByUserIdAndEventID(userId int, eventId []int) ([]model.Event, error) {
	var events []model.Event
	tx := repo.Connection.Where("deleted_at IS NULL").Where("user_id = ?", userId).Where("event_id in ", eventId).Find(&events)

	if tx.Error != nil {
		zapLogger.Logger.Error("EventRepository get events failed with error %", zap.Error(tx.Error))
		return events, tx.Error
	}
	return events, nil
}

func (repo *EventRepositoryImpl) UpdateEvent(event model.Event) (model.Event, error) {
	tx := repo.Connection.Model(&event).Where("deleted_at IS NULL").Where("id = ?", event.ID).Updates(&event)

	if tx.Error != nil {
		zapLogger.Logger.Error("EventRepository update event failed with error %", zap.Error(tx.Error))
		return event, tx.Error
	}
	return event, nil
}

func (repo *EventRepositoryImpl) DeleteEvent(userID, eventID int) error {
	tx := repo.Connection.Where("deleted_at IS NULL").Where("user_id = ?", userID).Where("id = ?", eventID).Delete(&model.Event{})

	if tx.Error != nil {
		zapLogger.Logger.Error("EventRepository delete event failed with error %", zap.Error(tx.Error))
		return tx.Error
	}
	return nil
}

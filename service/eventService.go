package service

import (
	"database/sql"
	"encoding/json"
	"gorm.io/gorm"
	"time"

	"github.com/SuperMatch/model"
	"github.com/SuperMatch/model/dto"
	elasticsearchPkg "github.com/SuperMatch/model/elasticSearch"
	"github.com/SuperMatch/pkg/db/dao"
	"github.com/SuperMatch/pkg/elasticSeach"
	"github.com/SuperMatch/zapLogger"
	"go.uber.org/zap"
)

type EventService interface {
	CreateUserEvent(createEventDTO dto.CreateEventDTO, userProfile model.UserProfile) (dto.CreateEventDTO, error)
	GetUserEvent(userProfile model.UserProfile) ([]dto.EventResponseDTO, error)
	SearchEvents(userProfile model.UserProfile, page model.Pagination, filters dto.EventFilterDTO) ([]dto.EventResponseDTO, error)
	CreateEventIndex() error
	createEventElasticDTO(event model.Event) (elasticsearchPkg.Event, error)
	UpdateUserEvent(createEventDTO dto.CreateEventDTO, userProfile model.UserProfile) (dto.CreateEventDTO, error)
	DeleteUserEvent(userID, eventID int) error
}

type EventServiceImpl struct {
	EventRepository dao.EventRepository
	EventIndexer    elasticSeach.EventIndexer
}

func NewEventServiceImpl() EventService {
	return &EventServiceImpl{
		EventRepository: dao.NewEventRepositoryImpl(),
		EventIndexer:    elasticSeach.NewEventIndexerImpl(),
	}
}

func (e *EventServiceImpl) CreateEventIndex() error {
	err := e.EventIndexer.CreateIndex()
	return err
}

func (e *EventServiceImpl) CreateUserEvent(createEventDTO dto.CreateEventDTO, userProfile model.UserProfile) (dto.CreateEventDTO, error) {
	event := model.Event{
		UserId:      userProfile.UserId,
		EventTime:   createEventDTO.EventTime,
		Type:        createEventDTO.Type,
		Description: createEventDTO.Description,
		Attendees:   createEventDTO.Attendees,
		ExpiresAt:   createEventDTO.ExpiresAt,
		Address1:    createEventDTO.Location.Address1,
		Address2:    createEventDTO.Location.Address2,
		City:        createEventDTO.Location.City,
		State:       createEventDTO.Location.State,
		Pincode:     createEventDTO.Location.Pincode,
		Latitude:    createEventDTO.Location.Latitude,
		Longitude:   createEventDTO.Location.Longitude,
	}
	event, err := e.EventRepository.InsertEvent(event)
	if err != nil {
		return dto.CreateEventDTO{}, err
	}
	createEventDTO.ID = event.ID

	//index event to elasticSearch
	elasticDTO, _ := e.createEventElasticDTO(event)
	bytes, err := json.Marshal(elasticDTO)
	if err != nil {
		zapLogger.Logger.Error("Error while marshalling event to elasticDTO", zap.Error(err))
		return dto.CreateEventDTO{}, err
	}
	zapLogger.Logger.Info("Event marshalled to elasticDTO", zap.String("event", string(bytes)))

	err = e.EventIndexer.IndexUserEvent(elasticDTO, bytes)
	return createEventDTO, err

}

func (e *EventServiceImpl) GetUserEvent(userProfile model.UserProfile) ([]dto.EventResponseDTO, error) {
	events, err := e.EventRepository.GetEventByUserId(userProfile.UserId)

	if err != nil {
		return nil, err
	}

	var eventDTOs []dto.EventResponseDTO

	for _, event := range events {
		location := dto.Location{
			Address1:  event.Address1,
			Address2:  event.Address2,
			City:      event.City,
			State:     event.State,
			Pincode:   event.Pincode,
			Latitude:  event.Latitude,
			Longitude: event.Longitude,
		}

		eventDTO := dto.EventResponseDTO{
			ID:          event.ID,
			EventTime:   event.EventTime,
			Type:        event.Type,
			Description: event.Description,
			Attendees:   event.Attendees,
			ExpiresAt:   event.ExpiresAt,
			Location:    location,
		}
		eventDTOs = append(eventDTOs, eventDTO)
	}
	return eventDTOs, nil
}

func (e *EventServiceImpl) SearchEvents(userProfile model.UserProfile, page model.Pagination, filters dto.EventFilterDTO) ([]dto.EventResponseDTO, error) {

	events, err := e.EventIndexer.SearchEvents(userProfile, page, filters)

	if err != nil {
		return nil, err
	}

	var eventDTOs []dto.EventResponseDTO

	for _, event := range events {
		location := dto.Location{
			Address1:  event.Address1,
			Address2:  event.Address2,
			City:      event.City,
			State:     event.State,
			Pincode:   event.Pincode,
			Latitude:  event.Location[0],
			Longitude: event.Location[1],
		}

		eventDTO := dto.EventResponseDTO{
			ID:          event.ID,
			EventTime:   event.EventTime,
			Type:        event.Type,
			Description: event.Description,
			Attendees:   event.Attendees,
			ExpiresAt:   event.ExpiresAt,
			Location:    location,
		}
		eventDTOs = append(eventDTOs, eventDTO)
	}
	return eventDTOs, nil
}

func (e *EventServiceImpl) createEventElasticDTO(event model.Event) (elasticsearchPkg.Event, error) {
	eventDTO := elasticsearchPkg.Event{
		ID:          event.ID,
		UserId:      event.UserId,
		EventTime:   event.EventTime,
		Description: event.Description,
		Attendees:   event.Attendees,
		Type:        event.Type,
		Address1:    event.Address1,
		Address2:    event.Address2,
		City:        event.City,
		State:       event.State,
		Pincode:     event.Pincode,
		Location:    []float32{event.Latitude, event.Longitude},
		ExpiresAt:   event.ExpiresAt,
		CreatedAt:   &event.CreatedAt,
		UpdatedAt:   &event.UpdatedAt,
		DeletedAt:   convertSqlNullTimeToTime(sql.NullTime(event.DeletedAt)),
	}
	return eventDTO, nil
}

func convertSqlNullTimeToTime(time sql.NullTime) *time.Time {
	if time.Valid {
		return &time.Time
	} else {
		return nil
	}
}

func (e *EventServiceImpl) UpdateUserEvent(createEventDTO dto.CreateEventDTO, userProfile model.UserProfile) (dto.CreateEventDTO, error) {
	event := model.Event{
		Model:       gorm.Model{ID: createEventDTO.ID},
		UserId:      userProfile.UserId,
		EventTime:   createEventDTO.EventTime,
		Type:        createEventDTO.Type,
		Description: createEventDTO.Description,
		Attendees:   createEventDTO.Attendees,
		ExpiresAt:   createEventDTO.ExpiresAt,
		Address1:    createEventDTO.Location.Address1,
		Address2:    createEventDTO.Location.Address2,
		City:        createEventDTO.Location.City,
		Country:     createEventDTO.Location.Country,
		State:       createEventDTO.Location.State,
		Pincode:     createEventDTO.Location.Pincode,
		Latitude:    createEventDTO.Location.Latitude,
		Longitude:   createEventDTO.Location.Longitude,
	}
	event, err := e.EventRepository.UpdateEvent(event)
	if err != nil {
		return dto.CreateEventDTO{}, err
	}
	createEventDTO.ID = event.ID

	//index event to elasticSearch
	elasticDTO, _ := e.createEventElasticDTO(event)
	bytes, err := json.Marshal(elasticDTO)
	if err != nil {
		zapLogger.Logger.Error("Error while marshalling event to elasticDTO", zap.Error(err))
		return dto.CreateEventDTO{}, err
	}
	zapLogger.Logger.Info("Event marshalled to elasticDTO", zap.String("event", string(bytes)))

	err = e.EventIndexer.UpdateUserEvent(elasticDTO, bytes)
	if err != nil {
		zapLogger.Logger.Error("Error while updating user event in ES", zap.Error(err))
		return dto.CreateEventDTO{}, err
	}

	return createEventDTO, nil
}

func (e *EventServiceImpl) DeleteUserEvent(userID, eventID int) error {
	err := e.EventRepository.DeleteEvent(userID, eventID)
	if err != nil {
		zapLogger.Logger.Error("error deleting user event from database", zap.Error(err))
		return err
	}

	err = e.EventIndexer.DeleteUserEvent(eventID)
	if err != nil {
		zapLogger.Logger.Error("error deleting user event from elasticsearch", zap.Error(err))
		return err
	}

	return nil
}

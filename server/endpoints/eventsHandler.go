package endpoints

import (
	"net/http"
	"strconv"

	utils "github.com/SuperMatch/utilities"

	"github.com/SuperMatch/model/dto"
	Service "github.com/SuperMatch/service"
	"github.com/SuperMatch/zapLogger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreateEventIndexHandler godoc
//
//	@Security		ApiKeyAuth
//	@Summary		CreateEventIndexHandler
//	@Description	Create an event index
//	@Tags			Events
//	@Accept			json
//	@Produce		json
//	@Success		200				{string}	string	"ok"
//	@Failure		400				{string}	string	"bad request"
//	@Failure		500				{string}	string	"internal server error"
//	@Router			/event/index	[POST]
func CreateEventIndexHandler(c *gin.Context) {

	eventService := Service.NewEventServiceImpl()
	err := eventService.CreateEventIndex()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error in creating index in elasticSearch",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{"message": "index created successfully"})

}

// CreateUserEventsHandler godoc
//
//	@Security		ApiKeyAuth
//	@Summary		CreateUserEventsHandler
//	@Description	Create an event for a user
//	@Tags			Events
//	@Accept			json
//	@Produce		json
//	@Param			user_id		header		int						true	"User ID"
//	@Param			event		body		dto.CreateEventDTO		true	"Event Data"
//	@Success		200			{object}	dto.EventResponseDTO	"event created successfully."
//	@Failure		400			{string}	string					Bad	request
//	@Failure		500			{string}	string					"internal server error"
//	@Router			/user/event	[POST]
func CreateUserEventsHandler(c *gin.Context) {

	id := c.Request.Header.Get("user_id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var eventDTO dto.CreateEventDTO
	if err := c.BindJSON(&eventDTO); err != nil {
		zapLogger.Logger.Debug("failed to parse event data", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	userProfileService := Service.NewUserProfileService()
	userProfile, err := userProfileService.GetUserProfileFromDB(userID)

	if err != nil {
		zapLogger.Logger.Debug("failed to get user profile .", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	eventService := Service.NewEventServiceImpl()
	eventDTO, err = eventService.CreateUserEvent(eventDTO, userProfile)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": eventDTO, "message": "event created successfully."})

}

// GetUserEventsHandler godoc
//
//	@Security		ApiKeyAuth
//	@Summary		GetUserEventsHandler
//	@Description	Get all events for a user
//	@Tags			Events
//	@Accept			json
//	@Produce		json
//	@Param			user_id			header		int	true	"User ID"
//	@Success		200				{object}	dto.EventResponseDTO
//	@Failure		400				{string}	string	Bad	request
//	@Failure		500				{string}	string	"internal server error"
//	@Router			/user/events	[GET]
func GetUserEventsHandler(c *gin.Context) {
	id := c.Request.Header.Get("user_id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userProfileService := Service.NewUserProfileService()
	userProfile, err := userProfileService.GetUserProfileFromDB(userID)

	if err != nil {
		zapLogger.Logger.Debug("failed to user profile .", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	eventService := Service.NewEventServiceImpl()
	events, err := eventService.GetUserEvent(userProfile)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": events, "message": "events fetched successfully."})

}

// SearchEventsHandler godoc
//
//	@Security		ApiKeyAuth
//	@Summary		SearchEventsHandler
//	@Description	Search events
//	@Tags			Events
//	@Accept			json
//	@Produce		json
//	@Param			eventFilter		body		dto.EventFilterDTO	true	"Event Filter"
//	@Param			user_id			header		int					true	"User ID"
//	@Success		200				{object}	dto.EventResponseDTO
//	@Failure		400				{string}	string	Bad	request
//	@Failure		500				{string}	string	"internal server error"
//	@Router			/events/search	[GET]
func SearchEventsHandler(c *gin.Context) {
	id := c.Request.Header.Get("user_id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	page, err := utils.ReadPaginationDataFromRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var eventFilterDTO dto.EventFilterDTO
	if err := c.BindJSON(&eventFilterDTO); err != nil {
		zapLogger.Logger.Debug("failed to parse event filters", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	userProfileService := Service.NewUserProfileService()
	userProfile, err := userProfileService.GetUserProfileFromDB(userID)

	if err != nil {
		zapLogger.Logger.Debug("failed to fetch user profile.", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	utils.ParseSearchEventFilters(&eventFilterDTO)

	eventService := Service.NewEventServiceImpl()
	events, err := eventService.SearchEvents(userProfile, page, eventFilterDTO)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": events, "message": "events fetched successfully."})
}

// UpdateUserEventHandler godoc
//
//	@Security		ApiKeyAuth
//	@Summary		UpdateUserEventHandler
//	@Description	Update event for a user
//	@Tags			Events
//	@Accept			json
//	@Produce		json
//	@Param			event		body		dto.CreateEventDTO	true	"Event Data"
//	@Param			user_id		header		int					true	"User ID"
//	@Success		200			{object}	dto.CreateEventDTO
//	@Failure		400			{string}	string	Bad	request
//	@Failure		500			{string}	string	"internal server error"
//	@Router			/user/event	[PUT]
func UpdateUserEventHandler(c *gin.Context) {
	id := c.Request.Header.Get("user_id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var eventDTO dto.CreateEventDTO
	if err := c.BindJSON(&eventDTO); err != nil {
		zapLogger.Logger.Debug("failed to parse event data", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	userProfileService := Service.NewUserProfileService()
	userProfile, err := userProfileService.GetUserProfileFromDB(userID)

	if err != nil {
		zapLogger.Logger.Debug("failed to fetch user profile.", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	eventService := Service.NewEventServiceImpl()
	eventDTO, err = eventService.UpdateUserEvent(eventDTO, userProfile)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": eventDTO, "message": "event updated successfully."})

}

// DeleteUserEventHandler godoc
//
//	@Security		ApiKeyAuth
//	@Summary		DeleteUserEventHandler
//	@Description	Delete event for a user
//	@Tags			Events
//	@Accept			json
//	@Produce		json
//	@Param			event		body		dto.CreateEventDTO	true	"Event Data"
//	@Param			user_id		header		int					true	"User ID"
//	@Param			event_id	header		int					true	"Event ID"
//	@Success		200			{string}	string				"event deleted successfully."
//	@Failure		400			{string}	string				Bad	request
//	@Failure		500			{string}	string				"internal server error"
//	@Router			/user/event	[DELETE]
func DeleteUserEventHandler(c *gin.Context) {
	id := c.Request.Header.Get("user_id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id = c.Request.Header.Get("event_id")
	eventID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	eventService := Service.NewEventServiceImpl()
	err = eventService.DeleteUserEvent(userID, eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error while deleting user event", "error": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "event deleted successfully."})
}

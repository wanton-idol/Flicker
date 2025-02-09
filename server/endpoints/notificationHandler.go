package endpoints

import (
	"github.com/SuperMatch/model"
	"github.com/SuperMatch/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetDeviceToken(c *gin.Context) {
	ID := c.Request.Header.Get("user_id")
	userID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var deviceToken model.UserDeviceToken
	if err := c.BindJSON(&deviceToken); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error in parsing request."})
		return
	}
	deviceToken.UserID = userID

	notificationService := service.NewNotificationService()
	err = notificationService.InsertDeviceToken(deviceToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error in inserting device token", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "device token inserted successfully"})

}

func SendNotificationToUser(c *gin.Context) {
	ID := c.Request.Header.Get("user_id")
	userID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var message model.NotificationData
	if err := c.BindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error in parsing request."})
		return
	}

	notificationService := service.NewNotificationService()
	err = notificationService.SendNotificationToUser(userID, message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error in sending notification to user's all devices", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "notification sent successfully"})

}

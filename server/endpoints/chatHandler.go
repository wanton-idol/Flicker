package endpoints

import (
	"github.com/SuperMatch/model"
	"github.com/SuperMatch/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// SaveMessage godoc
//
//	@Security		ApiKeyAuth
//	@Summary		SaveMessage
//	@Description	Save Message and Media
//	@Tags			Chat
//	@Accept			mpfd
//	@Produce		json
//	@Param			media			formData	file	true	"Media"
//	@Param			message			formData	string	true	"Message"
//	@Param			sender_id		header		int		true	"Sender ID"
//	@Param			receiver_id		header		int		true	"Receiver ID"
//	@Success		200				{string}	string	"chat saved successfully"
//	@Failure		400				{string}	string	Bad	request
//	@Failure		500				{string}	string	"internal server error"
//	@Router			/chat/message	[POST]
func SaveMessage(c *gin.Context) {
	userID := c.Request.Header.Get("sender_id")
	senderID, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID2 := c.Request.Header.Get("receiver_id")
	receiverID, err := strconv.Atoi(userID2)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	form, _ := c.MultipartForm()
	files := form.File["media"]
	message := form.Value["message"]

	chatService := service.NewChatService()
	err = chatService.SaveMessage(senderID, receiverID, message[0], files)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error in saving chat", "error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "chat saved successfully"})

}

// RetrieveUserChats godoc
//
//	@Security		ApiKeyAuth
//	@Summary		RetrieveUserChats
//	@Description	Retrieve User Chats
//	@Tags			Chat
//	@Accept			json
//	@Produce		json
//	@Param			chat_id				header		int		true	"Chat ID"
//	@Success		200					{string}	string	"chats retrieved successfully"
//	@Failure		400					{string}	string	Bad	request
//	@Failure		500					{string}	string	"internal server error"
//	@Router			/chat/user/chats	[GET]
func RetrieveUserChats(c *gin.Context) {
	chatID := c.Request.Header.Get("chat_id")

	chatService := service.NewChatService()
	chats, err := chatService.RetrieveUserChats(chatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error in retrieving chats", "error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "chats retrieved successfully", "data": chats})
}

// GetUserChatsList godoc
//
//	@Security		ApiKeyAuth
//	@Summary		GetUserChatsList
//	@Description	Get User Chats List
//	@Tags			Chat
//	@Accept			json
//	@Produce		json
//	@Param			user_id			header		int		true	"User ID"
//	@Success		200				{string}	string	"chats retrieved successfully
//	@Failure		400				{string}	string	Bad	request
//	@Failure		500				{string}	string	"internal server error"
//	@Router			/chat/user/list	[GET]
func GetUserChatsList(c *gin.Context) {
	userID := c.Request.Header.Get("user_id")
	ID, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chatService := service.NewChatService()
	chats, err := chatService.GetUserChatsList(ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error in getting chat list", "error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "chats retrieved successfully", "data": chats})
}

// UpdateMessagesStatus godoc
//
//	@Security		ApiKeyAuth
//	@Summary		UpdateMessagesStatus
//	@Description	Update Messages Status
//	@Tags			Chat
//	@Accept			json
//	@Produce		json
//	@Param			messageIDs				body		model.MessagesIDs	true	"Message IDs"
//	@Success		200						{string}	string				"messages status updated successfully"
//	@Failure		400						{string}	string				Bad	request
//	@Failure		500						{string}	string				"internal server error"
//	@Router			/chat/messages/status	[PUT]
func UpdateMessagesStatus(c *gin.Context) {
	var messageIDs model.MessagesIDs

	if err := c.BindJSON(&messageIDs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error in parsing request."})
		return
	}

	chatService := service.NewChatService()
	err := chatService.UpdateMessagesStatus(messageIDs.MessageIDS)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error in updating messages status", "error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "messages status updated successfully"})
}

// GetLastMessages godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Get Last Messages
//	@Description	Get Last Messages of Chats
//	@Tags			Chat
//	@Accept			json
//	@Produce		json
//	@Param			chatIDs				body		model.ChatIDs	true	"Chat IDs"
//	@Success		200					{string}	string			"last messages retrieved successfully"
//	@Failure		400					{string}	string			Bad	request
//	@Failure		500					{string}	string			"internal server error"
//	@Router			/chat/last/messages	[GET]
func GetLastMessages(c *gin.Context) {
	var chatIds model.ChatIDs

	if err := c.BindJSON(&chatIds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error in parsing request."})
		return
	}

	chatService := service.NewChatService()
	chats, err := chatService.RetrieveLastMessages(chatIds.ChatIDS)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error in retrieving last messages", "error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "last messages retrieved successfully", "data": chats})
}

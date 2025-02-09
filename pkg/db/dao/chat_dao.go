package dao

import (
	"errors"
	"github.com/SuperMatch/model"
	"github.com/SuperMatch/pkg/db"
	"github.com/SuperMatch/zapLogger"
	"gorm.io/gorm"
)

type ChatDao interface {
	Insert(chatDetails model.ChatDetails) error
	RetrieveUserChats(chatID string) ([]model.ChatDetails, error)
	GetUserChatsList(userID int) ([]model.ChatDetails, error)
	UpdateMessageReadStatus(messageIDs []int) error
	RetrieveLastMessages(chatIDs []string) ([]model.ChatDetails, error)
}

type ChatDaoImpl struct {
	Connection gorm.DB
}

func NewChatDaoImpl() *ChatDaoImpl {
	return &ChatDaoImpl{Connection: *db.GlobalOrm}
}

func (c *ChatDaoImpl) Insert(chatDetails model.ChatDetails) error {
	err := c.Connection.Table("user_chats").Create(&chatDetails)
	if err.Error != nil {
		zapLogger.Logger.Error("error inserting chat details in user_chats table")
		return err.Error
	}

	return nil
}

func (c *ChatDaoImpl) RetrieveUserChats(chatID string) ([]model.ChatDetails, error) {
	var chats []model.ChatDetails
	err := c.Connection.Table("user_chats").Where("chat_id = ?", chatID).Order("created_at DESC").Find(&chats)
	if errors.Is(err.Error, gorm.ErrRecordNotFound) {
		zapLogger.Logger.Error("no chats found for the user")
		return nil, err.Error
	} else if err.Error != nil {
		zapLogger.Logger.Error("error in retrieving chats")
		return nil, err.Error
	}

	return chats, nil
}

func (c *ChatDaoImpl) GetUserChatsList(userID int) ([]model.ChatDetails, error) {
	var chats []model.ChatDetails
	err := c.Connection.Debug().Table("user_chats").Where("sender_id = ? OR receiver_id = ?", userID, userID).Group("chat_id").Order("created_at DESC").Find(&chats)
	if errors.Is(err.Error, gorm.ErrRecordNotFound) {
		zapLogger.Logger.Error("no chats found for the user")
		return nil, err.Error
	} else if err.Error != nil {
		zapLogger.Logger.Error("error in retrieving chats")
		return nil, err.Error
	}

	return chats, nil
}

func (c *ChatDaoImpl) UpdateMessageReadStatus(messageIDs []int) error {
	err := c.Connection.Debug().Table("user_chats").Where("ID in ?", messageIDs).Update("is_read", true)
	if err.Error != nil {
		zapLogger.Logger.Error("error in updating message read status")
		return err.Error
	}

	return nil
}

func (c *ChatDaoImpl) RetrieveLastMessages(chatIDs []string) ([]model.ChatDetails, error) {
	var chats []model.ChatDetails

	query := "select * from user_chats uc where chat_id in ? and created_at  = (select max(created_at) from user_chats uc2 where uc.chat_id = uc2.chat_id);"

	err := c.Connection.Debug().Raw(query, chatIDs).Find(&chats)
	if errors.Is(err.Error, gorm.ErrRecordNotFound) {
		zapLogger.Logger.Error("no chats found for the user")
		return nil, err.Error
	} else if err.Error != nil {
		zapLogger.Logger.Error("error in retrieving chats")
		return nil, err.Error
	}

	return chats, nil
}

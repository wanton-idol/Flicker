package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/SuperMatch/model"
	"github.com/SuperMatch/pkg/db/dao"
	"github.com/SuperMatch/zapLogger"
	"go.uber.org/zap"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"
)

type ChatServiceInterface interface {
	SaveMessage(senderID, receiverID int, message string, media []*multipart.FileHeader) error
	RetrieveUserChats(chatID string) ([]model.ChatDetails, error)
	GetUserChatsList(userID int) ([]model.ChatValues, error)
	UpdateMessagesStatus(messageIDs []int) error
	RetrieveLastMessages(chatIDs []string) ([]model.ChatDetails, error)
}

type ChatService struct {
	s3Service          S3ServiceInterface
	userProfileService UserProfileInterface
	chatDao            dao.ChatDao
	userMatchDao       dao.UserMatchDao
	userProfileDao     dao.UserProfileRepository
	userMedia          dao.UserMediaRepository
}

func NewChatService() *ChatService {
	return &ChatService{
		s3Service:          NewS3Service(),
		userProfileService: NewUserProfileService(),
		chatDao:            dao.NewChatDaoImpl(),
		userMatchDao:       dao.NewUserMatchDaoImpl(),
		userProfileDao:     dao.NewUserProfileRepository(),
		userMedia:          dao.NewUserMediaRepository(),
	}
}

func (c *ChatService) SaveMessage(senderID, receiverID int, message string, media []*multipart.FileHeader) error {
	mediaUrl := ""
	if media != nil {
		fileExt := filepath.Ext(media[0].Filename)
		if !c.userProfileService.CheckAllowedFileType(fileExt) && !c.userProfileService.CheckAllowedAudioFileType(fileExt) {
			zapLogger.Logger.Error("file extension is not supported")
			return errors.New("file extension is not supported")
		}

		filename := fmt.Sprintf("%v", time.Now().UTC().Format("2006-01-02T15:04:05.00000")) + fileExt
		tempFile, _ := media[0].Open()
		S3filepath := fmt.Sprintf("%d/chatMedia/%s", senderID, filename)

		result, err := c.s3Service.UploadFileToS3(user_profile_S3_bucket, S3filepath, tempFile, filename)
		if err != nil {
			zapLogger.Logger.Error("error in uploading file to S3")
			return err
		}
		mediaUrl = result
	}

	userMatch, err := c.userMatchDao.FindByUserIdMatchId(context.Background(), senderID, receiverID)
	if err != nil {
		zapLogger.Logger.Error("error in finding user match")
		return err
	}

	chatDetails := model.ChatDetails{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Message:    message,
		MediaURL:   mediaUrl,
		ChatID:     *userMatch.ChatID,
	}

	err = c.chatDao.Insert(chatDetails)
	if err != nil {
		zapLogger.Logger.Error("error in inserting chat details")
		return err
	}

	return nil
}

func (c *ChatService) RetrieveUserChats(chatID string) ([]model.ChatDetails, error) {
	chats, err := c.chatDao.RetrieveUserChats(chatID)
	if err != nil {
		zapLogger.Logger.Error("error in retrieving chats")
		return nil, err
	}

	for idx := range chats {
		if chats[idx].MediaURL != "" {
			key := strings.ReplaceAll(strings.TrimPrefix(chats[idx].MediaURL, S3_BUCKET_PATH), "%3A", ":")
			signedURL, err := c.s3Service.SignS3FilesUrl(user_profile_S3_bucket, key)
			if err != nil {
				zapLogger.Logger.Error("error in getting signed url ", zap.Error(err))
				return nil, err
			}
			chats[idx].MediaURL = signedURL
		}
	}
	return chats, nil
}

func (c *ChatService) GetUserChatsList(userID int) ([]model.ChatValues, error) {
	userChats, err := c.chatDao.GetUserChatsList(userID)
	if err != nil {
		zapLogger.Logger.Error("error in retrieving user chats list")
		return nil, err
	}

	var chatList []model.ChatValues

	for _, chats := range userChats {
		var id int
		if chats.ReceiverID == userID {
			id = chats.SenderID
		} else {
			id = chats.ReceiverID
		}

		userProfile, err := c.userProfileDao.FindByUserId(context.Background(), id)
		if err != nil {
			zapLogger.Logger.Error("error while getting user profile: ", zap.Error(err))
			return nil, err
		}

		imageUrl, err := c.userMedia.FindByUserIDOrderID(id, 1)
		if err != nil {
			zapLogger.Logger.Error("error in getting image url")
			return nil, err
		}

		if imageUrl.URL != "" {
			key := strings.ReplaceAll(strings.TrimPrefix(imageUrl.URL, S3_BUCKET_PATH), "%3A", ":")
			signedURL, err := c.s3Service.SignS3FilesUrl(user_profile_S3_bucket, key)
			if err != nil {
				zapLogger.Logger.Error("error in getting signed url ", zap.Error(err))
				return nil, err
			}
			imageUrl.URL = signedURL
		}

		chatValue := model.ChatValues{
			ChatID:      chats.ChatID,
			UserProfile: userProfile,
			Image:       imageUrl.URL,
		}

		chatList = append(chatList, chatValue)

	}

	return chatList, nil
}

func (c *ChatService) UpdateMessagesStatus(messageIDs []int) error {
	err := c.chatDao.UpdateMessageReadStatus(messageIDs)
	if err != nil {
		zapLogger.Logger.Error("error in updating messages status")
		return err
	}

	return nil
}

func (c *ChatService) RetrieveLastMessages(chatIDs []string) ([]model.ChatDetails, error) {
	lastMessages, err := c.chatDao.RetrieveLastMessages(chatIDs)
	if err != nil {
		zapLogger.Logger.Error("error in retrieving last messages")
		return nil, err
	}

	for idx := range lastMessages {
		if lastMessages[idx].MediaURL != "" {
			key := strings.ReplaceAll(strings.TrimPrefix(lastMessages[idx].MediaURL, S3_BUCKET_PATH), "%3A", ":")
			signedURL, err := c.s3Service.SignS3FilesUrl(user_profile_S3_bucket, key)
			if err != nil {
				zapLogger.Logger.Error("error in getting signed url ", zap.Error(err))
				return nil, err
			}
			lastMessages[idx].MediaURL = signedURL
		}
	}

	return lastMessages, nil
}

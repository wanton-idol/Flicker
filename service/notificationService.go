package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SuperMatch/config"
	"github.com/SuperMatch/model"
	"github.com/SuperMatch/pkg/db/dao"
	"github.com/SuperMatch/zapLogger"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"go.uber.org/zap"
)

type NotificationServiceInterface interface {
	InsertDeviceToken(deviceToken model.UserDeviceToken) error
	SendNotificationToUser(userID int, message model.NotificationData) error
}

type NotificationService struct {
	s3Service          S3ServiceInterface
	userDeviceTokenDao dao.UserDeviceTokenDao
	userProfileService UserProfileInterface
}

func NewNotificationService() *NotificationService {
	return &NotificationService{
		s3Service:          NewS3Service(),
		userDeviceTokenDao: dao.NewUserDeviceTokenDaoImpl(),
		userProfileService: NewUserProfileService(),
	}
}

func (n *NotificationService) InsertDeviceToken(deviceToken model.UserDeviceToken) error {
	deviceTokens, err := n.userDeviceTokenDao.GetDeviceTokensByUserId(deviceToken.UserID)
	if err != nil {
		zapLogger.Logger.Error("error getting user device tokens from DB", zap.Error(err))
		return err
	}
	if err == nil && len(deviceTokens) > 0 {
		for _, token := range deviceTokens {
			if token.DeviceToken == deviceToken.DeviceToken {
				zapLogger.Logger.Info("device token already exists for user")
				return errors.New("device token already exists for user")
			}
		}
	}

	endpointARN, err := createPlatformEndpoint(deviceToken.DeviceToken)
	if err != nil {
		zapLogger.Logger.Error("error creating platform endpoint", zap.Error(err))
		return err
	}
	deviceToken.EndpointARN = endpointARN

	err = n.userDeviceTokenDao.Insert(deviceToken)
	if err != nil {
		zapLogger.Logger.Error("error inserting user device token in DB", zap.Error(err))
		return err
	}

	return nil
}

func (n *NotificationService) SendNotificationToUser(userID int, message model.NotificationData) error {
	deviceTokens, err := n.userDeviceTokenDao.GetDeviceTokensByUserId(userID)
	if err != nil {
		zapLogger.Logger.Error("error getting user device tokens from DB", zap.Error(err))
		return err
	}

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(config.AppConfig.AWSConfig.Region),
		Credentials: credentials.NewStaticCredentials(config.AppConfig.AWSConfig.AccessKeyID, config.AppConfig.AWSConfig.AccessKeySecret, ""),
	})
	if err != nil {
		zapLogger.Logger.Error("Failed to create session:", zap.Error(err))
		return err
	}

	svc := sns.New(sess)

	payload := convertMessageToPayload(message)

	for _, deviceToken := range deviceTokens {
		_, err := svc.Publish(&sns.PublishInput{
			Message:          aws.String(payload),
			MessageStructure: aws.String("json"),
			TargetArn:        aws.String(deviceToken.EndpointARN),
		})
		if err != nil {
			zapLogger.Logger.Error("error sending notification to user", zap.Error(err))
			return err
		}
	}

	return nil
}

func createPlatformEndpoint(deviceToken string) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(config.AppConfig.AWSConfig.Region),
		Credentials: credentials.NewStaticCredentials(config.AppConfig.AWSConfig.AccessKeyID, config.AppConfig.AWSConfig.AccessKeySecret, ""),
	})
	if err != nil {
		zapLogger.Logger.Error("Failed to create session:", zap.Error(err))
		return "", err
	}

	svc := sns.New(sess)

	resp, err := svc.CreatePlatformEndpoint(&sns.CreatePlatformEndpointInput{
		PlatformApplicationArn: aws.String(config.AppConfig.AWSConfig.PlatformApplicationArn),
		Token:                  aws.String(deviceToken),
	})
	if err != nil {
		zapLogger.Logger.Error("Failed to create platform endpoint:", zap.Error(err))
		return "", err
	}

	return *resp.EndpointArn, nil
}

func convertMessageToPayload(message model.NotificationData) string {
	notificationObject := model.Notification{
		Notification: model.NotificationData{
			Title: message.Title,
			Body:  message.Body,
		},
	}
	notificationBytes, _ := json.Marshal(notificationObject)

	messageObject := model.GCMNotification{
		GCM: string(notificationBytes),
	}

	messageBytes, _ := json.Marshal(messageObject)
	payload := string(messageBytes)
	fmt.Println(payload)

	return payload
}

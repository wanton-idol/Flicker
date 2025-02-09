package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SuperMatch/model"
	elasticsearchPkg "github.com/SuperMatch/model/elasticSearch"
	pkg "github.com/SuperMatch/pkg/elasticSeach"
	"github.com/SuperMatch/zapLogger"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"log"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"
)

const (
	user_stories_S3_bucket      = "user-stories-supermatch"
	user_Stories_S3_Bucket_Path = "https://" + user_stories_S3_bucket + ".s3.ap-south-1.amazonaws.com"
)

type UserStoriesInterface interface {
	CreateStoriesIndex() error
	UploadFileToS3(userId string, file *multipart.FileHeader) (string, string, error)
	IndexUserStories(userStories elasticsearchPkg.UserStories) error
	GetUserStoriesByProfileID(userProfileID int) ([]elasticsearchPkg.UserStories, error)
	GetUserStoriesByLocation(location model.UserLocation) ([]elasticsearchPkg.UserStories, error)
}

type UserStoriesService struct {
	esIndex            pkg.UserStoriesIndexer
	userProfileService UserProfileInterface
	s3Service          S3ServiceInterface
}

func NewUserStoriesService() *UserStoriesService {
	return &UserStoriesService{
		esIndex:            pkg.NewUserStoriesIndexerImpl(),
		userProfileService: NewUserProfileService(),
		s3Service:          NewS3Service(),
	}
}

func (u *UserStoriesService) CreateStoriesIndex() error {
	err := u.esIndex.CreateIndex()
	if err != nil {
		log.Println("error in creating index in elasticSearch")
		return err
	}

	return nil
}

func (u *UserStoriesService) UploadFileToS3(userId string, file *multipart.FileHeader) (string, string, error) {
	var mediaType string
	fileExt := filepath.Ext(file.Filename)
	if u.userProfileService.CheckAllowedFileType(fileExt) {
		mediaType = "video"
	} else if u.userProfileService.CheckAllowedAudioFileType(fileExt) {
		mediaType = "audio"
	} else {
		log.Println("file extension is not supported.")
		return "", "", errors.New("file extension is not supported")
	}

	filename := fmt.Sprintf("%v", time.Now().UTC().Format("2006-01-02T15:04:05.00000")) + fileExt
	tempFile, _ := file.Open()
	S3filepath := userId + "/user_stories/" + filename

	result, err := u.s3Service.UploadFileToS3(user_stories_S3_bucket, S3filepath, tempFile, filename)
	if err != nil {
		log.Println("error in uploading file to S3")
		return result, mediaType, err
	}

	return result, mediaType, nil
}

func (u *UserStoriesService) IndexUserStories(userStories elasticsearchPkg.UserStories) error {
	userStories.ID = uuid.New().String()
	userStories.CreatedAt = time.Now()
	userStories.ExpiresAt = time.Now().Add(time.Hour * 24 * 1)

	fmt.Println(userStories)

	jsonString, err := json.Marshal(userStories)
	log.Println(string(jsonString))
	if err != nil {
		return err
	}

	err = u.esIndex.IndexUserStories(userStories, jsonString)
	if err != nil {
		log.Println("error in indexing user stories")
		return err
	}

	return nil
}

func (u *UserStoriesService) GetUserStoriesByProfileID(userProfileID int) ([]elasticsearchPkg.UserStories, error) {
	userStories, err := u.esIndex.GetUserStoriesByProfileID(userProfileID)
	if err != nil {
		zapLogger.Logger.Error("error in getting user stories from elastic search")
		return nil, err
	}

	for idx := range userStories {
		if userStories[idx].MediaURL != "" {
			key := strings.ReplaceAll(strings.TrimPrefix(userStories[idx].MediaURL, user_Stories_S3_Bucket_Path), "%3A", ":")
			signedURL, err := u.s3Service.SignS3FilesUrl(user_stories_S3_bucket, key)
			if err != nil {
				zapLogger.Logger.Error("error in getting signed url ", zap.Error(err))
				return nil, err
			}
			userStories[idx].MediaURL = signedURL
		}
	}
	return userStories, nil
}

func (u *UserStoriesService) GetUserStoriesByLocation(location model.UserLocation) ([]elasticsearchPkg.UserStories, error) {
	userStories, err := u.esIndex.GetUserStoriesByLocation(location)
	if err != nil {
		zapLogger.Logger.Error("error in getting user stories from elastic search based on location")
		return userStories, err
	}

	for idx := range userStories {
		if userStories[idx].MediaURL != "" {
			key := strings.ReplaceAll(strings.TrimPrefix(userStories[idx].MediaURL, user_Stories_S3_Bucket_Path), "%3A", ":")
			signedURL, err := u.s3Service.SignS3FilesUrl(user_stories_S3_bucket, key)
			if err != nil {
				zapLogger.Logger.Error("error in getting signed url ", zap.Error(err))
				return userStories, err
			}
			userStories[idx].MediaURL = signedURL
		}
	}
	return userStories, nil
}

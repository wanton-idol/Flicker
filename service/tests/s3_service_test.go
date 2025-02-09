package tests

import (
	"github.com/SuperMatch/service/mocks"
	"github.com/golang/mock/gomock"
	"strings"
	"testing"
)

const (
	userProfileS3Bucket = "user-profile-supermatch"
	S3BucketPath        = "https://" + userProfileS3Bucket + ".s3.ap-south-1.amazonaws.com"
)

func TestUploadFileToS3(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctrl.Finish()
	filename := "2023-04-24T03:38:23-images3.webp"
	userID := "1"
	path := userID + "/profile/" + filename
	url := "https://user-profile-supermatch.s3.ap-south-1.amazonaws.com/1/profile/2023-04-24T03%3A38%3A23-images3.webp"

	mockS3 := mocks.NewMockS3ServiceInterface(ctrl)
	mockS3.EXPECT().UploadFileToS3(gomock.Eq(userProfileS3Bucket), gomock.Eq(path), gomock.Any(), gomock.Eq(filename)).
		Return(url, nil)
}

func TestGetFilesInFolder(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctrl.Finish()

	signedURLs := make([]string, 0)
	filename := "2023-04-24T03:38:23-images3.webp"
	userID := "1"
	folder := userID + "/profile/" + filename
	mockS3 := mocks.NewMockS3ServiceInterface(ctrl)
	mockS3.EXPECT().GetFilesInFolder(gomock.Eq(userProfileS3Bucket), gomock.Eq(folder)).
		Return(signedURLs, nil)
}

func TestSignS3FilesUrl(t *testing.T) {
	url := "https://user-profile-supermatch.s3.ap-south-1.amazonaws.com/1/profile/2023-04-24T03%3A38%3A23-images3.webp"
	key := strings.ReplaceAll(strings.TrimPrefix(url, S3BucketPath), "%3A", ":")
	signedURL := ""
	ctrl := gomock.NewController(t)
	ctrl.Finish()

	mockS3 := mocks.NewMockS3ServiceInterface(ctrl)
	mockS3.EXPECT().SignS3FilesUrl(gomock.Eq(userProfileS3Bucket), gomock.Eq(key)).
		Return(signedURL, nil)
}

func TestDeleteFile(t *testing.T) {
	url := "https://user-profile-supermatch.s3.ap-south-1.amazonaws.com/1/profile/2023-04-24T03%3A38%3A23-images3.webp"
	key := strings.ReplaceAll(strings.TrimPrefix(url, S3BucketPath), "%3A", ":")

	ctrl := gomock.NewController(t)
	ctrl.Finish()

	mockS3 := mocks.NewMockS3ServiceInterface(ctrl)
	mockS3.EXPECT().DeleteFile(gomock.Eq(userProfileS3Bucket), gomock.Eq(key)).
		Return(nil)
}

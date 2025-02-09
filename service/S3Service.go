package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/SuperMatch/config"
	"github.com/SuperMatch/zapLogger"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"mime/multipart"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ses"
	"go.uber.org/zap"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3ServiceInterface interface {
	UploadFileToS3(bucket string, path string, file multipart.File, fileName string) (string, error)
	GetFilesInFolder(bucket string, folder string) ([]string, error)
	SignS3FilesUrl(bucket string, url string) (string, error)
	DeleteFile(bucket string, key string) error
	SendEmailInput(receiverEmail, htmlBody, title string) *ses.SendEmailInput
	SendEmail(input *ses.SendEmailInput) error
}

type S3Service struct {
}

func NewS3Service() *S3Service {
	return &S3Service{}
}

const (
	sender = "whonishchal@gmail.com"
)

func (s *S3Service) UploadFileToS3(bucket string, path string, file multipart.File, fileName string) (string, error) {

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(config.AppConfig.AWSConfig.Region), // Replace with your desired region
		Credentials: credentials.NewStaticCredentials(config.AppConfig.AWSConfig.AccessKeyID, config.AppConfig.AWSConfig.AccessKeySecret, ""),
	})
	if err != nil {
		zapLogger.Logger.Error("Failed to create session:", zap.Error(err))
		return "", errors.New("failed to connect s3")
	}

	uploader := s3manager.NewUploader(sess)
	result, err := uploader.UploadWithContext(context.Background(), &s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(path),
		Body:   file,
	})
	if err != nil {
		return "", err
	}
	zapLogger.Logger.Debug(fmt.Sprintf("file uploaded to=%s ", result.Location))
	return result.Location, nil
}

func (s *S3Service) GetFilesInFolder(bucket string, folder string) ([]string, error) {

	signedURLs := make([]string, 0)
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(config.AppConfig.AWSConfig.Region), // Replace with your desired region
		Credentials: credentials.NewStaticCredentials(config.AppConfig.AWSConfig.AccessKeyID, config.AppConfig.AWSConfig.AccessKeySecret, ""),
	})
	if err != nil {
		zapLogger.Logger.Error("Failed to create session:", zap.Error(err))
		return nil, errors.New("failed to connect s3")
	}
	svc := s3.New(sess)

	params := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(folder),
	}

	resp, err := svc.ListObjectsV2(params)

	for _, key := range resp.Contents {
		signedURL, err := s.SignS3FilesUrl(bucket, *key.Key)

		if err != nil {
			zapLogger.Logger.Error("error in signing url for ", zap.Any("key", *key.Key), zap.Error(err))
			break
		}
		signedURLs = append(signedURLs, signedURL)
	}

	if err != nil {
		return nil, err
	}
	return signedURLs, nil
}

func (s *S3Service) SignS3FilesUrl(bucket string, url string) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(config.AppConfig.AWSConfig.Region), // Replace with your desired region
		Credentials: credentials.NewStaticCredentials(config.AppConfig.AWSConfig.AccessKeyID, config.AppConfig.AWSConfig.AccessKeySecret, ""),
	})
	if err != nil {
		zapLogger.Logger.Error("Failed to create session:", zap.Error(err))
		return "", errors.New("failed to connect s3")
	}

	if err != nil {
		return "", err
	}
	svc := s3.New(sess)
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(url),
	})

	urlStr, err := req.Presign(15 * time.Minute)

	if err != nil {
		return "", err
	}
	return urlStr, nil
}

func (s *S3Service) DeleteFile(bucket string, key string) error {

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(config.AppConfig.AWSConfig.Region), // Replace with your desired region
		Credentials: credentials.NewStaticCredentials(config.AppConfig.AWSConfig.AccessKeyID, config.AppConfig.AWSConfig.AccessKeySecret, ""),
	})

	if err != nil {
		zapLogger.Logger.Error("Failed to create session:", zap.Error(err))
		return errors.New("failed to connect s3")
	}

	svc := s3.New(sess)
	zapLogger.Logger.Debug(fmt.Sprintf("deleting file from s3 bucket=%s for key=%s", bucket, key))

	params := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	_, err = svc.DeleteObject(params)
	if err != nil {
		zapLogger.Logger.Error("error in deleting file from s3", zap.Error(err))
		return err
	}
	return nil
}

func (s *S3Service) SendEmailInput(receiverEmail, htmlBody, title string) *ses.SendEmailInput {
	return &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(receiverEmail),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("utf-8"),
					Data:    aws.String(htmlBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("utf-8"),
				Data:    aws.String(title),
			},
		},
		Source: aws.String(sender),
	}

}

func (s *S3Service) SendEmail(input *ses.SendEmailInput) error {

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(config.AppConfig.AWSConfig.Region), // Replace with your desired region
		Credentials: credentials.NewStaticCredentials(config.AppConfig.AWSConfig.AccessKeyID, config.AppConfig.AWSConfig.AccessKeySecret, ""),
	})

	if err != nil {
		zapLogger.Logger.Error("Failed to create session:", zap.Error(err))
		return errors.New("failed to connect email service")
	}

	svc := ses.New(sess)

	_, err = svc.SendEmail(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				zapLogger.Logger.Error(fmt.Sprintf(ses.ErrCodeMessageRejected+": %s", aerr.Error()))
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				zapLogger.Logger.Error(fmt.Sprintf(ses.ErrCodeMailFromDomainNotVerifiedException+": %s", aerr.Error()))
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				zapLogger.Logger.Error(fmt.Sprintf(ses.ErrCodeConfigurationSetDoesNotExistException+": %s", aerr.Error()))
			default:
				zapLogger.Logger.Error(fmt.Sprintf(aerr.Error()))
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			zapLogger.Logger.Error(fmt.Sprintf(err.Error()))
		}
	}
	zapLogger.Logger.Info("Email sent successfully.")
	return nil
}

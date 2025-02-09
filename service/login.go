package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/SuperMatch/config"
	"github.com/SuperMatch/model/dto"
	"github.com/SuperMatch/utilities"
	"github.com/SuperMatch/zapLogger"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/twilio/twilio-go"
	"go.uber.org/zap"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"log"
	"math/rand"
	"strings"
	"time"

	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
	"gorm.io/gorm"

	"github.com/SuperMatch/model"
	"github.com/SuperMatch/pkg/db"
	"github.com/SuperMatch/pkg/db/dao"
	"google.golang.org/api/oauth2/v2"
)

type LoginInterface interface {
	GoogleLogin(idToken string) (*dto.TokenInfo, error)
	UserSignIN(user model.User) (string, time.Time, error)
	UserSignUP(userModel model.User) (int, string, time.Time, error)
	GenerateAndSaveToken(user model.User) (string, time.Time, error)
	SendOTPService(phoneNumber string) (string, error)
	VerifyOTPService(phoneNumber string, OTP string) (bool, error)
	GetUserDetailsFromGoogle(googleResponse *dto.TokenInfo) model.User
	CheckUserExistOrNot(userMobile string) (string, time.Time, error)
	SendVerificationEmail(emailId string) error
	VerifyEmailService(verificationCode string) error
	DeleteUser(userID int, email string) error
}

type LoginService struct{}

func NewLoginService() *LoginService {
	return &LoginService{}
}

func twilioClient() *twilio.RestClient {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: config.ConfigValue.TwilioConfig.AccountSID,
		Password: config.ConfigValue.TwilioConfig.AuthToken,
	})

	return client
}

func (l *LoginService) GoogleLogin(idToken string) (*dto.TokenInfo, error) {
	authService, err := oauth2.NewService(context.Background(), option.WithoutAuthentication())
	if err != nil {
		zapLogger.Logger.Error("oauth2.NewService failed", zap.Error(err))
		return nil, err
	}

	tokenInfoCall := authService.Tokeninfo()
	tokenInfoCall.IdToken(idToken)
	ctx, cancelFunc := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancelFunc()
	tokenInfoCall.Context(ctx)
	_, err = tokenInfoCall.Do()
	if err != nil {
		var e *googleapi.Error
		_ = errors.As(err, &e)
		return nil, e
	}

	token, _, err := new(jwt.Parser).ParseUnverified(idToken, &dto.TokenInfo{})
	if tokenInfo, ok := token.Claims.(*dto.TokenInfo); ok {
		return tokenInfo, nil
	} else {
		zapLogger.Logger.Error("parse token.payload failed", zap.Error(err))
		return nil, err
	}
}

func (l *LoginService) UserSignIN(user model.User) (string, time.Time, error) {
	token, expiresAt, err := l.GenerateAndSaveToken(user)
	if err != nil {
		zapLogger.Logger.Error("error in generating and saving auth token for user", zap.Error(err))
		return token, expiresAt, err
	}

	return token, expiresAt, nil
}

func (l *LoginService) UserSignUP(userModel model.User) (int, string, time.Time, error) {
	userDao := &dao.UserDao{
		Connection: *db.GlobalOrm,
	}

	user, err := userDao.Insert(userModel)
	if err != nil {
		zapLogger.Logger.Error("error in creating user", zap.Error(err))
		return 0, "", time.Time{}, err
	}

	token, expiresAt, err := l.GenerateAndSaveToken(user)
	if err != nil {
		zapLogger.Logger.Error("error in generating and saving auth token for user", zap.Error(err))
		return 0, token, expiresAt, err
	}

	return int(user.ID), token, expiresAt, nil
}

func (l *LoginService) GenerateAndSaveToken(user model.User) (string, time.Time, error) {
	JWTService := JWTImpl{}
	token, expiresAt, err := JWTService.GenerateToken(int(user.ID), user.Email, true)
	if err != nil {
		zapLogger.Logger.Error("error in generating auth token", zap.Error(err))
		return token, expiresAt, err
	}

	err = JWTService.SaveToken(user, token, expiresAt)
	if err != nil {
		zapLogger.Logger.Error("error in saving auth token", zap.Error(err))
		return token, expiresAt, err
	}

	return token, expiresAt, nil
}

func generateOTP() string {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	r := rand.Intn(900000) + 100000
	otp := utilities.ConvertIntToString(r)

	return otp
}

func (l *LoginService) SendOTPService(phoneNumber string) (string, error) {
	client := twilioClient()
	otp := generateOTP()
	body := fmt.Sprintf("Your Pluto verification code is %s. This code will expire in 10 minutes.", otp)
	newPhoneNumber := fmt.Sprintf("+91" + phoneNumber)
	params := &twilioApi.CreateMessageParams{}
	params.SetTo(newPhoneNumber)
	params.SetFrom(config.ConfigValue.TwilioConfig.PhoneNumber)
	params.SetBody(body)

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		zapLogger.Logger.Error("failed to send verification SMS", zap.Error(err))
		return "", err
	}

	otpDao := &dao.UserVerificationOTPDao{
		Connection: *db.GlobalOrm,
	}
	userOtp := model.UserVerificationOTP{
		PhoneNumber: phoneNumber,
		OTP:         otp,
		ExpiresAt:   time.Now().Add(10 * time.Minute),
	}
	err = otpDao.Insert(userOtp)
	if err != nil {
		zapLogger.Logger.Error("error in inserting user otp to DB.", zap.Error(err))
		return "", err
	}

	return *resp.Sid, nil
}

func (l *LoginService) VerifyOTPService(phoneNumber string, OTP string) (bool, error) {
	otpDao := &dao.UserVerificationOTPDao{
		Connection: *db.GlobalOrm,
	}

	userOtpDetails, err := otpDao.FindByPhoneAndOTP(phoneNumber, OTP)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		zapLogger.Logger.Error("no record found, wrong otp")
		return false, errors.New("no record found, wrong otp")
	} else if err != nil {
		zapLogger.Logger.Error(fmt.Sprintf("error in finding otp details for phone number: %s", phoneNumber))
		return false, err
	}

	if time.Now().Before(userOtpDetails.ExpiresAt) {
		return true, nil
	} else {
		return false, errors.New("verification OTP expired")
	}

}

func (l *LoginService) CheckUserExistOrNot(userMobile string) (string, time.Time, error) {
	userDao := &dao.UserDao{
		Connection: *db.GlobalOrm,
	}

	user, err := userDao.FindByMobile(userMobile)
	if err == nil {
		if strings.ToLower(user.SignUpType) != "otp" {
			return "", time.Time{}, fmt.Errorf("otp signin didn't exists for this account. please sign in using %s method", user.SignUpType)
		}

		token, expiresAt, err := l.UserSignIN(user)
		if err != nil {
			log.Fatalln("error in user signin")
			return token, expiresAt, err
		} else {
			return token, expiresAt, nil
		}
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		userDetails := model.User{
			Mobile:     userMobile,
			SignUpType: "otp",
			IsActive:   true,
		}
		_, token, expiresAt, err := l.UserSignUP(userDetails)
		if err != nil {
			zapLogger.Logger.Error("error in registering user: ", zap.Error(err))
		}
		return token, expiresAt, nil
	}

	return "", time.Time{}, err
}

func (l *LoginService) GetUserDetailsFromGoogle(tokenInfo *dto.TokenInfo) model.User {
	userDetails := model.User{
		FirstName:  tokenInfo.GivenName,
		LastName:   tokenInfo.FamilyName,
		Mobile:     "",
		Code:       "",
		Email:      tokenInfo.Email,
		Password:   "",
		IsActive:   true,
		SignUpType: "social",
	}

	return userDetails
}

func (l *LoginService) SendVerificationEmail(emailId string) error {
	userDao := &dao.UserDao{
		Connection: *db.GlobalOrm,
	}

	user, err := userDao.FindByEmail(emailId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		zapLogger.Logger.Error("email does not exist. Please signup first")
		return err
	} else if err != nil {
		zapLogger.Logger.Error(fmt.Sprintf("error in getting user details from DB for emailId : %s", emailId))
		return err
	}

	verificationCode := uuid.New().String()
	name := user.FirstName
	verifyUrl := fmt.Sprintf("%s/user/verify/email?verification_code=%s", config.ConfigValue.BaseURL.URL, verificationCode)
	htmlBody := fmt.Sprintf(`<h1>Email Verification</h1>
    <p>Hello %s,</p>
    <p>Thank you for registering with Pluto. To complete your registration, please verify your email address by clicking the button below:</p>
    <a href="%s" style="display: inline-block; padding: 12px 24px; background-color: #4CAF50; color: white; text-decoration: none; border-radius: 4px;">Verify Email</a>
    <p>Best regards,</p>
    <p>Pluto</p>`, name, verifyUrl)
	title := "Verify you Pluto account"

	verificationDetails := model.EmailVerification{
		UserId:           int(user.ID),
		VerificationCode: verificationCode,
		ExpiresAt:        time.Now().Add(24 * time.Hour),
	}

	emailVerificationDao := &dao.EmailVerification{
		Connection: *db.GlobalOrm,
	}

	err = emailVerificationDao.Insert(verificationDetails)
	if err != nil {
		zapLogger.Logger.Error("error inserting email verification details")
		return err
	}

	s3Service := S3Service{}
	input := s3Service.SendEmailInput(emailId, htmlBody, title)

	err = s3Service.SendEmail(input)
	if err != nil {
		zapLogger.Logger.Error(fmt.Sprintf("Error sending email to account: %s", emailId))
		return err
	}

	return nil
}

func (l *LoginService) VerifyEmailService(verificationCode string) error {
	emailVerificationDao := &dao.EmailVerification{
		Connection: *db.GlobalOrm,
	}

	verificationDetails, err := emailVerificationDao.FindByVerificationCode(verificationCode)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		zapLogger.Logger.Error("verification code not found, wrong code")
		return err
	} else if err != nil {
		zapLogger.Logger.Error("error in getting email verification details")
	}

	if time.Now().Before(verificationDetails.ExpiresAt) {
		verificationDetails.IsVerified = true
	} else {
		zapLogger.Logger.Error("verification code expired")
		return errors.New("verification code expired")
	}

	userDao := &dao.UserDao{
		Connection: *db.GlobalOrm,
	}

	user, err := userDao.FindById(verificationDetails.UserId)
	if err != nil {
		zapLogger.Logger.Error(fmt.Sprintf("error in getting user details for userID: %d", verificationDetails.UserId))
	}

	if user.EmailVerified {
		zapLogger.Logger.Error("user email is already verified")
		return errors.New("email is already verified")
	} else {
		user.EmailVerified = true
	}

	verificationDetails, err = emailVerificationDao.UpdateEmailVerificationDetails(verificationDetails)
	if err != nil {
		zapLogger.Logger.Error("error updating email verification details")
		return err
	}

	user, err = userDao.UpdateUser(user)
	if err != nil {
		zapLogger.Logger.Error("error updating user details")
		return err
	}

	return nil
}

func (l *LoginService) DeleteUser(userID int, email string) error {
	userDao := &dao.UserDao{
		Connection: *db.GlobalOrm,
	}

	err := userDao.DeleteUser(userID, email)
	if err != nil {
		zapLogger.Logger.Error("error deleting user")
		return err
	}

	return nil
}

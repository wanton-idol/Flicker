package tests

import (
	"errors"
	"github.com/SuperMatch/model"
	mockDao "github.com/SuperMatch/pkg/db/dao/mocks"
	"github.com/SuperMatch/service/mocks"
	"github.com/golang/mock/gomock"
	"gorm.io/gorm"
	"net/http"
	"testing"
	"time"
)

func TestGoogleLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctrl.Finish()
	url := ""
	mockLogin := mocks.NewMockLoginInterface(ctrl)
	mockLogin.EXPECT().GoogleLogin().Return(url).AnyTimes()
}

func TestGoogleCallback(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctrl.Finish()
	mockLogin := mocks.NewMockLoginInterface(ctrl)
	mockLogin.EXPECT().GoogleCallback(gomock.Any(), gomock.Any()).
		Return(&http.Response{}, nil).
		AnyTimes()

}

func TestUserSignIN(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctrl.Finish()
	user := user()
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjc2ODE5OTA2NTgsInN1YiI6IjIzMTBuaWNAZ21haWwuY29tIiwidXNlcmRJZCI6OH0.Zm0jE9f5PzKzgaOuQK6yiQfXdtqFt_zcjJiqXS9OSVY"
	mockLogin := mocks.NewMockLoginInterface(ctrl)
	mockLogin.EXPECT().GenerateAndSaveToken(gomock.Eq(user)).
		Return(token, time.Now(), nil).
		AnyTimes()
}

func TestUserSignUP(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctrl.Finish()
	user := user()
	mockUserDao := mockDao.NewMockUserRepository(ctrl)
	mockUserDao.EXPECT().Insert(gomock.Any(), gomock.Eq(user)).
		Return(user, nil).
		AnyTimes()

	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjc2ODE5OTA2NTgsInN1YiI6IjIzMTBuaWNAZ21haWwuY29tIiwidXNlcmRJZCI6OH0.Zm0jE9f5PzKzgaOuQK6yiQfXdtqFt_zcjJiqXS9OSVY"
	mockLogin := mocks.NewMockLoginInterface(ctrl)
	mockLogin.EXPECT().GenerateAndSaveToken(gomock.Eq(user)).
		Return(token, time.Now(), nil).
		AnyTimes()
}

func TestGenerateAndSaveToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctrl.Finish()
	user := user()
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjc2ODE5OTA2NTgsInN1YiI6IjIzMTBuaWNAZ21haWwuY29tIiwidXNlcmRJZCI6OH0.Zm0jE9f5PzKzgaOuQK6yiQfXdtqFt_zcjJiqXS9OSVY"
	mockLogin := mocks.NewMockLoginInterface(ctrl)
	mockLogin.EXPECT().GenerateAndSaveToken(gomock.Eq(user)).
		Return(token, time.Now(), nil).
		AnyTimes()

}

func TestSendOTPService(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctrl.Finish()
	phoneNumber := "9874563210"
	SID := ""
	mockLogin := mocks.NewMockLoginInterface(ctrl)
	mockLogin.EXPECT().SendOTPService(gomock.Eq(phoneNumber)).
		Return(SID, nil).
		AnyTimes()
}

func TestVerifyOTPService(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctrl.Finish()
	userOtpDetails := model.UserVerificationOTP{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		PhoneNumber: "9874563210",
		OTP:         "715346",
		ExpiresAt:   time.Now().Add(10 * time.Minute),
	}
	mockOtp := mockDao.NewMockUserVerificationOTPRepository(ctrl)
	mockOtp.EXPECT().FindByPhoneAndOTP(gomock.Eq(userOtpDetails.PhoneNumber), gomock.Eq(userOtpDetails.OTP)).
		Return(userOtpDetails, nil).
		AnyTimes()

	if time.Now().After(userOtpDetails.ExpiresAt) {
		t.Errorf("otp expired. user verification failed.")
	}
}

func TestCheckUserExistOrNot(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctrl.Finish()
	user := user()
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjc2ODE5OTA2NTgsInN1YiI6IjIzMTBuaWNAZ21haWwuY29tIiwidXNlcmRJZCI6OH0.Zm0jE9f5PzKzgaOuQK6yiQfXdtqFt_zcjJiqXS9OSVY"
	mockUserDao := mockDao.NewMockUserRepository(ctrl)
	var err error
	mockUserDao.EXPECT().FindByMobile(gomock.Eq(user.Mobile)).Return(user, err).AnyTimes()
	mockLogin := mocks.NewMockLoginInterface(ctrl)

	if err == nil {
		mockLogin.EXPECT().UserSignIN(gomock.Eq(user)).Return(token, time.Now(), nil).AnyTimes()

	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		userDetails := model.User{
			Mobile:     user.Mobile,
			SignUpType: "otp",
			IsActive:   true,
		}
		mockLogin.EXPECT().UserSignUP(gomock.Eq(userDetails)).Return(token, time.Now(), nil).AnyTimes()
	}
}

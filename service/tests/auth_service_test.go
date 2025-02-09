package tests

import (
	"github.com/SuperMatch/model"
	mocks "github.com/SuperMatch/pkg/db/dao/mocks"
	"github.com/SuperMatch/service"
	"github.com/golang/mock/gomock"
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestGenerateToken(t *testing.T) {
	userId := 1
	email := "test@example.com"
	jwtService := service.JWTImpl{}
	token, expiresAt, err := jwtService.GenerateToken(userId, email, true)
	claims, err := jwtService.ValidateToken(token)
	if claims.UserID != userId {
		t.Errorf("token is not valid %q", token)
	}

	if expiresAt.IsZero() {
		t.Errorf("expiresAt is Zero %q", expiresAt)
	}

	if err != nil {
		t.Errorf("error in validate token: %v", err)
	}
}

func TestValidateToken(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjc2ODE5OTA2NTgsInN1YiI6IjIzMTBuaWNAZ21haWwuY29tIiwidXNlcmRJZCI6OH0.Zm0jE9f5PzKzgaOuQK6yiQfXdtqFt_zcjJiqXS9OSVY"
	jwtService := service.JWTImpl{}
	claims, err := jwtService.ValidateToken(token)
	if claims.UserID == 0 {
		t.Errorf("Token is not valid %q", token)
	}
	if err != nil {
		t.Errorf("error in validate token: %v", err)
	}
}

func TestValidateTokenFromDatabase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUser := mocks.NewMockUserRepository(ctrl)
	claims := service.AuthClaims{
		UserID: 1,
	}
	user := user()
	mockUser.EXPECT().
		FindById(gomock.Any(), gomock.Eq(claims.UserID)).
		Return(user, nil).
		AnyTimes()

	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjc2ODE5OTA2NTgsInN1YiI6IjIzMTBuaWNAZ21haWwuY29tIiwidXNlcmRJZCI6OH0.Zm0jE9f5PzKzgaOuQK6yiQfXdtqFt_zcjJiqXS9OSVY"
	mockUserToken := mocks.NewMockUserTokenRepository(ctrl)
	userToken := userToken()
	mockUserToken.EXPECT().
		FindByTokenAndUserId(gomock.Any(), gomock.Eq(token), gomock.Eq(claims.UserID)).
		Return(userToken, nil).
		AnyTimes()
}

func user() model.User {
	return model.User{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		FirstName:  "John",
		LastName:   "Doe",
		Email:      "test@example.com",
		Password:   "password",
		Code:       "91",
		Mobile:     "9874563210",
		IsActive:   true,
		SignUpType: "social",
		CreatedAt:  time.Now(),
		DeletedAt:  nil,
		UpdatedAt:  time.Now(),
	}
}

func userToken() model.UserToken {
	return model.UserToken{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		UserId:    1,
		Token:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjc2ODE5OTA2NTgsInN1YiI6IjIzMTBuaWNAZ21haWwuY29tIiwidXNlcmRJZCI6OH0.Zm0jE9f5PzKzgaOuQK6yiQfXdtqFt_zcjJiqXS9OSVY",
		IsActive:  true,
		ExpiresAt: time.Time{},
	}
}

func TestSaveToken(t *testing.T) {
	user := user()
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjc2ODE5OTA2NTgsInN1YiI6IjIzMTBuaWNAZ21haWwuY29tIiwidXNlcmRJZCI6OH0.Zm0jE9f5PzKzgaOuQK6yiQfXdtqFt_zcjJiqXS9OSVY"
	expiresAt := time.Now()

	userToken := model.UserToken{
		UserId:    int(user.Model.ID),
		Token:     token,
		IsActive:  true,
		ExpiresAt: expiresAt,
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockToken := mocks.NewMockUserTokenRepository(ctrl)
	mockToken.EXPECT().Insert(gomock.Any(), gomock.Eq(userToken)).
		Return(userToken, nil).
		AnyTimes()
}

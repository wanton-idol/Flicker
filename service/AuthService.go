package service

import (
	"context"
	"errors"
	"github.com/SuperMatch/zapLogger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"time"

	"github.com/SuperMatch/model"
	"github.com/SuperMatch/pkg/db"
	"github.com/SuperMatch/pkg/db/dao"
	"github.com/golang-jwt/jwt"
)

type JWTService interface {
	GenerateToken(userId int64, email string, isUser bool) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
	SaveToken(user model.User, token string, expiresAt time.Time) error
	ValidateTokenFromDatabase(token string) (model.UserToken, error)
}

const SecretKey = "secretKey"

type JWTImpl struct {
}

type AuthClaims struct {
	jwt.StandardClaims
	UserID int `json:"userdId"`
}

func (j *JWTImpl) GenerateToken(userId int, email string, isUser bool) (string, time.Time, error) {

	expiresAt := time.Now().Add(100000000 * time.Minute)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, AuthClaims{
		StandardClaims: jwt.StandardClaims{
			Subject:   email,
			ExpiresAt: expiresAt.Unix(),
		},
		UserID: int(userId),
	})

	signedToken, err := token.SignedString([]byte(SecretKey))

	if err != nil {
		return "", time.Now(), err
	}
	return signedToken, expiresAt, nil
}

func (j *JWTImpl) ValidateToken(token string) (AuthClaims, error) {

	var claims AuthClaims

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("invalid token")
		}

		return []byte(SecretKey), nil
	}

	_, err := jwt.ParseWithClaims(token, &claims, keyFunc)

	if err != nil {
		return AuthClaims{}, err
	}

	return claims, nil
}

func (j *JWTImpl) ValidateTokenFromDatabase(c *gin.Context, claims AuthClaims, token string) {

	newUserRepo := dao.UserDao{
		Connection: *db.GlobalOrm,
	}

	_, err := newUserRepo.FindById(claims.UserID)

	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"message": "no active user found."})
		return
	}

	tokenDao := &dao.UserTokenDao{
		Connection: *db.GlobalOrm,
	}

	userToken, err := tokenDao.FindByTokenAndUserId(context.Background(), token, claims.UserID)

	if err != nil || !userToken.IsActive {
		zapLogger.Logger.Debug("no user found with active token in database")
		c.AbortWithStatusJSON(401, gin.H{"message": "token expird."})
		return
	}
}

func (j *JWTImpl) SaveToken(user model.User, token string, expiresAt time.Time) error {

	tokenDao := &dao.UserTokenDao{
		Connection: *db.GlobalOrm,
	}

	userToken := model.UserToken{
		UserId:    int(user.ID),
		Token:     token,
		IsActive:  true,
		ExpiresAt: expiresAt,
	}

	_, err := tokenDao.Insert(context.Background(), userToken)
	if err != nil {
		zapLogger.Logger.Error("userToken is not inserted", zap.Error(err))
		return err
	}
	return nil
}

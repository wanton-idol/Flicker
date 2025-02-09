package middleware

import (
	"github.com/SuperMatch/zapLogger"
	"strconv"

	Service "github.com/SuperMatch/service"
	"github.com/gin-gonic/gin"
)

func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		zapLogger.Logger.Debug("auth middleware is checking authentication for this request")
		authToken := c.GetHeader("token")
		userId := c.GetHeader("user_id")
		if authToken == "" {
			zapLogger.Logger.Debug("auth middleware has found no token in this request")
			c.AbortWithStatusJSON(401, gin.H{"message": "token is missing"})
			return
		}
		if userId == "" {
			zapLogger.Logger.Debug("auth middleware has found no userID in this request")
			c.AbortWithStatusJSON(401, gin.H{"message": "user_id is missing"})
			return
		}

		validateToken(c, userId, authToken)

		c.Next()
		zapLogger.Logger.Debug("auth middleware has checked authentication for this request")
	}
}

func validateToken(c *gin.Context, userId, token string) {

	zapLogger.Logger.Debug("auth middleware is validating token")
	authservice := &Service.JWTImpl{}

	claims, err := authservice.ValidateToken(token)

	if err != nil {
		zapLogger.Logger.Debug("auth middleware has found invalid token")
		c.AbortWithStatusJSON(401, gin.H{"message": "invalid token"})
		return
	}
	userIDD, err := strconv.Atoi(userId)

	if err != nil {
		zapLogger.Logger.Debug("auth middleware has found invalid token")
		c.AbortWithStatusJSON(401, gin.H{"message": "invalid userId"})
		return
	}

	if claims.UserID != userIDD {
		zapLogger.Logger.Debug("auth middleware has found invalid token")
		c.AbortWithStatusJSON(401, gin.H{"message": "invalid token"})
		return
	}

	authservice.ValidateTokenFromDatabase(c, claims, token)

	zapLogger.Logger.Debug("auth middleware has found valid token")
}

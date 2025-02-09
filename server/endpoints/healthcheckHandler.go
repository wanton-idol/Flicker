package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthCheckHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Status(http.StatusOK)
	}
}

func CauseErr() func(c *gin.Context) {
	return func(c *gin.Context) {
		panic("artificial health check failure")
	}
}

func CausePanic() func(c *gin.Context) {
	return func(c *gin.Context) {
		panic("artificial panic")
	}
}

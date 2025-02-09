package server

import (
	"github.com/SuperMatch/config"
	"github.com/SuperMatch/server/endpoints"
	"github.com/SuperMatch/zapLogger"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type RouterDeps struct {
	AppEnv    string
	AppConfig config.Config
}

func ConfigRouter(deps *RouterDeps) (*gin.Engine, error) {

	router := gin.New()

	router.Use(zapLogger.LogMiddleware())
	addHealthRouter(router)

	//swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL(deps.AppConfig.HTTP.IP+":"+deps.AppConfig.HTTP.PORT+"/swagger/doc.json"),
		ginSwagger.DefaultModelsExpandDepth(-1))

	router.POST("/user/register", endpoints.UserRegistrationHandler)
	router.POST("/user/login", endpoints.LoginUser)
	router.POST("/user/google/login", endpoints.GoogleLoginHandler)
	router.POST("/user/send/otp", endpoints.SendOTP)
	router.POST("/user/verify/otp", endpoints.VerifyOTP)
	router.POST("/user/verify/email", endpoints.SendVerificationEmail)
	router.GET("/user/verify/email", endpoints.VerifyEmail)

	//elasticsearch indexing
	router.POST("/create/user_profile/index", endpoints.CreateProfileIndex)
	router.POST("/create/user_stories/index", endpoints.CreateStoriesIndex)

	//Image uploader

	//router.Use(middleware.AuthMiddleWare())
	//User APIs
	router.POST("/user/upgradePremium", endpoints.PremiumUpgradeHandler)
	router.POST("/user/profile", endpoints.CreateProfileHandler)
	router.PUT("/user/profile", endpoints.UpdateUserProfileHandler)
	router.GET("/user/profile", endpoints.GetUserProfileHandler)
	router.PUT("/user/updateSearchProfile", endpoints.UpdateSearchProfileHandler)
	router.GET("/user/searchProfile", endpoints.GetUserSearchProfileHandler)
	router.POST("/user/updateLocation", endpoints.UpdateLocationHandler)
	router.POST("/user/delete", endpoints.DeleteUserHandler)

	//user media
	router.POST("/user/profileMedia", endpoints.SaveMediaHandler)
	router.POST("/user/update/profileMedia", endpoints.UpdateMediaHandler)
	router.DELETE("/user/profileMedia", endpoints.DeleteMediaHandler)
	router.GET("/user/profileMedia", endpoints.GetUserProfileMediaHandler)
	router.PUT("/user/advanced-filters", endpoints.UpdateAdvancedFilterHandler)
	router.GET("/user/matches", endpoints.GetUserMatchHandler)
	router.GET("/user/likes", endpoints.GetUserLikesHandler)

	//Public use APIS
	router.GET("/searchProfile", endpoints.SearchProfileHandler)
	router.POST("/user/swipe", endpoints.SwipeHandler)
	router.GET("/interests", endpoints.GetInterests)
	router.POST("/user/interests", endpoints.CreateUserInterests)
	router.GET("/user/interests", endpoints.GetUserInterests)
	router.PUT("/user/interests", endpoints.UpdateUserInterests)
	router.GET("/nudges", endpoints.GetNudges)
	router.POST("/user/nudge", endpoints.CreateUserNudge)
	router.POST("/user/nudge/media", endpoints.CreateUserNudgeWithMedia)
	router.GET("/user/nudges", endpoints.GetUserNudges)
	router.PUT("/user/nudge", endpoints.UpdateUserNudge)
	router.DELETE("/user/nudge", endpoints.DeleteUserNudge)
	router.GET("/filters", endpoints.GetFiltersHandler)

	//Chat APIs
	router.POST("/chat/message", endpoints.SaveMessage)
	router.GET("/chat/user/chats", endpoints.RetrieveUserChats)
	router.GET("/chat/user/list", endpoints.GetUserChatsList)
	router.PUT("/chat/messages/status", endpoints.UpdateMessagesStatus)
	router.GET("/chat/last/messages", endpoints.GetLastMessages)

	//Stories APIs
	router.POST("/user/stories/index", endpoints.IndexUserStories)
	router.GET("/user/stories/search/profileID", endpoints.GetUserStoriesByProfileID)
	router.GET("/user/stories/search/location", endpoints.GetUserStoriesByLocation)

	//event APIs
	router.POST("/event/index", endpoints.CreateEventIndexHandler)
	router.POST("/user/event", endpoints.CreateUserEventsHandler)
	router.GET("/user/event", endpoints.GetUserEventsHandler)
	router.PUT("/user/event", endpoints.UpdateUserEventHandler)
	router.DELETE("/user/event", endpoints.DeleteUserEventHandler)
	router.GET("/events/search", endpoints.SearchEventsHandler)

	router.POST("/user/device/token", endpoints.GetDeviceToken)
	router.POST("/user/send/notification", endpoints.SendNotificationToUser)

	return router, nil
}

func addHealthRouter(router *gin.Engine) {
	router.GET("/healthCheck", endpoints.HealthCheckHandler())
}

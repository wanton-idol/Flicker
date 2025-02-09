package endpoints

import (
	"encoding/json"
	"github.com/SuperMatch/model"
	elasticsearchPkg "github.com/SuperMatch/model/elasticSearch"
	Service "github.com/SuperMatch/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// CreateStoriesIndex godoc
//
//	@Security		ApiKeyAuth
//	@Summary		CreateStoriesIndex
//	@Description	Create User Stories Index
//	@Tags			Stories
//	@Accept			json
//	@Produce		json
//	@Success		200							{string}	string	"index created successfully"
//	@Failure		400							{string}	string	Bad	request
//	@Failure		500							{string}	string	"internal server error"
//	@Router			/create/user_stories/index	[POST]
func CreateStoriesIndex(c *gin.Context) {
	userStoriesService := Service.NewUserStoriesService()
	err := userStoriesService.CreateStoriesIndex()
	if err != nil {
		c.JSON(500, gin.H{
			"message": "error in creating index in elasticSearch",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{"message": "index created successfully"})
}

// CreateUserStories godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Create User Stories
//	@Description	Create User Stories
//	@Tags			Stories
//	@Accept			mpfd
//	@Produce		json
//	@Param			media				formData	file					true	"Media"
//	@Param			values				formData	elasticsearchPkg.Values	true	"Values"
//	@Param			user_id				header		int						true	"User ID"
//	@Param			user_profile_id		header		int						true	"User Profile ID"
//	@Success		200					{string}	string					"user stories indexed successfully"
//	@Failure		400					{string}	string					Bad	request
//	@Failure		500					{string}	string					"internal server error"
//	@Router			/user/stories/index	[POST]
func IndexUserStories(c *gin.Context) {
	id := c.GetHeader("user_id")

	id2 := c.GetHeader("user_profile_id")
	userProfileID, err := strconv.Atoi(id2)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	form, _ := c.MultipartForm()
	files := form.File["media"]
	values := form.Value["values"]

	valuesData := []byte(values[0])
	var storyValues elasticsearchPkg.Values
	err = json.Unmarshal(valuesData, &storyValues)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var (
		mediaUrl  string
		mediaType string
	)
	userStoriesService := Service.NewUserStoriesService()
	if files != nil {
		url, fileType, err := userStoriesService.UploadFileToS3(id, files[0])
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "error in uploading file to s3", "error": err.Error()})
			return
		}
		mediaUrl = url
		mediaType = fileType
	}

	userStories := elasticsearchPkg.UserStories{
		UserProfileID: userProfileID,
		Text:          storyValues.Text,
		MediaURL:      mediaUrl,
		MediaType:     mediaType,
		Location:      storyValues.Location,
	}

	err = userStoriesService.IndexUserStories(userStories)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error in indexing user stories",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user stories indexed successfully"})

}

// GetUserStoriesByProfileID godoc
//
//	@Security		ApiKeyAuth
//	@Summary		GetUserStoriesByProfileID
//	@Description	Get user stories based on profile id
//	@Tags			Stories
//	@Accept			json
//	@Produce		json
//	@Param			user_profile_id					header		int	true	"User Profile ID"
//	@Success		200								{object}	elasticsearchPkg.UserStories
//	@Failure		400								{string}	string	Bad	request
//	@Failure		500								{string}	string	"internal server error"
//	@Router			/user/stories/search/profileID	[GET]
func GetUserStoriesByProfileID(c *gin.Context) {
	id := c.GetHeader("user_profile_id")
	userProfileID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userStoriesService := Service.NewUserStoriesService()
	userStories, err := userStoriesService.GetUserStoriesByProfileID(userProfileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error in getting user stories by user profile id",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"userStories": userStories})
}

// GetUserStoriesByLocation godoc
//
//	@Security		ApiKeyAuth
//	@Summary		GetUserStoriesByLocation
//	@Description	Get user stories based on location
//	@Tags			Stories
//	@Accept			json
//	@Produce		json
//	@Param			location						body		model.UserLocation	true	"Location"
//	@Success		200								{object}	elasticsearchPkg.UserStories
//	@Failure		400								{string}	string	Bad	request
//	@Failure		500								{string}	string	"internal server error"
//	@Router			/user/stories/search/location	[GET]
func GetUserStoriesByLocation(c *gin.Context) {
	var location model.UserLocation
	if err := c.ShouldBindJSON(&location); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userStoriesService := Service.NewUserStoriesService()
	userStories, err := userStoriesService.GetUserStoriesByLocation(location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error in getting user stories by location",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"userStories": userStories})
}

package endpoints

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/SuperMatch/zapLogger"

	"github.com/SuperMatch/model"
	"gorm.io/gorm"

	dto "github.com/SuperMatch/model/dto"
	Service "github.com/SuperMatch/service"
	"github.com/gin-gonic/gin"
)

const user_profile_S3_bucket = "user-profile-supermatch"

// CreateProfileHandler godoc
//
//	@Security		ApiKeyAuth
//	@Summary		CreateUserProfile
//	@Description	API to create the user profile in elasticsearch and database for the first time
//	@Tags			Profile
//	@Accept			json
//	@Produce		json
//	@Param			userProfile	body		dto.UserProfile	true	"userProfile"
//	@Param			user_id		header		string			true	"user_id"
//	@Success		200			{object}	dto.UserProfile	"user Profile created."
//	@Failure		400			{string}	string			"Bad request"
//	@Failure		500			{string}	string			"Internal Server Error"
//	@Router			/user/profile [post]

func CreateProfileHandler(c *gin.Context) {
	userID, _ := strconv.Atoi(c.GetHeader("user_id"))

	var profileDTO dto.UserProfile
	if err := c.ShouldBindJSON(&profileDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userProfileService := Service.NewUserProfileService()
	actualUserProfile, err := userProfileService.GetUserProfileFromDB(userID)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	}

	actualUserProfile.UserId = userID
	profileDTO.UserId = userID
	profileDTO, err = userProfileService.CreateUserProfile(actualUserProfile, profileDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user Profile created successfully.", "data": profileDTO})
}

// UpdateUserProfileHandler godoc
//
//	@Security		ApiKeyAuth
//	@Summary		updateUserProfile
//	@Description	API to store the user profile in elasticsearch and database
//	@Tags			Profile
//	@Accept			json
//	@Produce		json
//	@Param			userProfile	body		dto.UserProfile	true	"userProfile"
//	@Param			userID		header		string			true	"user_id"
//	@Success		200			{object}	dto.UserProfile	"user Profile updated."
//	@Failure		400			{string}	string			"Bad request"
//	@Failure		500			{string}	string			"Internal Server Error"
//	@Router			/user/profile [put]
func UpdateUserProfileHandler(c *gin.Context) {
	id := c.Request.Header.Get("user_id")
	zapLogger.Logger.Debug(fmt.Sprintf("user_id=%s", id))
	userID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var profileDTO dto.UserProfile
	if err := c.ShouldBindJSON(&profileDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userProfileService := Service.NewUserProfileService()
	actualUserProfile, err := userProfileService.GetUserProfileFromDB(userID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	actualUserProfile.UserId = userID
	profileDTO, err = userProfileService.UpdateUserProfile(actualUserProfile, profileDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user Profile updated.", "data": profileDTO})
}

// GetUserProfileHandler godoc
//
//	@Security		ApiKeyAuth
//	@Summary		getUserProfile
//	@Description	API to get the user profile
//	@Tags			Profile
//	@Produce		json
//	@Param			user_id	header		string			true	"user_id"
//	@Success		200		{object}	dto.UserProfile	"successfully received profile."
//	@Failure		400		{string}	Bad				request
//	@Failure		500		{string}	Internal		Server	Error
//	@Header			all		{string}	token			"token"
//	@Header			all		{string}	userID			"user_id"
//	@Router			/user/profile [get]
func GetUserProfileHandler(c *gin.Context) {

	userID := c.Request.Header.Get("user_id")
	id, _ := strconv.Atoi(userID)
	userProfileService := Service.NewUserProfileService()
	userProfile, err := userProfileService.GetUserProfileFromDB(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error in fetching user profile", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully received profile.", "data": userProfile})
}

// UpdateSearchProfileHandler godoc
//
//	@Security		ApiKeyAuth
//	@Summary		updateUserSearchProfile
//	@Description	API to update the search profile for a user
//	@Tags			Profile
//	@Accept			json
//	@Produce		json
//	@Param			userSearchProfile	body		dto.UserSearchProfile	true	"userSearchProfile"
//	@Param			userID				header		string					true	"user_id"
//	@Success		200					{object}	dto.UserProfile			"user Search Profile updated."
//	@Failure		400					{string}	string					"Bad request"
//	@Failure		500					{string}	string					"Internal Server Error"
//	@Router			/user/updateSearchProfile [put]
func UpdateSearchProfileHandler(c *gin.Context) {

	userID := c.Request.Header.Get("user_id")
	id, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userSearchProfile dto.UserSearchProfile
	if err := c.ShouldBindJSON(&userSearchProfile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userProfileService := Service.NewUserProfileService()

	userProfile, err := userProfileService.GetUserProfileFromDB(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error in fetching user profile", "error": err.Error()})
		return
	}
	userSearchProfileService := Service.NewUserSearchProfileService()
	userSearchProfile, err = userSearchProfileService.UpdateUserSearchProfile(context.Background(), userProfile, userSearchProfile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user Search Profile updated.", "data": userSearchProfile})
}

// GetUserSearchProfileHandler godoc
//
//	@Security		ApiKeyAuth
//	@Summary		getUserSearchProfile
//	@Description	API to get the search profile for a user
//	@Tags			Profile
//	@Accept			json
//	@Produce		json
//	@Param			userID	header		string			true	"user_id"
//	@Success		200		{object}	dto.UserProfile	"user Search Profile found."
//	@Failure		400		{string}	string			"Bad request"
//	@Failure		500		{string}	string			"Internal Server Error"
//	@Router			/user/searchProfile [get]
func GetUserSearchProfileHandler(c *gin.Context) {
	userID := c.Request.Header.Get("user_id")
	_, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	profileID := c.Request.Header.Get("profile_id")
	profileId, err := strconv.Atoi(profileID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userSearchProfileService := Service.NewUserSearchProfileService()
	userSearchProfile, err := userSearchProfileService.FindByProfileId(c, profileId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error in fetching user profile", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully received user search profile.", "data": userSearchProfile})
}

// SearchProfileHandler godoc
//
//	@Security		ApiKeyAuth
//	@Summary		UserSearchProfile
//	@Description	API to get user search profile
//	@Tags			Profile
//	@Produce		json
//	@Param			userID	header		string							true	"user_id"
//	@Success		200		{object}	elasticsearchPkg.UserProfile	"successfully received profiles."
//	@Failure		400		{string}	string							"Bad request"
//	@Failure		500		{string}	string							"Internal Server Error"
//	@Router			/searchProfile [get]
func SearchProfileHandler(c *gin.Context) {

	userID := c.Request.Header.Get("user_id")

	id, err := strconv.Atoi(userID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userProfileService := Service.NewUserProfileService()
	userProfile, err := userProfileService.GetUserProfile(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profiles, err := userProfileService.SearchProfile(userProfile)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error in fetching user profile", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully received profiles.", "data": profiles})
}

// PremiumUpgradeHandler godoc
//
//	@Security		ApiKeyAuth
//	@Summary		PremiumUpgrade
//	@Description	API to upgrade the profile to premium
//	@Tags			Profile
//	@Accept			json
//	@Produce		json
//	@Param			userID	header		string			true	"user_id"
//	@Success		200		{object}	dto.UserProfile	"successfully received profile."
//	@Failure		400		{string}	string			"Bad request"
//	@Failure		500		{string}	string			"Internal Server Error"
//	@Router			/user/upgradePremium [post]
func PremiumUpgradeHandler(c *gin.Context) {
	userID := c.Request.Header.Get("user_id")
	id, _ := strconv.Atoi(userID)

	userProfileService := Service.NewUserProfileService()
	userProfile, err := userProfileService.GetUserProfileFromDB(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedProfile, err := userProfileService.UpdatePremiumProfile(userProfile)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "successfully updated user profile.", "data": updatedProfile})
}

// SaveMediaHandler godoc
//
//	@Security		ApiKeyAuth
//	@Summary		SaveProfileMedia
//	@Description	API to save images to S3 and DB
//	@Tags			Media
//	@Accept			mpfd
//	@Produce		json
//	@Param			images	formData	file	true	"ImageToUpload"
//	@Param			userID	header		string	true	"user_id"
//	@Success		200		{string}	string	"successfully uploaded files."
//	@Failure		400		{string}	string	"Bad request"
//	@Failure		500		{string}	string	"Internal Server Error"
//	@Router			/user/profileMedia [post]
func SaveMediaHandler(c *gin.Context) {

	userId := c.Request.Header.Get("user_id")
	userProfileId := c.Request.Header.Get("user_profile_id")

	id, err := strconv.Atoi(userId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//to be used multiple profile
	_, err = strconv.Atoi(userProfileId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error in getting user_profile_id.", "error": err.Error()})
		return
	}

	userProfileService := Service.NewUserProfileService()
	s3service := Service.NewS3Service()

	userProfile, _ := userProfileService.GetUserProfileFromDB(id)

	if err != nil {
		zapLogger.Logger.Error(fmt.Sprintf("error in fetching user profile for user_profile_id=%s", userProfileId))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error in fetching user profile", "error": err.Error()})
		return
	}

	form, _ := c.MultipartForm()
	files := form.File["image"]

	var profileRequest dto.ProfileMediaRequestDTO
	if err := c.ShouldBind(&profileRequest); err != nil {
		c.String(http.StatusBadRequest, "Invalid form data: %v", err)
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fileExt := filepath.Ext(files[0].Filename)
	if !userProfileService.CheckAllowedFileType(fileExt) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "file extension is not supported.", "error": err.Error()})
		return
	}

	filename := fmt.Sprintf("%v", time.Now().UTC().Format("2006-01-02T15:04:05.00000")) + fileExt
	tempFile, _ := files[0].Open()
	S3filepath := userId + "/profile/" + filename

	result, err := s3service.UploadFileToS3(user_profile_S3_bucket, S3filepath, tempFile, filename)
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"message": "error in uploading file to S3", "error": err.Error()})
		return
	}

	userMedia, err := userProfileService.SaveProfileMedia(userProfile, result, profileRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error in saving profile media.", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "file uploaded successfully.", "data": userMedia})
}

// GetUserProfileMediaHandler godoc
//
//	@Security		ApiKeyAuth
//	@Summary		GetUserProfileMedia
//	@Description	Get Profile media with order and s3 signed URL
//	@Tags			Media
//	@Produce		json
//	@Param			userID	header		string	true	"user_id"
//	@Success		200		{string}	string	"successfully received profile media."
//	@Failure		400		{string}	string	"Bad request"
//	@Failure		500		{string}	string	"Internal Server Error"
//	@Router			/user/profileMedia [get]
func GetUserProfileMediaHandler(c *gin.Context) {
	//get userID from Query Params

	userID := c.Request.URL.Query().Get("user_id")

	id, _ := strconv.Atoi(userID)
	userProfileService := Service.NewUserProfileService()
	data, err := userProfileService.GetProfileMediaByUserId(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error in fetching user media.", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "successfully received profile media.", "data": data})
}

// DeleteMediaHandler godoc
//
//	@Security		ApiKeyAuth
//	@Summary		DeleteMediaAPI
//	@Description	API to Delete Profile media.
//	@Tags			Media
//	@Accept			json
//	@Produce		json
//	@Param			mediaId	query		int		true	"MediaID"
//	@Param			userID	header		string	true	"user_id"
//	@Success		200		{string}	string	"successfully deleted media."
//	@Failure		400		{string}	string	"Bad request"
//	@Failure		500		{string}	string	"Internal Server Error"
//	@Router			/user/profileMedia [delete]
func DeleteMediaHandler(c *gin.Context) {
	userID := c.Request.Header.Get("user_id")
	mediaID, ok := c.GetQuery("mediaId")

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "mediaId is required."})
		return
	}
	id, _ := strconv.Atoi(userID)
	media, _ := strconv.Atoi(mediaID)

	userProfileService := Service.NewUserProfileService()
	err := userProfileService.RemoveProfileMedia(id, media)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "successfully deleted media."})
}

// UpdateAdvancedFilterHandler godoc
//
//	@Security		ApiKeyAuth
//	@Summary		updateAdvancedFilter
//	@Description	API to update the advanced filter for a user
//	@Tags			Profile
//	@Accept			json
//	@Produce		json
//	@Param			advancedFilter	body		dto.AdvancedFilter	true	"advancedFilter"
//	@Param			userID			header		string				true	"user_id"
//	@Success		200				{object}	dto.AdvancedFilter	"advanced filters updated."
//	@Failure		400				{string}	string				"Bad request"
//	@Failure		500				{string}	string				"Internal Server Error"
//	@Router			/user/advancedFilter [post]
func UpdateAdvancedFilterHandler(c *gin.Context) {

	userID := c.Request.Header.Get("user_id")
	id, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var advancedFilters dto.AdvancedFilter

	if err := c.ShouldBindJSON(&advancedFilters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userProfileService := Service.NewUserProfileService()
	userProfile, err := userProfileService.GetUserProfileFromDB(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error in getting user profile.", "error": err.Error()})
		return
	}

	userSearchProfileService := Service.NewUserSearchProfileService()
	advFilters, err := userSearchProfileService.UpdateAdvancedFilters(advancedFilters, userProfile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error in updating advanced filters.", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "advanced filters updated.", "data": advFilters})
}

// GetInterests godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Get Interests_category List
//	@Description	API to get interest list
//	@Tags			Profile
//	@Accept			json
//	@Produce		json
//	@Param			userID	header		string						true	"user_id"
//	@Success		200		{object}	model.InterestsListResponse	"Fetched interests list"
//	@Failure		400		{string}	string						"Bad request"
//	@Failure		500		{string}	string						"Internal Server Error"
//	@Router			/interests [get]
func GetInterests(c *gin.Context) {
	userProfileService := Service.NewUserProfileService()
	interests, err := userProfileService.GetInterestsCategory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error in getting interests list.", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"interests": interests})
}

// CreateUserInterests godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Create User Interests_category
//	@Description	API to create user interests
//	@Tags			Profile
//	@Accept			json
//	@Produce		json
//	@Param			interests	body		model.InterestsCategory	true	"userInterests"
//	@Param			userID		header		string					true	"user_id"
//	@Success		200			{string}	string					"UserInterests created successfully"
//	@Failure		400			{string}	string					"Bad request"
//	@Failure		500			{string}	string					"Internal Server Error"
//	@Router			/user/interests [post]
func CreateUserInterests(c *gin.Context) {
	userID := c.Request.Header.Get("user_id")
	id, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var interests model.InterestsCategory
	if err := c.ShouldBindJSON(&interests); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userProfileService := Service.NewUserProfileService()
	userInterests, err := userProfileService.CreateUserInterest(interests, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error in creating user interests", "error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "user interests are created successfully", "data": userInterests})

}

// GetUserInterests godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Get User Interests_category List
//	@Description	API to get user interest list
//	@Tags			Profile
//	@Accept			json
//	@Produce		json
//	@Param			userID	header		string					true	"user_id"
//	@Success		200		{object}	model.InterestsCategory	"Fetched user interests list"
//	@Failure		400		{string}	string					"Bad request"
//	@Failure		500		{string}	string					"Internal Server Error"
//	@Router			/user/interests [get]
func GetUserInterests(c *gin.Context) {
	userID := c.Request.Header.Get("user_id")
	id, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userProfileService := Service.NewUserProfileService()
	userInterests, err := userProfileService.GetUserInterests(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "couldn't find user interests", "error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"data": userInterests})
}

// GetNudges godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Get Nudges List
//	@Description	API to get nudges list
//	@Tags			Profile
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	model.Nudge	"Fetched nudges list"
//	@Failure		400	{string}	string		"Bad request"
//	@Failure		500	{string}	string		"Internal Server Error"
//	@Router			/nudges [get]
func GetNudges(c *gin.Context) {
	userProfileService := Service.NewUserProfileService()
	nudges, err := userProfileService.GetNudgesService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error in getting nudges list.", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"nudges": nudges})
}

// CreateUserNudge godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Create User Nudge
//	@Description	API to create user nudge
//	@Tags			Profile
//	@Accept			json
//	@Produce		json
//	@Param			nudge	body		model.NudgeDetail	true	"userNudge"
//	@Param			userID	header		string				true	"user_id"
//	@Success		200		{object}	model.NudgeDetail	"UserNudge created successfully"
//	@Failure		400		{string}	string				"Bad request"
//	@Failure		500		{string}	string				"Internal Server Error"
//	@Router			/user/nudge [post]
func CreateUserNudge(c *gin.Context) {
	userID := c.Request.Header.Get("user_id")
	id, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var nudgesDetail model.NudgeDetail
	if err := c.ShouldBindJSON(&nudgesDetail); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userProfileService := Service.NewUserProfileService()
	userNudge, err := userProfileService.CreateUserNudgeService(nudgesDetail, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error in creating user nudge", "error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "user nudge created successfully", "data": userNudge})

}

// GetUserNudges godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Get User nudges List
//	@Description	API to get user nudges list
//	@Tags			Profile
//	@Accept			json
//	@Produce		json
//	@Param			userID	header		string				true	"user_id"
//	@Success		200		{object}	model.NudgeDetail	"Fetched user nudges list"
//	@Failure		400		{string}	string				"Bad request"
//	@Failure		500		{string}	string				"Internal Server Error"
//	@Router			/user/nudges [get]
func GetUserNudges(c *gin.Context) {
	userID := c.Request.Header.Get("user_id")
	id, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userProfileService := Service.NewUserProfileService()
	userNudges, err := userProfileService.GetUserNudgesService(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "couldn't find user nudges", "error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"data": userNudges})
}

// UpdateMediaHandler godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Update Media Handler
//	@Description	API to update user media order id
//	@Tags			Profile
//	@Accept			json
//	@Produce		json
//	@Param			mediaDetails	body		model.MediaOrderId	true	"mediaDetails"
//	@Param			userID			header		string				true	"user_id"
//	@Success		200				{string}	string				"user media profile updated successfully"
//	@Failure		400				{string}	string				"Bad request"
//	@Failure		500				{string}	string				"Internal Server Error"
//	@Router			/user/update/profileMedia [post]
func UpdateMediaHandler(c *gin.Context) {

	userId := c.Request.Header.Get("user_id")
	id, err := strconv.Atoi(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var mediaDetails model.MediaOrderId
	if err := c.ShouldBindJSON(&mediaDetails); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userProfileService := Service.NewUserProfileService()
	userProfile, _ := userProfileService.GetUserProfileFromDB(id)
	if err != nil {
		zapLogger.Logger.Error(fmt.Sprintf("error in fetching user profile for user_id=%d", id))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error in fetching user profile", "error": err.Error()})
		return
	}

	_, err = userProfileService.UpdateMediaProfile(userProfile, mediaDetails)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error in updating profile media"})
	}

	c.JSON(http.StatusOK, gin.H{"message": "user media profile updated successfully"})
}

// CreateUserNudgeWithMedia godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Create user nudges with media
//	@Description	API to create user nudges with media
//	@Tags			Profile
//	@Accept			mpfd
//	@Produce		json
//	@Param			media			formData	file				true	"nudgeMedia"
//	@Param			nudgeRequest	body		model.NudgeRequest	true	"nudgeRequest"
//	@Param			userID			header		string				true	"user_id"
//	@Success		200				{object}	model.NudgeDetail	"user nudge created successfully"
//	@Failure		400				{string}	string				"Bad request"
//	@Failure		500				{string}	string				"Internal Server Error"
//	@Router			/user/nudge/media [post]
func CreateUserNudgeWithMedia(c *gin.Context) {
	userID := c.Request.Header.Get("user_id")
	id, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	form, _ := c.MultipartForm()
	files := form.File["media"]

	var nudgeRequest model.NudgeRequest
	if err := c.ShouldBind(&nudgeRequest); err != nil {
		c.String(http.StatusBadRequest, "Invalid form data: %v", err)
		return
	}

	nudgeDetails := model.NudgeDetail{
		Question: nudgeRequest.Question,
		Answer:   nudgeRequest.Answer,
		Order:    nudgeRequest.Order,
		Type:     nudgeRequest.Type,
	}
	userProfileService := Service.NewUserProfileService()

	if len(files) != 0 {
		updatedNudgeDetails, err := userProfileService.UploadNudgeMediaToS3(userID, files[0], nudgeDetails)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "error in uploading media file", "error": err.Error()})
			return
		}
		nudgeDetails = updatedNudgeDetails
	} else {
		nudgeDetails.Type = "text"
	}

	userNudge, err := userProfileService.CreateUserNudgeService(nudgeDetails, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error in creating user nudge", "error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "user nudge created successfully", "data": userNudge})
}

// CreateProfileIndex godoc
//
//	@Security		ApiKeyAuth
//	@Summary		CreateProfileIndex
//	@Description	Create User Profile Index
//	@Tags			Profile
//	@Accept			json
//	@Produce		json
//	@Success		200							{string}	string	"index created successfully"
//	@Failure		400							{string}	string	Bad	request
//	@Failure		500							{string}	string	"internal server error"
//	@Router			/create/user_profile/index	[POST]
func CreateProfileIndex(c *gin.Context) {
	userProfileService := Service.NewUserProfileService()
	err := userProfileService.CreateProfileIndex()
	if err != nil {
		c.JSON(500, gin.H{
			"message": "error in creating index in elasticSearch",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{"message": "index created successfully"})
}

// UpdateUserNudge godoc
//
//	@Security		ApiKeyAuth
//	@Summary		UpdateUserNudge
//	@Description	Update user nudges
//	@Tags			Profile
//	@Accept			mpfd
//	@Produce		json
//	@Param			media			formData	file				true	"nudgeMedia"
//	@Param			nudgeRequest	body		model.NudgeRequest	true	"nudgeRequest"
//	@Param			userID			header		string				true	"user_id"
//	@Param			nudgeID			header		string				true	"nudge_id"
//	@Success		200				{object}	model.NudgeDetail	"user nudge updated successfully"
//	@Failure		400				{string}	string				"Bad request"
//	@Failure		500				{string}	string				"Internal Server Error"
//	@Router			/user/nudge		[PUT]
func UpdateUserNudge(c *gin.Context) {
	userID := c.Request.Header.Get("user_id")

	id := c.Request.Header.Get("nudge_id")
	nudgeID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	form, _ := c.MultipartForm()
	files := form.File["media"]

	var nudgeRequest model.NudgeRequest
	if err := c.ShouldBind(&nudgeRequest); err != nil {
		c.String(http.StatusBadRequest, "Invalid form data: %v", err)
		return
	}

	nudgeDetails := model.NudgeDetail{
		Question: nudgeRequest.Question,
		Answer:   nudgeRequest.Answer,
		Order:    nudgeRequest.Order,
		Type:     nudgeRequest.Type,
	}

	userProfileService := Service.NewUserProfileService()

	if len(files) != 0 {
		updatedNudgeDetails, err := userProfileService.UploadNudgeMediaToS3(userID, files[0], nudgeDetails)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "error in uploading media file", "error": err.Error()})
			return
		}
		nudgeDetails = updatedNudgeDetails
	} else {
		nudgeDetails.Type = "text"
	}

	userNudge, err := userProfileService.UpdateUserNudge(nudgeDetails, nudgeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error in updating user nudge", "error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "user nudge updated successfully", "data": userNudge})

}

// DeleteUserNudge godoc
//
//	@Security		ApiKeyAuth
//	@Summary		DeleteUserNudge
//	@Description	Delete user nudges
//	@Tags			Profile
//	@Accept			json
//	@Produce		json
//	@Param			nudge_id	header		int		true	"Nudge ID"
//	@Success		200			{string}	string	"user nudge deleted successfully"
//	@Failure		400			{string}	string	Bad	request
//	@Failure		500			{string}	string	"internal server error"
//	@Router			/user/nudge	[DELETE]
func DeleteUserNudge(c *gin.Context) {
	id := c.Request.Header.Get("nudge_id")
	nudgeID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userProfileService := Service.NewUserProfileService()
	err = userProfileService.DeleteUserNudge(nudgeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error in deleting user nudge", "error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "user nudge deleted successfully"})
}

// UpdateLocationHandler godoc
//
//	@Security		ApiKeyAuth
//	@Summary		UpdateLocationHandler
//	@Description	Update user location
//	@Tags			Profile
//	@Accept			json
//	@Produce		json
//	@Param			location		body		model.UserLocation	true	"location"
//	@Success		200				{string}	string				"user location updated successfully"
//	@Failure		400				{string}	string				Bad	request
//	@Failure		500				{string}	string				"internal server error"
//	@Router			/user/location	[PUT]
func UpdateLocationHandler(c *gin.Context) {
	id := c.Request.Header.Get("user_id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var location model.UserLocation
	if err := c.ShouldBindJSON(&location); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userProfileService := Service.NewUserProfileService()
	userProfile, err := userProfileService.GetUserProfileFromDB(userId)
	err = userProfileService.UpdateUserLocation(userProfile, location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error in updating user location", "error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "user location updated successfully"})
}

// GetFiltersHandler godoc
//
//	@Security		None
//	@Summary		GetFiltersHandler
//	@Description	Get filters list
//	@Tags			Profile
//	@Accept			json
//	@Produce		json
//	@Success		200			{string}	string	"filters list"
//	@Failure		400			{string}	string	Bad	request
//	@Failure		500			{string}	string	"internal server error"
//	@Router			/filters	[GET]
func GetFiltersHandler(c *gin.Context) {
	userProfileService := Service.NewUserProfileService()
	filtersList, err := userProfileService.GetFiltersList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error in getting filters list", "error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"data": filtersList})
}

// UpdateUserInterests godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Update User Interests_category
//	@Description	API to update user interests
//	@Tags			Profile
//	@Accept			json
//	@Produce		json
//	@Param			interests	body		model.InterestData	true	"userInterestsData"
//	@Param			userID		header		string				true	"user_id"
//	@Success		200			{string}	string				"User interests updated successfully"
//	@Failure		400			{string}	string				"Bad request"
//	@Failure		500			{string}	string				"Internal Server Error"
//	@Router			/user/interests [put]
func UpdateUserInterests(c *gin.Context) {
	userID := c.Request.Header.Get("user_id")
	id, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var interests model.InterestData
	if err := c.ShouldBindJSON(&interests); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userProfileService := Service.NewUserProfileService()
	userInterests, err := userProfileService.UpdateUserInterest(interests.InterestDetails, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error in updating user interests", "error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "user interests are updated successfully", "data": userInterests})

}

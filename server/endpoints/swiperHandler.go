package endpoints

import (
	"github.com/SuperMatch/model/dto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/SuperMatch/service"
)

// SwipeHandler godoc
//
//	@Security		ApiKeyAuth
//
//	@Summary		swipe
//	@Description	swipe
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			userLike	body		dto.UserLikeDTO	true	"userLike"
//	@Param			user_id		header		string			true	"user_id"
//	@Success		200			{string}	string			"success"
//	@Failure		400			{string}	string			"Bad request"
//	@Failure		500			{string}	string			"Internal Server Error"
//	@Router			/user/swipe [post]
func SwipeHandler(c *gin.Context) {

	var userLike dto.UserLikeDTO
	if err := c.ShouldBindJSON(&userLike); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// swipe service
	swipeService := service.NewSwipeService()
	isMatch, err := swipeService.Swipe(userLike)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	response := map[string]string{"isMatch": strconv.FormatBool(isMatch)}
	c.JSON(200, gin.H{"message": "success", "data": response})
}

// GetUserMatchHandler godoc
//
//	@Security		ApiKeyAuth
//	@Summary		GetUserMatch
//	@Description	Get User Match
//	@Tags			user
//	@Produce		json
//	@Param			user_id	header		string				true	"user_"
//	@Success		200		{array}		dto.UserMatchDTO	"successfully received user match"
//	@Failure		400		{string}	string				"Bad request"
//	@Failure		500		{string}	string				"error in fetching user match."
//	@Router			/user/match [get]
func GetUserMatchHandler(c *gin.Context) {
	userID := c.Request.Header.Get("user_id")
	id, _ := strconv.Atoi(userID)
	swipeService := service.NewSwipeService()
	data, err := swipeService.GetUserMatchListFromDB(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error in fetching user match.", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "successfully received user match", "data": data})
}

// GetUserLikesHandler godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Get user likes
//	@Description	Get User Likes
//	@Tags			user
//	@Produce		json
//	@Param			user_id	header		string				true	"user_"
//	@Success		200		{array}		model.UserLikers	"successfully received user likes"
//	@Failure		400		{string}	string				"Bad request"
//	@Failure		500		{string}	string				"error in fetching user match."
//	@Router			/user/likes [get]
func GetUserLikesHandler(c *gin.Context) {
	userID := c.Request.Header.Get("user_id")
	id, _ := strconv.Atoi(userID)
	swipeService := service.NewSwipeService()
	data, err := swipeService.GetUserLikes(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error in fetching user likes.", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "successfully received user likes", "data": data})
}

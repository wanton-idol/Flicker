package tests

import (
	"math/rand"
	"strings"
	"testing"

	"github.com/SuperMatch/model"
	"github.com/SuperMatch/model/dto"
	"github.com/SuperMatch/pkg/db/dao"
	mockdao "github.com/SuperMatch/pkg/db/dao/mocks"
	mockredis "github.com/SuperMatch/pkg/redis/mocks"
	"github.com/SuperMatch/service/mocks"
	"github.com/SuperMatch/utilities"
	"github.com/golang/mock/gomock"
)

func TestCheckForExistingResponse(t *testing.T) {
	dislikerUserId := 1
	dislikedUserId := 2
	key := utilities.ConvertIntToString(dislikerUserId) + ":" + utilities.ConvertIntToString(dislikedUserId)

	ctrl := gomock.NewController(t)
	ctrl.Finish()

	mockCache := mockredis.NewMockLikeDislikeCacheInterface(ctrl)
	mockCache.EXPECT().GetLikeDislike(gomock.Eq(key)).
		Return(nil, nil)
}

func TestLike(t *testing.T) {
	userId := 1
	likedUserId := 2
	key := utilities.ConvertIntToString(userId) + ":" + utilities.ConvertIntToString(likedUserId)

	ctrl := gomock.NewController(t)
	ctrl.Finish()

	mockCache := mockredis.NewMockLikeDislikeCacheInterface(ctrl)
	mockCache.EXPECT().PutLikeDislike(gomock.Eq(key), gomock.Eq(utilities.ConvertIntToString(1))).
		Return(nil)

}

func TestDislike(t *testing.T) {
	userId := 1
	dislikedUserId := 2
	key := utilities.ConvertIntToString(userId) + ":" + utilities.ConvertIntToString(dislikedUserId)

	ctrl := gomock.NewController(t)
	ctrl.Finish()

	mockCache := mockredis.NewMockLikeDislikeCacheInterface(ctrl)
	mockCache.EXPECT().PutLikeDislike(gomock.Eq(key), gomock.Eq(utilities.ConvertIntToString(0))).
		Return(nil)

}

func TestAddToUserMatchList(t *testing.T) {
	userLike := dto.UserLikeDTO{
		LikeeID: 1,
		LikerID: 2,
		Type:    1,
	}

	match1 := model.UserMatch{
		UserID:     userLike.LikerID,
		MatchID:    userLike.LikeeID,
		Match_type: 1,
	}

	match2 := model.UserMatch{
		UserID:     userLike.LikeeID,
		MatchID:    userLike.LikerID,
		Match_type: 1,
	}

	ctrl := gomock.NewController(t)
	ctrl.Finish()

	mockUserMatch := mockdao.NewMockUserMatchDao(ctrl)
	mockUserMatch.EXPECT().Insert(gomock.Any(), gomock.Eq(match1)).
		Return(match1, nil)

	mockUserMatch.EXPECT().Insert(gomock.Any(), gomock.Eq(match2)).
		Return(match2, nil)

	mockCache := mockredis.NewMockLikeDislikeCacheInterface(ctrl)
	mockCache.EXPECT().AddToUserMatchList(gomock.Eq(utilities.ConvertIntToString(userLike.LikerID)), gomock.Eq(utilities.ConvertIntToString(userLike.LikeeID))).
		Return(nil)

	mockCache.EXPECT().AddToUserMatchList(gomock.Eq(utilities.ConvertIntToString(userLike.LikeeID)), gomock.Eq(utilities.ConvertIntToString(userLike.LikerID))).
		Return(nil)

}

func TestGetUserMatchListFromCache(t *testing.T) {
	userID := 1
	key := utilities.ConvertIntToString(userID)

	ctrl := gomock.NewController(t)
	ctrl.Finish()

	mockCache := mockredis.NewMockLikeDislikeCacheInterface(ctrl)
	mockCache.EXPECT().GetMatchList(gomock.Eq(key)).
		Return(nil, nil)
}

func TestGetUserMatchListFromDB(t *testing.T) {
	userID := 1
	userMatch := []model.UserMatch{
		{
			ID:      1,
			UserID:  userID,
			MatchID: 2,
		},
		{
			ID:      2,
			UserID:  userID,
			MatchID: 3,
		},
	}
	userMatchUserMedia := []dao.UserMatchUserMediaDTO{
		{
			ID:      1,
			UserId:  userID,
			MatchId: 1,
			URL:     "https://user-profile-supermatch.s3.ap-south-1.amazonaws.com/1/profile/2023-04-24T03%3A38%3A23-images3.webp",
		},
		{
			ID:      2,
			UserId:  userID,
			MatchId: 2,
			URL:     "https://user-profile-supermatch.s3.ap-south-1.amazonaws.com/1/profile/2023-04-24T03%3A38%3A23-images3.webp",
		},
	}
	ctrl := gomock.NewController(t)
	ctrl.Finish()

	mockUserMatch := mockdao.NewMockUserMatchDao(ctrl)
	mockUserMatch.EXPECT().FindByUserId(gomock.Any(), gomock.Eq(userID)).
		Return(userMatch, nil)
	matchIDs := make([]int, 0)
	for _, val := range userMatch {
		matchIDs = append(matchIDs, val.MatchID)
	}

	if len(matchIDs) == 0 {
		t.Error("No match IDs found.")
	}

	mockUserMedia := mockdao.NewMockUserMediaRepository(ctrl)
	mockUserMedia.EXPECT().FindFirstGroupByUserID(gomock.Any(), gomock.Eq(matchIDs)).
		Return(userMatchUserMedia, nil)

	userMatches := make([]dto.UserMatchDTO, 0)
	for _, match := range userMatchUserMedia {
		key := strings.ReplaceAll(strings.TrimPrefix(match.URL, S3BucketPath), "%3A", ":")
		signedURL := ""
		mockS3Service := mocks.NewMockS3ServiceInterface(ctrl)
		mockS3Service.EXPECT().SignS3FilesUrl(gomock.Eq(userProfileS3Bucket), gomock.Eq(key)).
			Return(signedURL, nil)
		x := dto.UserMatchDTO{
			ID:        match.ID,
			UserId:    match.UserId,
			MatchId:   match.MatchId,
			OrderId:   match.OrderId,
			MediaId:   match.MediaId,
			URL:       signedURL,
			ChatId:    match.ChatId,
			CreatedAt: match.CreatedAt,
			DeletedAt: match.DeletedAt,
		}
		userMatches = append(userMatches, x)
	}

}

func TestRemoveMatch(t *testing.T) {
	userLike := dto.UserLikeDTO{
		LikeeID: 1,
		LikerID: 2,
		Type:    1,
	}

	ctrl := gomock.NewController(t)
	ctrl.Finish()

	mockCache := mockredis.NewMockLikeDislikeCacheInterface(ctrl)
	key := utilities.ConvertIntToString(userLike.LikerID) + ":" + utilities.ConvertIntToString(userLike.LikeeID)
	mockCache.EXPECT().RemoveFromUserMatchList(gomock.Eq(key)).Return(nil)

	key = utilities.ConvertIntToString(userLike.LikeeID) + ":" + utilities.ConvertIntToString(userLike.LikerID)
	mockCache.EXPECT().RemoveFromUserMatchList(gomock.Eq(key)).Return(nil)

	mockUserMatch := mockdao.NewMockUserMatchDao(ctrl)
	mockUserMatch.EXPECT().DeleteByUserID(gomock.Any(), gomock.Eq(userLike.LikerID), gomock.Eq(userLike.LikeeID)).
		Return(nil)
}

// check-This
func TestSwipe(t *testing.T) {
	userLike := dto.UserLikeDTO{
		LikeeID: 1,
		LikerID: 2,
		Type:    1,
	}
	response := rand.Int()%2 == 0
	TestCheckForExistingResponse(t)

	if userLike.Type == 1 {
		if !response {
			TestLike(t)
		} else {
			TestAddToUserMatchList(t)
			TestLike(t)
		}
	} else {
		TestDislike(t)
	}
}

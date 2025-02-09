package service

import (
	"context"
	"fmt"
	"github.com/SuperMatch/zapLogger"
	"go.uber.org/zap"
	"strings"

	"github.com/SuperMatch/model"
	"github.com/SuperMatch/model/dto"
	"github.com/SuperMatch/pkg/db/dao"
	"github.com/SuperMatch/pkg/redis"
	"github.com/SuperMatch/utilities"
)

const TIME_FORMAT = "2006-01-02T15:04:05"

type SwipeServiceInterface interface {
	Swipe(userActionDTO dto.UserLikeDTO) (bool, error)
	Like(userId int, likedUserId int) error
	Dislike(userId int, dislikedUserId int) error
	checkForExistingResponse(userId, likeeID int) (bool, error)
	addToUserMatchList(userActionDTO dto.UserLikeDTO) error
	GetUserMatchListFromCache(userID int) ([]int, error)
	GetUserMatchListFromDB(userID int) ([]dto.UserMatchDTO, error)
	RemoveMatch(userActionDTO dto.UserLikeDTO) error
	GetUserLikes(userID int) ([]model.UserLikers, error)
	PutLiker(likerID, likeeID int) error
	RemoveLikerFromLikeeList(likerID, likeeID int) error
}

type SwipeService struct {
	LikeDislikeCache    redis.LikeDislikeCacheInterface
	UserMatchDao        dao.UserMatchDao
	UserMediaRepository dao.UserMediaRepository
	S3Service           S3ServiceInterface
}

func NewSwipeService() *SwipeService {
	return &SwipeService{
		LikeDislikeCache:    redis.LikeDislikeCacheConstructor(),
		UserMatchDao:        dao.NewUserMatchDaoImpl(),
		UserMediaRepository: dao.NewUserMediaRepository(),
		S3Service:           NewS3Service(),
	}
}

func (s *SwipeService) Swipe(userActionDTO dto.UserLikeDTO) (bool, error) {

	isLiked, err := s.checkAlreadyLiked(userActionDTO.LikerID, userActionDTO.LikeeID)
	if err != nil {
		zapLogger.Logger.Error("error in checking liker already liked the likee:", zap.Error(err))
		return false, err
	}
	if isLiked && userActionDTO.Type == 1 {
		zapLogger.Logger.Error("liker already liked the likee")
		return false, fmt.Errorf("liker: %d already liked the likee: %d", userActionDTO.LikerID, userActionDTO.LikeeID)
	}

	response, err := s.checkForExistingResponse(userActionDTO.LikeeID, userActionDTO.LikerID)
	zapLogger.Logger.Info("swipe response:", zap.Any("response", response))
	if err != nil {
		zapLogger.Logger.Error("error in checking likee already liked the liker:", zap.Error(err))
		return false, err
	}

	if userActionDTO.Type == 1 {
		//if already disiked by other user,just don't match
		if !response {
			zapLogger.Logger.Debug("already not liked")
			err := s.Like(userActionDTO.LikerID, userActionDTO.LikeeID)
			if err != nil {
				zapLogger.Logger.Error("error in liking:", zap.Error(err))
				return false, err
			}
			err = s.PutLiker(userActionDTO.LikerID, userActionDTO.LikeeID)
			if err != nil {
				zapLogger.Logger.Error("error in putting liker:", zap.Error(err))
				return false, err
			}

		} else {
			zapLogger.Logger.Debug("already liked")
			//match if already liked by other user & create chat room
			err := s.addToUserMatchList(userActionDTO)
			if err != nil {
				zapLogger.Logger.Error("error in adding to user match list:", zap.Error(err))
				return false, err
			}
			err = s.Like(userActionDTO.LikerID, userActionDTO.LikeeID)
			if err != nil {
				zapLogger.Logger.Error("error in liking:", zap.Error(err))
				return false, err
			}

			err = s.RemoveLikerFromLikeeList(userActionDTO.LikerID, userActionDTO.LikeeID)
			if err != nil {
				zapLogger.Logger.Error("error in removing liker from likee list:", zap.Error(err))
				return false, err
			}

			return true, nil
		}
	} else {
		zapLogger.Logger.Debug("Dislike")
		err := s.Dislike(userActionDTO.LikerID, userActionDTO.LikeeID)
		if err != nil {
			zapLogger.Logger.Error("error in disliking:", zap.Error(err))
			return false, err
		}
		if response {
			err = s.RemoveLikerFromLikeeList(userActionDTO.LikerID, userActionDTO.LikeeID)
			if err != nil {
				zapLogger.Logger.Error("error in removing liker from likee list:", zap.Error(err))
				return false, err
			}
		}

	}
	return false, nil
}

func (s *SwipeService) Like(userId int, likedUserId int) error {
	key := utilities.ConvertIntToString(userId) + ":" + utilities.ConvertIntToString(likedUserId)
	err := s.LikeDislikeCache.PutLikeDislike(key, utilities.ConvertIntToString(1))
	if err != nil {
		zapLogger.Logger.Error("error in putting like", zap.Error(err))
		return err
	}

	return nil
}

func (s *SwipeService) PutLiker(likerID, likeeID int) error {
	err := s.LikeDislikeCache.PutLiker(utilities.ConvertIntToString(likeeID), utilities.ConvertIntToString(likerID))
	if err != nil {
		zapLogger.Logger.Error("error in putting likerID into the likee's likers list", zap.Error(err))
		return err
	}

	return nil
}

func (s *SwipeService) RemoveLikerFromLikeeList(likerID, likeeID int) error {
	err := s.LikeDislikeCache.RemoveLikerFromLikeeList(utilities.ConvertIntToString(likerID), utilities.ConvertIntToString(likeeID))
	if err != nil {
		zapLogger.Logger.Error("error in removing likerID from the likee's likers list", zap.Error(err))
		return err
	}

	return nil
}

func (s *SwipeService) Dislike(userId int, dislikedUserId int) error {
	key := utilities.ConvertIntToString(userId) + ":" + utilities.ConvertIntToString(dislikedUserId)
	err := s.LikeDislikeCache.PutLikeDislike(key, utilities.ConvertIntToString(0))
	return err
}

func (s *SwipeService) checkForExistingResponse(dislikerUserId, dislikedUserId int) (bool, error) {
	key := utilities.ConvertIntToString(dislikerUserId) + ":" + utilities.ConvertIntToString(dislikedUserId)
	val, err := s.LikeDislikeCache.GetLikeDislike(key)

	if err != nil {
		zapLogger.Logger.Error("Error in checking existing like dislike", zap.Error(err))
		return false, err
	}

	if val == nil {
		return false, nil
	} else if *val == "0" {
		return false, nil
	} else {
		return true, nil
	}
}

func (s *SwipeService) addToUserMatchList(userActionDTO dto.UserLikeDTO) error {

	// matches := make([]model.UserMatch, 2)
	chatId := utilities.ConvertStringToStringPointer(fmt.Sprintf("%d_%d", userActionDTO.LikerID, userActionDTO.LikeeID))
	match1 := model.UserMatch{
		UserID:     userActionDTO.LikerID,
		MatchID:    userActionDTO.LikeeID,
		Match_type: 1, // right now it's 1 but depends on matchtype
		ChatID:     chatId,
	}

	match2 := model.UserMatch{
		UserID:     userActionDTO.LikeeID,
		MatchID:    userActionDTO.LikerID,
		Match_type: 1,
		ChatID:     chatId,
	}

	_, _ = s.UserMatchDao.Insert(context.Background(), match1)
	_, err := s.UserMatchDao.Insert(context.Background(), match2)

	if err != nil {
		zapLogger.Logger.Error("Error in adding to user match list", zap.Error(err))
		return err
	}

	return nil
}

func (s *SwipeService) GetUserMatchListFromCache(userID int) ([]int, error) {
	key := utilities.ConvertIntToString(userID)
	x, err := s.LikeDislikeCache.GetMatchList(key)
	return x, err
}

func (s *SwipeService) GetUserMatchListFromDB(userID int) ([]dto.UserMatchDTO, error) {

	matchList, err := s.UserMatchDao.FindByUserId(context.Background(), userID)

	userMatches := make([]dto.UserMatchDTO, 0)

	if err != nil {
		zapLogger.Logger.Error("Error in getting match list from db", zap.Error(err))
		return userMatches, err
	}

	matchIdS := make([]int, 0)

	for _, val := range matchList {
		matchIdS = append(matchIdS, val.MatchID)
	}

	if len(matchIdS) == 0 {
		return userMatches, nil
	}

	matchListMedia, err := s.UserMediaRepository.FindFirstGroupByUserID(context.Background(), matchIdS)

	if err != nil {
		zapLogger.Logger.Error("Error in getting match list from db", zap.Error(err))
		return nil, err
	}

	for _, match := range matchListMedia {
		key := strings.ReplaceAll(strings.TrimPrefix(match.URL, S3_BUCKET_PATH), "%3A", ":")
		signedURL, err := s.S3Service.SignS3FilesUrl(user_profile_S3_bucket, key)
		if err != nil {
			zapLogger.Logger.Error("Error in getting signed url", zap.Error(err))
			return nil, err
		}
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

	return userMatches, nil
}

func (s *SwipeService) RemoveMatch(userActionDTO dto.UserLikeDTO) error {

	err := s.LikeDislikeCache.RemoveFromUserMatchList(utilities.ConvertIntToString(userActionDTO.LikerID) + ":" + utilities.ConvertIntToString(userActionDTO.LikeeID))

	if err != nil {
		zapLogger.Logger.Error("Error in removing from liker match list", zap.Error(err))
		return err
	}

	err = s.LikeDislikeCache.RemoveFromUserMatchList(utilities.ConvertIntToString(userActionDTO.LikeeID) + ":" + utilities.ConvertIntToString(userActionDTO.LikerID))
	if err != nil {
		zapLogger.Logger.Error("Error in removing from likee match list", zap.Error(err))
		return err
	}

	err = s.UserMatchDao.DeleteByUserID(context.Background(), userActionDTO.LikerID, userActionDTO.LikeeID)

	if err != nil {
		zapLogger.Logger.Error("Error in removing from db", zap.Error(err))
		return err
	}

	return nil
}

func (s *SwipeService) checkAlreadyLiked(likerID, likeeID int) (bool, error) {
	key := utilities.ConvertIntToString(likerID) + ":" + utilities.ConvertIntToString(likeeID)
	val, err := s.LikeDislikeCache.GetLikeDislike(key)
	if err != nil {
		zapLogger.Logger.Error("Error in checking existing like dislike", zap.Error(err))
		return false, err
	}

	if val == nil {
		return false, nil
	} else if *val == "0" {
		return false, nil
	} else {
		return true, nil
	}
}

func (s *SwipeService) GetUserLikes(userID int) ([]model.UserLikers, error) {
	key := utilities.ConvertIntToString(userID)
	likes, err := s.LikeDislikeCache.GetUserLikes(key)
	if err != nil {
		zapLogger.Logger.Error("Error in getting user likes", zap.Error(err))
		return nil, err
	}

	userMedia, err := s.UserMediaRepository.FindByUserIDs(likes)
	if err != nil {
		zapLogger.Logger.Error("Error in getting user media", zap.Error(err))
		return nil, err
	}

	userLikers := make([]model.UserLikers, 0)
	for _, media := range userMedia {
		key := strings.ReplaceAll(strings.TrimPrefix(media.URL, S3_BUCKET_PATH), "%3A", ":")
		signedURL, err := s.S3Service.SignS3FilesUrl(user_profile_S3_bucket, key)
		if err != nil {
			zapLogger.Logger.Error("Error in getting signed url", zap.Error(err))
			return nil, err
		}

		userLikers = append(userLikers, model.UserLikers{
			UserID: media.UserId,
			Image:  signedURL,
		})

	}

	return userLikers, err
}

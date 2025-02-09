package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/twpayne/go-geom"
	"gorm.io/gorm"

	"github.com/SuperMatch/zapLogger"
	"go.uber.org/zap"

	"github.com/SuperMatch/model"
	"github.com/SuperMatch/model/dto"
	elasticsearchPkg "github.com/SuperMatch/model/elasticSearch"
	"github.com/SuperMatch/pkg/db/dao"
	pkg "github.com/SuperMatch/pkg/elasticSeach"
	_ "github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/wkb"
)

const (
	ISO_DATE_FORMAT        = "2006-01-02"
	user_profile_S3_bucket = "user-profile-supermatch"
	S3_BUCKET_PATH         = "https://" + user_profile_S3_bucket + ".s3.ap-south-1.amazonaws.com"
)

type UserProfileInterface interface {
	GetUserProfile(userProfileId int) (elasticsearchPkg.UserProfile, error)
	GetUserProfileFromDB(userId int) (model.UserProfile, error)
	UpdatePremiumProfile(profile model.UserProfile) (model.UserProfile, error)
	UpdateUserProfileFirst(userProfileES elasticsearchPkg.UserProfile, userprofileDB model.UserProfile) (elasticsearchPkg.UserProfile, error)

	CreateUserProfile(profile model.UserProfile, userProfileDTO dto.UserProfile) (dto.UserProfile, error)
	UpdateUserProfile(profile model.UserProfile, userProfileDTO dto.UserProfile) (dto.UserProfile, error)
	SearchProfile(user elasticsearchPkg.UserProfile) ([]elasticsearchPkg.UserProfile, error)
	SaveProfileMedia(user model.UserProfile, urls string, mediaRequestDTO dto.ProfileMediaRequestDTO) (model.UserMedia, error)
	GetProfileMediaByUserId(userID int) ([]model.UserMedia, error)
	RemoveProfileMedia(userID int, ImageId int) error
	CreateUserProfileES(user dto.UserProfile, userProfile elasticsearchPkg.UserProfile) (elasticsearchPkg.UserProfile, error)
	GetInterestsCategory() (model.InterestsListResponse, error)
	CreateUserInterest(interests model.InterestsCategory, userID int) ([]model.UserInterests, error)
	GetUserInterests(userID int) (model.InterestsCategory, error)
	CheckAllowedFileType(extension string) bool
	CheckAllowedAudioFileType(extension string) bool
	UpdateMediaProfile(user model.UserProfile, mediaDetails model.MediaOrderId) (model.UserMedia, error)
	UploadNudgeMediaToS3(userId string, file *multipart.FileHeader, nudgeDetails model.NudgeDetail) (model.NudgeDetail, error)
	CreateProfileIndex() error
	DeleteUserNudge(nudgeID int) error
	GetFiltersList() (model.Filters, error)
	UpdateUserInterest(interests []model.InterestDetails, userID int) ([]model.UserInterests, error)
	CheckAllowedVideoFileType(extension string) bool
}

type UserProfileService struct {
	esIndex              pkg.ElasticSearchIndexer
	userMedia            dao.UserMediaRepository
	s3Service            S3ServiceInterface
	swipeService         SwipeServiceInterface
	userProfileDao       dao.UserProfileRepository
	userSearchProfileDao dao.UserSearchProfileRepository
	advancedFilerDao     dao.AdvancedFilterRepository
	interestsDao         dao.InterestsDao
	userNudgesDao        dao.UserNudgesDao
	filtersDao           dao.FiltersDao
}

func NewUserProfileService() *UserProfileService {
	return &UserProfileService{
		esIndex:              pkg.NewElasticSearchIndexerImpl(),
		userMedia:            dao.NewUserMediaRepository(),
		s3Service:            NewS3Service(),
		swipeService:         NewSwipeService(),
		userProfileDao:       dao.NewUserProfileRepository(),
		userSearchProfileDao: dao.NewUserSearchProfile(),
		advancedFilerDao:     dao.NewAdvancedFilterRepository(),
		interestsDao:         dao.NewInterestsDaoImpl(),
		userNudgesDao:        dao.NewUserNudgesDaoImpl(),
		filtersDao:           dao.NewFiltersDaoImpl(),
	}
}

func (u *UserProfileService) GetUserProfile(userProfileId int) (elasticsearchPkg.UserProfile, error) {

	userProfile, err := u.esIndex.GetUserProfile(userProfileId)
	if err != nil {
		zapLogger.Logger.Error("Error while getting user profile: ", zap.Error(err))
		return elasticsearchPkg.UserProfile{}, err
	}
	return userProfile, nil
}

func (u *UserProfileService) GetUserProfileFromDB(userId int) (model.UserProfile, error) {
	userProfile, err := u.userProfileDao.FindByUserId(context.Background(), userId)
	if err != nil {
		zapLogger.Logger.Error("Error while getting user profile: ", zap.Error(err))
		return model.UserProfile{}, err
	}
	return userProfile, nil
}

func (u *UserProfileService) UpdatePremiumProfile(userProfile model.UserProfile) (model.UserProfile, error) {
	userProfile.IsPremium = true
	return u.userProfileDao.UpdateUserProfile(context.Background(), userProfile)
}

func (u *UserProfileService) UpdateUserProfileFirst(userProfileES elasticsearchPkg.UserProfile, userprofileDB model.UserProfile) (elasticsearchPkg.UserProfile, error) {

	updatedProfile, err := u.userProfileDao.UpdateUserProfile(context.Background(), userprofileDB)
	if err != nil {
		zapLogger.Logger.Error("error while creating user profile in DB: ", zap.Error(err))
		return userProfileES, err
	}

	//update user profile
	userProfileES.Id = updatedProfile.UserId

	jsonString, err := json.Marshal(userProfileES)

	zapLogger.Logger.Debug(string(jsonString))
	if err != nil {
		return userProfileES, err
	}

	err = u.esIndex.IndexUserProfile(userProfileES, jsonString)

	if err != nil {
		zapLogger.Logger.Error("Error while indexing user profile: ", zap.Error(err))
		return userProfileES, err
	}

	return userProfileES, nil
}

func (u *UserProfileService) CreateUserProfile(profile model.UserProfile, userProfileDTO dto.UserProfile) (dto.UserProfile, error) {

	profile, err := u.userProfileDao.CreateUserProfile(context.Background(), profile)
	if err != nil {
		zapLogger.Logger.Error("error while creating user profile in DB:", zap.Error(err))
		return userProfileDTO, err
	}

	//create User Search Profile
	mp := make(map[string]interface{})
	mp["user_profile_id"] = int(profile.ID)
	mp["user_id"] = profile.UserId

	userProfileDTO.UserId = profile.UserId
	userProfileDTO.Id = int(profile.ID)

	_, err = u.userSearchProfileDao.CreateUserSearchProfile(context.Background(), mp)
	if err != nil {
		zapLogger.Logger.Error("error while creating user search profile in DB:", zap.Error(err))
		return dto.UserProfile{}, err
	}

	advancedFilter := model.AdvancedFilter{
		UserId:        profile.UserId,
		UserProfileId: int(profile.ID),
	}
	err = u.advancedFilerDao.CreateAdvancedFilter(context.Background(), advancedFilter)

	if err != nil {
		zapLogger.Logger.Error("error while creating advanced filter in DB:", zap.Error(err))
		return dto.UserProfile{}, err
	}

	userProfileES := elasticsearchPkg.UserProfile{
		Id:     int(profile.ID),
		UserId: profile.UserId,
		UserSearchProfile: elasticsearchPkg.UserSearchProfile{
			Gender:         []elasticsearchPkg.Gender{elasticsearchPkg.MALE, elasticsearchPkg.FEMALE, elasticsearchPkg.BINARY},
			AdvancedFilter: elasticsearchPkg.AdvancedFilter{},
		},
	}

	jsonString, err := json.Marshal(userProfileES)

	fmt.Println(string(jsonString))
	zapLogger.Logger.Debug(string(jsonString))

	if err != nil {
		return userProfileDTO, err
	}

	err = u.esIndex.IndexUserProfile(userProfileES, jsonString)
	if err != nil {
		zapLogger.Logger.Error("error while indexing user profile: ", zap.Error(err))
		return userProfileDTO, err
	}
	return userProfileDTO, nil
}

func (u *UserProfileService) UpdateUserProfile(profile model.UserProfile, userProfileDTO dto.UserProfile) (dto.UserProfile, error) {

	// Map of interface of user profile
	profileMap := make(map[string]interface{})
	// to be set here
	profileMap["ID"] = int(profile.ID)
	profileMap["user_id"] = profile.UserId
	profileMap["first_name"] = userProfileDTO.FirstName
	profileMap["last_name"] = userProfileDTO.LastName
	profileMap["date_of_birth"] = userProfileDTO.DateOfBirth
	profileMap["gender"] = userProfileDTO.Gender
	profileMap["sexual_orientation"] = userProfileDTO.SexualOrientation
	profileMap["latitude"] = userProfileDTO.Latitude
	profileMap["longitude"] = userProfileDTO.Longitude
	if userProfileDTO.Education != nil && userProfileDTO.Education.Level != nil {
		profileMap["education_level"] = userProfileDTO.Education.Level
	}
	if userProfileDTO.Education != nil && userProfileDTO.Education.College != nil {
		profileMap["college"] = userProfileDTO.Education.College
	}
	profileMap["occupation"] = userProfileDTO.Occupation
	profileMap["marital_status"] = userProfileDTO.MaritalStatus
	profileMap["religion"] = userProfileDTO.Religion
	profileMap["height"] = userProfileDTO.Height
	profileMap["weight"] = userProfileDTO.Weight
	profileMap["looking_for"] = userProfileDTO.LookingFor
	profileMap["exercise"] = userProfileDTO.Exercise
	profileMap["drink"] = userProfileDTO.Drink
	profileMap["smoke"] = userProfileDTO.Smoke
	profileMap["about"] = userProfileDTO.About
	profileMap["pronoun"] = userProfileDTO.Pronoun

	var err error
	_, err = u.userProfileDao.UpdateProfileByMap(context.Background(), profileMap)

	if err != nil {
		zapLogger.Logger.Error("Error while creating user profile: ", zap.Error(err))
		return userProfileDTO, err
	}

	userProfileDTO.Id = int(profile.ID)

	oldProfileES, err := u.esIndex.GetUserProfile(int(profile.ID))

	//create ElasticSearch UserProfile object
	userProfileES, err := u.CreateUserProfileES(userProfileDTO, oldProfileES)
	if err != nil {
		zapLogger.Logger.Error("Error while creating user profile: ", zap.Error(err))
		return userProfileDTO, err
	}

	if err != nil {
		return dto.UserProfile{}, err
	}

	jsonString, err := json.Marshal(userProfileES)
	zapLogger.Logger.Debug(string(jsonString))
	if err != nil {
		return userProfileDTO, err
	}

	err = u.esIndex.UpdateUserProfile(userProfileES, jsonString)
	if err != nil {
		zapLogger.Logger.Error("error in updating userSearchProfile", zap.Error(err))
		return userProfileDTO, err
	}
	return userProfileDTO, nil
}

func (u *UserProfileService) SearchProfile(user elasticsearchPkg.UserProfile) ([]elasticsearchPkg.UserProfile, error) {

	pagination := model.Pagination{
		TotalCount: 20,
		PageSize:   20,
		PageNumber: 0,
		Sort:       "asc",
	}

	query := generateQuery(user, &pagination)
	tmp, err := u.esIndex.SearchProfile(query)
	if err != nil {
		zapLogger.Logger.Error("error in updating userSearchProfile", zap.Error(err))
	}
	return tmp, nil
}

func (u *UserProfileService) SaveProfileMedia(user model.UserProfile, url string, requestDTO dto.ProfileMediaRequestDTO) (model.UserMedia, error) {
	userMedia := model.UserMedia{
		UserProfileId: int(user.ID),
		UserId:        user.UserId,
		URL:           url,
		ImageText:     requestDTO.ImageText,
		Latitude:      requestDTO.Latitude,
		Longitude:     requestDTO.Longitude,
		City:          requestDTO.City,
		OrderId:       requestDTO.OrderId,
	}
	userMedia, err := u.userMedia.Insert(context.Background(), userMedia)
	if err != nil {
		zapLogger.Logger.Error("error in saving user media ", zap.Error(err))
		return userMedia, errors.New("error in saving user media")
	}

	key := strings.ReplaceAll(strings.TrimPrefix(userMedia.URL, S3_BUCKET_PATH), "%3A", ":")
	signedURL, err := u.s3Service.SignS3FilesUrl(user_profile_S3_bucket, key)
	userMedia.URL = signedURL

	return userMedia, nil
}

func (u *UserProfileService) GetProfileMediaByUserId(userID int) ([]model.UserMedia, error) {
	data, err := u.userMedia.FindByUserId(context.Background(), userID)
	if err != nil {
		zapLogger.Logger.Error("error in getting profile media from DB.", zap.Error(err))
		return nil, err
	}

	for idx := range data {
		key := strings.ReplaceAll(strings.TrimPrefix(data[idx].URL, S3_BUCKET_PATH), "%3A", ":")
		signedURL, err := u.s3Service.SignS3FilesUrl(user_profile_S3_bucket, key)
		if err != nil {
			zapLogger.Logger.Error("error in getting signed url ", zap.Error(err))
			return nil, err
		}
		zapLogger.Logger.Debug(fmt.Sprintf("normal URL: %v\n", data[idx].URL))
		data[idx].URL = signedURL
		zapLogger.Logger.Debug(fmt.Sprintf("signed URL: %v\n", data[idx].URL))
	}
	return data, nil
}

func (u *UserProfileService) RemoveProfileMedia(userID int, mediaId int) error {

	userMedia, err := u.userMedia.FindByIdAndUserId(context.Background(), mediaId, userID)

	if err != nil {
		zapLogger.Logger.Error("error in getting profile media from DB.", zap.Error(err))
		return err
	}

	if userMedia == nil {
		return nil
	}

	key := strings.ReplaceAll(strings.TrimPrefix(userMedia.URL, S3_BUCKET_PATH), "%3A", ":")
	err = u.s3Service.DeleteFile(user_profile_S3_bucket, key)
	if err != nil {
		zapLogger.Logger.Error("error in deleting file from S3.", zap.Error(err))
		return err
	}

	err = u.userMedia.DeleteById(context.Background(), mediaId, userID)
	if err != nil {
		zapLogger.Logger.Error("error in deleting profile media from DB.", zap.Error(err))
		return errors.New("error in deleting profile media")
	}
	return nil
}

func generateQuery(user elasticsearchPkg.UserProfile, pagination *model.Pagination) map[string]interface{} {

	userSearchProfile := user.UserSearchProfile
	mustMap := []map[string]interface{}{}

	if userSearchProfile.Gender != nil {
		x := map[string]interface{}{
			"terms": map[string]interface{}{
				"gender": userSearchProfile.Gender,
			},
		}
		mustMap = append(mustMap, x)
	}

	if userSearchProfile.MinAge != 0 && userSearchProfile.MaxAge != 0 {

		maxAge := time.Now().AddDate(-userSearchProfile.MaxAge, 0, 0).Format(ISO_DATE_FORMAT)
		minage := time.Now().AddDate(-userSearchProfile.MinAge, 0, 0).Format(ISO_DATE_FORMAT)
		x := map[string]interface{}{
			"range": map[string]interface{}{
				"date_of_birth": map[string]interface{}{
					"gte": maxAge,
					"lte": minage,
				},
			},
		}
		mustMap = append(mustMap, x)
	}

	if userSearchProfile.Distance != 0 {
		x := map[string]interface{}{
			"geo_distance": map[string]interface{}{
				"distance": userSearchProfile.Distance,
				"location": map[string]interface{}{
					"lat": user.Location[1],
					"lon": user.Location[0],
				},
			},
		}
		mustMap = append(mustMap, x)
	}

	if userSearchProfile.Language != nil {
		x := map[string]interface{}{
			"terms": map[string]interface{}{
				"language": userSearchProfile.Language,
			},
		}
		mustMap = append(mustMap, x)
	}

	if user.IsPremium != nil && !*user.IsPremium {

		advancedFilters := user.UserSearchProfile.AdvancedFilter

		if advancedFilters.IsProfileVerified != nil {
			x := map[string]interface{}{
				"terms": map[string]interface{}{
					"advancedFilter.is_profile_verified": advancedFilters.IsProfileVerified,
				},
			}
			mustMap = append(mustMap, x)
		}

		if advancedFilters.IsOnline != nil {
			x := map[string]interface{}{
				"terms": map[string]interface{}{
					"advancedFilter.is_online": advancedFilters.IsOnline,
				},
			}
			mustMap = append(mustMap, x)
		}

		if advancedFilters.Height != nil {
			x := map[string]interface{}{
				"range": map[string]interface{}{
					"advancedFilter.height": map[string]interface{}{
						"gte": advancedFilters.Height[0],
						"lte": advancedFilters.Height[1],
					},
				},
			}
			mustMap = append(mustMap, x)
		}

		if advancedFilters.Exercise != nil {
			x := map[string]interface{}{
				"terms": map[string]interface{}{
					"advancedFilter.exercise": advancedFilters.Exercise,
				},
			}
			mustMap = append(mustMap, x)
		}

		if advancedFilters.Religion != nil {
			x := map[string]interface{}{
				"terms": map[string]interface{}{
					"advancedFilter.religion": advancedFilters.Religion,
				},
			}
			mustMap = append(mustMap, x)
		}

		if advancedFilters.EducationLevel != nil {
			x := map[string]interface{}{
				"terms": map[string]interface{}{
					"advancedFilter.education": advancedFilters.EducationLevel,
				},
			}
			mustMap = append(mustMap, x)
		}

		if advancedFilters.Occupation != nil {
			x := map[string]interface{}{
				"terms": map[string]interface{}{
					"advancedFilter.occupation": advancedFilters.Occupation,
				},
			}
			mustMap = append(mustMap, x)
		}

		if advancedFilters.MaritalStatus != nil {
			x := map[string]interface{}{
				"terms": map[string]interface{}{
					"advancedFilter.marital_status": advancedFilters.MaritalStatus,
				},
			}
			mustMap = append(mustMap, x)
		}

		if advancedFilters.Smoke != nil {
			x := map[string]interface{}{
				"terms": map[string]interface{}{
					"advancedFilter.smoke": advancedFilters.Smoke,
				},
			}
			mustMap = append(mustMap, x)
		}
	}

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": mustMap,
			},
		},
		"sort": []map[string]interface{}{
			{
				"_geo_distance": map[string]interface{}{
					"location": map[string]interface{}{
						"lat": user.Location[1],
						"lon": user.Location[0],
					},
					"order":         "asc",
					"unit":          "km",
					"distance_type": "plane",
				},
			},
		},
	}

	if pagination != nil {
		query["from"] = (pagination.PageNumber)*(pagination.PageSize) + 1
		query["size"] = pagination.PageSize
	}
	return query
}

func (u *UserProfileService) CreateUserProfileES(userProfileDTO dto.UserProfile, userProfile elasticsearchPkg.UserProfile) (elasticsearchPkg.UserProfile, error) {

	userProfileES := elasticsearchPkg.UserProfile{
		Id:                userProfileDTO.Id,
		UserId:            userProfileDTO.UserId,
		FirstName:         userProfileDTO.FirstName,
		LastName:          userProfileDTO.LastName,
		IsVerified:        userProfileDTO.IsVerified,
		IsPremium:         userProfileDTO.IsPremium,
		DateOfBirth:       userProfileDTO.DateOfBirth,
		Gender:            userProfileDTO.Gender,
		SexualOrientation: userProfileDTO.SexualOrientation,
		Location:          []float64{*userProfileDTO.Longitude, *userProfileDTO.Latitude},
		Education:         userProfileDTO.Education,
		Occupation:        userProfileDTO.Occupation,
		MaritalStatus:     userProfileDTO.MaritalStatus,
		Religion:          userProfileDTO.Religion,
		Height:            userProfileDTO.Height,
		Weight:            userProfileDTO.Weight,
		LookingFor:        userProfileDTO.LookingFor,
		Exercise:          userProfileDTO.Exercise,
		Drink:             userProfileDTO.Drink,
		Smoke:             userProfileDTO.Smoke,
		About:             userProfileDTO.About,
		Pronoun:           userProfileDTO.Pronoun,
		//Images:            userImagesES,
		//Nudges:            userNudgesES,
		UserSearchProfile: userProfile.UserSearchProfile,
	}
	return userProfileES, nil
}

func getLocation(loc []float64) *model.Location {

	return &model.Location{
		Point: wkb.Point{
			Point: geom.NewPoint(geom.XY).MustSetCoords(loc).SetSRID(4326),
		},
	}
}

func (u *UserProfileService) GetInterestsCategory() (model.InterestsListResponse, error) {
	var interests model.InterestsListResponse

	interestsDB, err := u.interestsDao.GetInterestsList()
	if err != nil {
		zapLogger.Logger.Error("error in getting interests details from DB")
		return interests, err
	}

	for _, val := range interestsDB {
		interestResp := model.InterestsResponse{
			InterestValue: val.InterestValues,
			Emoticon:      val.Emoticon,
		}
		switch {
		case val.InterestID == 1:
			interests.SelfCare = append(interests.SelfCare, interestResp)
		case val.InterestID == 2:
			interests.Sports = append(interests.Sports, interestResp)
		case val.InterestID == 3:
			interests.Creativity = append(interests.Creativity, interestResp)
		case val.InterestID == 4:
			interests.GoingOut = append(interests.GoingOut, interestResp)
		case val.InterestID == 5:
			interests.FilmAndTV = append(interests.FilmAndTV, interestResp)
		case val.InterestID == 6:
			interests.StayingIn = append(interests.StayingIn, interestResp)
		case val.InterestID == 7:
			interests.Reading = append(interests.Reading, interestResp)
		case val.InterestID == 8:
			interests.Music = append(interests.Music, interestResp)
		case val.InterestID == 9:
			interests.FoodAndDrink = append(interests.FoodAndDrink, interestResp)
		case val.InterestID == 10:
			interests.Travelling = append(interests.Travelling, interestResp)
		case val.InterestID == 11:
			interests.Pets = append(interests.Pets, interestResp)
		case val.InterestID == 12:
			interests.ValuesAndTraits = append(interests.ValuesAndTraits, interestResp)
		case val.InterestID == 13:
			interests.PlutoValuesAndAllyship = append(interests.PlutoValuesAndAllyship, interestResp)
		}
	}

	return interests, nil
}

func (u *UserProfileService) CreateUserInterest(interests model.InterestsCategory, userID int) ([]model.UserInterests, error) {
	userInterests := make([]model.UserInterests, 0)
	var tempUI model.UserInterests

	userInterests, err := u.interestsDao.GetUserInterests(userID)
	if err == nil {
		zapLogger.Logger.Error("user interests already exist in DB")
		return userInterests, errors.New("user interests already exist in DB")
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		zapLogger.Logger.Error("error in getting user interests from DB")
		return userInterests, err
	}

	if len(interests.SelfCare) != 0 {

		for _, interestVal := range interests.SelfCare {
			tempUI.UserID = userID
			tempUI.InterestID = 1
			tempUI.InterestValues = interestVal
			userInterests = append(userInterests, tempUI)
		}
	}
	if len(interests.Sports) != 0 {

		for _, interestVal := range interests.Sports {
			tempUI.UserID = userID
			tempUI.InterestID = 2
			tempUI.InterestValues = interestVal
			userInterests = append(userInterests, tempUI)
		}
	}
	if len(interests.Creativity) != 0 {

		for _, interestVal := range interests.Creativity {
			tempUI.UserID = userID
			tempUI.InterestID = 3
			tempUI.InterestValues = interestVal
			userInterests = append(userInterests, tempUI)
		}
	}
	if len(interests.GoingOut) != 0 {

		for _, interestVal := range interests.GoingOut {
			tempUI.UserID = userID
			tempUI.InterestID = 4
			tempUI.InterestValues = interestVal
			userInterests = append(userInterests, tempUI)
		}
	}
	if len(interests.FilmAndTV) != 0 {

		for _, interestVal := range interests.FilmAndTV {
			tempUI.UserID = userID
			tempUI.InterestID = 5
			tempUI.InterestValues = interestVal
			userInterests = append(userInterests, tempUI)
		}
	}
	if len(interests.StayingIn) != 0 {

		for _, interestVal := range interests.StayingIn {
			tempUI.UserID = userID
			tempUI.InterestID = 6
			tempUI.InterestValues = interestVal
			userInterests = append(userInterests, tempUI)
		}
	}
	if len(interests.Reading) != 0 {

		for _, interestVal := range interests.Reading {
			tempUI.UserID = userID
			tempUI.InterestID = 7
			tempUI.InterestValues = interestVal
			userInterests = append(userInterests, tempUI)
		}
	}
	if len(interests.Music) != 0 {

		for _, interestVal := range interests.Music {
			tempUI.UserID = userID
			tempUI.InterestID = 8
			tempUI.InterestValues = interestVal
			userInterests = append(userInterests, tempUI)
		}
	}
	if len(interests.FoodAndDrink) != 0 {

		for _, interestVal := range interests.FoodAndDrink {
			tempUI.UserID = userID
			tempUI.InterestID = 9
			tempUI.InterestValues = interestVal
			userInterests = append(userInterests, tempUI)
		}
	}
	if len(interests.Travelling) != 0 {

		for _, interestVal := range interests.Travelling {
			tempUI.UserID = userID
			tempUI.InterestID = 10
			tempUI.InterestValues = interestVal
			userInterests = append(userInterests, tempUI)
		}
	}
	if len(interests.Pets) != 0 {

		for _, interestVal := range interests.Pets {
			tempUI.UserID = userID
			tempUI.InterestID = 11
			tempUI.InterestValues = interestVal
			userInterests = append(userInterests, tempUI)
		}
	}
	if len(interests.ValuesAndTraits) != 0 {

		for _, interestVal := range interests.ValuesAndTraits {
			tempUI.UserID = userID
			tempUI.InterestID = 12
			tempUI.InterestValues = interestVal
			userInterests = append(userInterests, tempUI)
		}
	}
	if len(interests.PlutoValuesAndAllyship) != 0 {

		for _, interestVal := range interests.PlutoValuesAndAllyship {
			tempUI.UserID = userID
			tempUI.InterestID = 13
			tempUI.InterestValues = interestVal
			userInterests = append(userInterests, tempUI)
		}
	}

	userInterests, err = u.interestsDao.CreateUserInterests(userInterests, userID)
	if err != nil {
		zapLogger.Logger.Error("[CreateUserInterest] error in creating user interests")
		return userInterests, err
	}

	return userInterests, nil
}

func (u *UserProfileService) GetUserInterests(userID int) (model.InterestsCategory, error) {
	var interests model.InterestsCategory

	userInterests, err := u.interestsDao.GetUserInterests(userID)
	if err != nil {
		zapLogger.Logger.Error("error in getting user interests from DB")
		return interests, err
	}

	for _, val := range userInterests {
		switch {
		case val.InterestID == 1:
			interests.SelfCare = append(interests.SelfCare, val.InterestValues)
		case val.InterestID == 2:
			interests.Sports = append(interests.Sports, val.InterestValues)
		case val.InterestID == 3:
			interests.Creativity = append(interests.Creativity, val.InterestValues)
		case val.InterestID == 4:
			interests.GoingOut = append(interests.GoingOut, val.InterestValues)
		case val.InterestID == 5:
			interests.FilmAndTV = append(interests.FilmAndTV, val.InterestValues)
		case val.InterestID == 6:
			interests.StayingIn = append(interests.StayingIn, val.InterestValues)
		case val.InterestID == 7:
			interests.Reading = append(interests.Reading, val.InterestValues)
		case val.InterestID == 8:
			interests.Music = append(interests.Music, val.InterestValues)
		case val.InterestID == 9:
			interests.FoodAndDrink = append(interests.FoodAndDrink, val.InterestValues)
		case val.InterestID == 10:
			interests.Travelling = append(interests.Travelling, val.InterestValues)
		case val.InterestID == 11:
			interests.Pets = append(interests.Pets, val.InterestValues)
		case val.InterestID == 12:
			interests.ValuesAndTraits = append(interests.ValuesAndTraits, val.InterestValues)
		case val.InterestID == 13:
			interests.PlutoValuesAndAllyship = append(interests.PlutoValuesAndAllyship, val.InterestValues)
		}
	}

	return interests, nil
}

func (u *UserProfileService) CheckAllowedFileType(extension string) bool {
	extensions := [5]string{".jpeg", ".png", ".gif", ".webp", ".jpg"}
	for _, ext := range extensions {
		if ext == extension {
			return true
		}
	}
	return false
}

func (u *UserProfileService) GetNudgesService() ([]model.Nudge, error) {
	nudges, err := u.userNudgesDao.GetNudgesDB()
	if err != nil {
		zapLogger.Logger.Error("error in getting nudges list from DB")
		return nudges, err
	}

	return nudges, nil
}

func (u *UserProfileService) CreateUserNudgeService(nudges model.NudgeDetail, userID int) (model.UserNudge, error) {
	userNudge := model.UserNudge{
		UserID:   userID,
		Question: nudges.Question,
		Answer:   nudges.Answer,
		Order:    nudges.Order,
		MediaURL: nudges.MediaURL,
		Type:     nudges.Type,
	}

	userNudge, err := u.userNudgesDao.CreateUserNudgesDB(userNudge)
	if err != nil {
		zapLogger.Logger.Error("error in creating user nudge in DB")
		return userNudge, err
	}

	return userNudge, nil
}

func (u *UserProfileService) GetUserNudgesService(userID int) ([]model.UserNudge, error) {
	userNudge, err := u.userNudgesDao.GetUserNudgesDB(userID)
	if err != nil {
		zapLogger.Logger.Error("error in getting user nudge from DB")
		return userNudge, err
	}

	for idx := range userNudge {
		if userNudge[idx].MediaURL != "" {
			key := strings.ReplaceAll(strings.TrimPrefix(userNudge[idx].MediaURL, S3_BUCKET_PATH), "%3A", ":")
			signedURL, err := u.s3Service.SignS3FilesUrl(user_profile_S3_bucket, key)
			if err != nil {
				zapLogger.Logger.Error("error in getting signed url ", zap.Error(err))
				return nil, err
			}
			userNudge[idx].MediaURL = signedURL
		}
	}

	return userNudge, nil
}

func (u *UserProfileService) UpdateMediaProfile(user model.UserProfile, mediaDetails model.MediaOrderId) (model.UserMedia, error) {
	userMedia, err := u.userMedia.FindById(context.Background(), mediaDetails.MediaID)
	userMedia.OrderId = mediaDetails.OrderID
	userMedia, err = u.userMedia.UpdateProfileMedia(context.Background(), userMedia, mediaDetails.MediaID)
	if err != nil {
		zapLogger.Logger.Error("error in updating user media profile")
		return userMedia, err
	}

	return userMedia, nil
}

func (u *UserProfileService) CheckAllowedAudioFileType(extension string) bool {
	extensions := [4]string{".mp3", ".wav", ".aac", ".m4a"}
	for _, ext := range extensions {
		if ext == extension {
			return true
		}
	}
	return false
}

func (u *UserProfileService) UploadNudgeMediaToS3(userId string, file *multipart.FileHeader, nudgeDetails model.NudgeDetail) (model.NudgeDetail, error) {
	fileExt := filepath.Ext(file.Filename)

	if u.CheckAllowedAudioFileType(fileExt) {
		nudgeDetails.Type = "audio"
	} else if u.CheckAllowedVideoFileType(fileExt) {
		nudgeDetails.Type = "video"
	} else {
		zapLogger.Logger.Error("file extension is not supported", zap.Any("extension_type:", fileExt))
		return nudgeDetails, errors.New("file extension is not supported")
	}

	filename := fmt.Sprintf("%v", time.Now().UTC().Format("2006-01-02T15:04:05.00000")) + fileExt
	tempFile, _ := file.Open()
	S3filepath := userId + "/nudge_media/" + filename

	result, err := u.s3Service.UploadFileToS3(user_profile_S3_bucket, S3filepath, tempFile, filename)
	if err != nil {
		zapLogger.Logger.Error("error in uploading file to S3:", zap.Error(err))
		return nudgeDetails, err
	}
	nudgeDetails.MediaURL = result
	return nudgeDetails, nil
}

func (u *UserProfileService) CreateProfileIndex() error {
	err := u.esIndex.CreateIndex()
	if err != nil {
		log.Println("error in creating index in elasticSearch")
		return err
	}

	return nil
}

func (u *UserProfileService) UpdateUserNudge(nudgeDetails model.NudgeDetail, nudgeID int) (model.UserNudge, error) {
	userNudge, err := u.userNudgesDao.GetUserNudgeById(nudgeID)
	if err != nil {
		zapLogger.Logger.Error("error in getting user nudge from DB")
		return userNudge, err
	}

	if (nudgeDetails.MediaURL != "" && userNudge.MediaURL != "") || (nudgeDetails.Type == "text" && userNudge.Type != "text") {
		key := strings.ReplaceAll(strings.TrimPrefix(userNudge.MediaURL, S3_BUCKET_PATH), "%3A", ":")
		err := u.s3Service.DeleteFile(user_profile_S3_bucket, key)
		if err != nil {
			zapLogger.Logger.Error("error in deleting file from S3:", zap.Error(err))
			return userNudge, err
		}
	}
	userNudge.Question = nudgeDetails.Question
	userNudge.Answer = nudgeDetails.Answer
	userNudge.Order = nudgeDetails.Order
	userNudge.MediaURL = nudgeDetails.MediaURL
	userNudge.Type = nudgeDetails.Type
	userNudge, err = u.userNudgesDao.UpdateUserNudge(userNudge, nudgeID)
	if err != nil {
		zapLogger.Logger.Error("error in updating user nudge in DB")
		return userNudge, err
	}

	return userNudge, nil
}

func (u *UserProfileService) DeleteUserNudge(nudgeID int) error {
	userNudge, err := u.userNudgesDao.GetUserNudgeById(nudgeID)
	if err != nil {
		zapLogger.Logger.Error("error in getting user nudge from DB")
		return err
	}

	if userNudge.MediaURL != "" {
		key := strings.ReplaceAll(strings.TrimPrefix(userNudge.MediaURL, S3_BUCKET_PATH), "%3A", ":")
		err := u.s3Service.DeleteFile(user_profile_S3_bucket, key)
		if err != nil {
			zapLogger.Logger.Error("error in deleting file from S3:", zap.Error(err))
			return err
		}
	}

	err = u.userNudgesDao.DeleteUserNudge(nudgeID)
	if err != nil {
		zapLogger.Logger.Error("error in deleting user nudge from DB")
		return err
	}

	return nil
}

func (u *UserProfileService) UpdateUserLocation(userProfile model.UserProfile, location model.UserLocation) error {
	userProfile.Latitude = &location.Latitude
	userProfile.Longitude = &location.Longitude

	_, err := u.userProfileDao.UpdateUserProfile(context.Background(), userProfile)
	if err != nil {
		zapLogger.Logger.Error("error in updating user location in DB")
		return err
	}

	//update location in ElasticSearch also by creating elastic object=
	//and calling update function

	userProfileES, err := u.esIndex.GetUserProfile(int(userProfile.ID))

	if err != nil {
		zapLogger.Logger.Error("error in getting user profile from DB")
		return err
	}

	userProfileES.Location = []float64{location.Latitude, location.Longitude}
	jsonString, _ := json.Marshal(userProfileES)
	err = u.esIndex.IndexUserProfile(userProfileES, jsonString)

	if err != nil {
		zapLogger.Logger.Error("error in updating user location in ElasticSearch")
		return err
	}
	return nil
}

func (u *UserProfileService) GetFiltersList() (model.Filters, error) {
	var filterLists model.Filters
	exercise, err := u.filtersDao.GetExerciseFilterList()
	if err != nil {
		zapLogger.Logger.Error("error in getting exercise filter list from DB")
		return model.Filters{}, err
	}
	for _, val := range exercise {
		filterLists.Exercise = append(filterLists.Exercise, val.DoTheyExercise)
	}

	starSign, err := u.filtersDao.GetStarSignFilterList()
	if err != nil {
		zapLogger.Logger.Error("error in getting star sign filter list from DB")
		return model.Filters{}, err
	}
	for _, val := range starSign {
		filterLists.StarSign = append(filterLists.StarSign, val.StarSign)
	}

	education, err := u.filtersDao.GetEducationFilterList()
	if err != nil {
		zapLogger.Logger.Error("error in getting education filter list from DB")
		return model.Filters{}, err
	}
	for _, val := range education {
		filterLists.Education = append(filterLists.Education, val.Education)
	}

	drink, err := u.filtersDao.GetDrinkFilterList()
	if err != nil {
		zapLogger.Logger.Error("error in getting drink filter list from DB")
		return model.Filters{}, err
	}
	for _, val := range drink {
		filterLists.Drink = append(filterLists.Drink, val.DoTheyDrink)
	}

	smoke, err := u.filtersDao.GetSmokeFilterList()
	if err != nil {
		zapLogger.Logger.Error("error in getting smoke filter list from DB")
		return model.Filters{}, err
	}
	for _, val := range smoke {
		filterLists.Smoke = append(filterLists.Smoke, val.DoTheySmoke)
	}

	lookingFor, err := u.filtersDao.GetLookingForFilterList()
	if err != nil {
		zapLogger.Logger.Error("error in getting looking for filter list from DB")
		return model.Filters{}, err
	}
	for _, val := range lookingFor {
		filterLists.LookingFor = append(filterLists.LookingFor, val.LookingFor)
	}

	religion, err := u.filtersDao.GetReligionFilterList()
	if err != nil {
		zapLogger.Logger.Error("error in getting religion filter list from DB")
		return model.Filters{}, err
	}
	for _, val := range religion {
		filterLists.Religion = append(filterLists.Religion, val.Religion)
	}

	politicsLikes, err := u.filtersDao.GetPoliticsLikesFilterList()
	if err != nil {
		zapLogger.Logger.Error("error in getting politics likes filter list from DB")
		return model.Filters{}, err
	}
	for _, val := range politicsLikes {
		filterLists.PoliticsLikes = append(filterLists.PoliticsLikes, val.PoliticsLikes)
	}

	childrenFilter, err := u.filtersDao.GetChildrenFilterList()
	if err != nil {
		zapLogger.Logger.Error("error in getting children filter list from DB")
		return model.Filters{}, err
	}
	for _, val := range childrenFilter {
		filterLists.Children = append(filterLists.Children, val.HaveOrWantChildren)
	}

	return filterLists, nil

}

func (u *UserProfileService) UpdateUserInterest(interests []model.InterestDetails, userID int) ([]model.UserInterests, error) {
	userInterests, err := u.interestsDao.GetUserInterests(userID)
	if err != nil {
		zapLogger.Logger.Error("error in getting user interests from DB")
		return nil, err
	}

	if len(userInterests) != len(interests) {
		return nil, errors.New("interests length mismatch")
	}

	for i, interest := range interests {
		userInterests[i].InterestID = interest.InterestID
		userInterests[i].InterestValues = interest.InterestValues
	}

	updatedInterest, err := u.interestsDao.UpdateUserInterests(userInterests, userID)
	if err != nil {
		zapLogger.Logger.Error("error in updating user interests in DB")
		return nil, err
	}

	return updatedInterest, nil

}

func (u *UserProfileService) CheckAllowedVideoFileType(extension string) bool {
	extensions := [5]string{".mp4", ".mov", ".flv", ".avi", ".wmv"}
	for _, ext := range extensions {
		if ext == extension {
			return true
		}
	}
	return false
}

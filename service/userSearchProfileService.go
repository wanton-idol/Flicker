package service

import (
	"context"
	"encoding/json"
	"github.com/SuperMatch/model"
	"github.com/SuperMatch/model/dto"
	elasticsearchPkg "github.com/SuperMatch/model/elasticSearch"
	"github.com/SuperMatch/pkg/db/dao"
	pkg "github.com/SuperMatch/pkg/elasticSeach"
	utils "github.com/SuperMatch/utilities"
	"github.com/SuperMatch/zapLogger"
	"go.uber.org/zap"
)

type UserSearchProfileServiceInterface interface {
	CreateUserSearchProfile(ctx context.Context, profile model.UserProfile, userProfile dto.UserSearchProfile) error
	UpdateUserSearchProfile(ctx context.Context, profile model.UserProfile, userSearchProfile dto.UserSearchProfile) (dto.UserSearchProfile, error)
	FindByProfileId(ctx context.Context, profileId int) (model.UserSearchProfile, error)
	FindByUserId(ctx context.Context, userId int) (model.UserSearchProfile, error)
	UpdateAdvancedFilters(advFil dto.AdvancedFilter, profile model.UserProfile) (dto.AdvancedFilter, error)
}

type UserSearchProfileService struct {
	esIndex              pkg.ElasticSearchIndexer
	userSearchProfileDao dao.UserSearchProfileRepository
	advancedFilterDao    dao.AdvancedFilterRepository
}

func NewUserSearchProfileService() *UserSearchProfileService {
	return &UserSearchProfileService{
		esIndex:              pkg.NewElasticSearchIndexerImpl(),
		userSearchProfileDao: dao.NewUserSearchProfile(),
		advancedFilterDao:    dao.NewAdvancedFilterRepository(),
	}
}

func (u *UserSearchProfileService) CreateUserSearchProfile(ctx context.Context, profile model.UserProfile, userSearchProfile dto.UserSearchProfile) (dto.UserSearchProfile, error) {

	// create map of interfaces for user search profile
	var userSearchProfileMap = make(map[string]interface{})
	userSearchProfileMap["user_profile_id"] = profile.ID
	userSearchProfileMap["user_id"] = profile.UserId
	userSearchProfileMap["min_age"] = userSearchProfile.MinAge
	userSearchProfileMap["max_age"] = userSearchProfile.MaxAge
	userSearchProfileMap["distance"] = userSearchProfile.Distance
	genderString, err := utils.ConvertCustomTypeGenderToString(userSearchProfile.Gender)

	if err != nil {
		return dto.UserSearchProfile{}, err
	}
	userSearchProfileMap["gender"] = genderString
	userSearchProfileMap["snooze"] = userSearchProfile.Snooze
	userSearchProfileMap["hide_my_name"] = userSearchProfile.HideMyName

	_, err = u.userSearchProfileDao.CreateUserSearchProfile(ctx, userSearchProfileMap)

	if err != nil {
		return userSearchProfile, err
	}
	searchProfileES := elasticsearchPkg.UserSearchProfile{
		MinAge:     userSearchProfile.MinAge,
		MaxAge:     userSearchProfile.MaxAge,
		Distance:   userSearchProfile.Distance,
		Gender:     userSearchProfile.Gender,
		Snooze:     userSearchProfile.Snooze,
		HideMyName: userSearchProfile.HideMyName,
	}

	userProfileES := elasticsearchPkg.UserProfile{
		UserSearchProfile: searchProfileES,
	}

	jsonString, err := json.Marshal(userProfileES)
	userProfileES.Id = int(profile.ID)
	zapLogger.Logger.Debug(string(jsonString))
	if err != nil {
		return userSearchProfile, err
	}

	//update search Profile in elasticsearch
	err = u.esIndex.UpdateSearchProfile(userProfileES, jsonString)

	if err != nil {
		return userSearchProfile, err
	}

	return userSearchProfile, nil
}

func (u *UserSearchProfileService) UpdateUserSearchProfile(ctx context.Context, profile model.UserProfile, userSearchProfile dto.UserSearchProfile) (dto.UserSearchProfile, error) {
	var userSearchProfileMap = make(map[string]interface{})
	userSearchProfileMap["user_profile_id"] = profile.ID
	userSearchProfileMap["user_id"] = profile.UserId
	userSearchProfileMap["min_age"] = userSearchProfile.MinAge
	userSearchProfileMap["max_age"] = userSearchProfile.MaxAge
	userSearchProfileMap["distance"] = userSearchProfile.Distance
	genderString, err := utils.ConvertCustomTypeGenderToString(userSearchProfile.Gender)

	if err != nil {
		return dto.UserSearchProfile{}, err
	}
	userSearchProfileMap["gender"] = genderString
	userSearchProfileMap["snooze"] = userSearchProfile.Snooze
	userSearchProfileMap["hide_my_name"] = userSearchProfile.HideMyName

	//update user search profile in database
	_, err = u.userSearchProfileDao.UpdateUserSearchProfile(ctx, userSearchProfileMap)

	if err != nil {
		return userSearchProfile, err
	}
	searchProfileES := elasticsearchPkg.UserSearchProfile{
		MinAge:     userSearchProfile.MinAge,
		MaxAge:     userSearchProfile.MaxAge,
		Distance:   userSearchProfile.Distance,
		Gender:     userSearchProfile.Gender,
		Snooze:     userSearchProfile.Snooze,
		HideMyName: userSearchProfile.HideMyName,
	}

	userProfileES, err := u.esIndex.GetUserProfile(int(profile.ID))

	//setting advanced filter in search profile
	searchProfileES.AdvancedFilter = userProfileES.UserSearchProfile.AdvancedFilter
	userProfileES.UserSearchProfile = searchProfileES

	jsonString, err := json.Marshal(userProfileES)
	userProfileES.Id = int(profile.ID)
	zapLogger.Logger.Debug(string(jsonString))
	if err != nil {
		return userSearchProfile, err
	}

	//update search Profile in elasticsearch
	err = u.esIndex.UpdateSearchProfile(userProfileES, jsonString)

	if err != nil {
		return userSearchProfile, err
	}

	return userSearchProfile, nil
}

func (u *UserSearchProfileService) FindByProfileId(ctx context.Context, profileId int) (dto.UserSearchProfile, error) {

	profile, err := u.userSearchProfileDao.FindByProfileId(ctx, profileId)
	if err != nil {
		return dto.UserSearchProfile{}, err
	}
	userSearchProfile := dto.UserSearchProfile{
		MinAge:     profile.MinAge,
		MaxAge:     profile.MaxAge,
		Distance:   profile.Distance,
		Gender:     utils.ConvertStringToCustomGender(profile.Gender),
		Snooze:     profile.Snooze,
		HideMyName: profile.HideMyName,
	}
	return userSearchProfile, nil
}

func (u *UserSearchProfileService) FindByUserId(ctx context.Context, userId int) (dto.UserSearchProfile, error) {

	profile, err := u.userSearchProfileDao.FindByUserId(ctx, userId)
	if err != nil {
		return dto.UserSearchProfile{}, err
	}
	userSearchProfile := dto.UserSearchProfile{
		MinAge:     profile.MinAge,
		MaxAge:     profile.MaxAge,
		Distance:   profile.Distance,
		Gender:     utils.ConvertStringToCustomGender(profile.Gender),
		Snooze:     profile.Snooze,
		HideMyName: profile.HideMyName,
	}
	return userSearchProfile, nil
}

func (u *UserSearchProfileService) UpdateAdvancedFilters(advFil dto.AdvancedFilter, profile model.UserProfile) (dto.AdvancedFilter, error) {

	advancedFilterDB := model.AdvancedFilter{
		UserId:             profile.UserId,
		UserProfileId:      int(profile.ID),
		IsProfileVerified:  advFil.IsProfileVerified,
		IsOnline:           advFil.IsOnline,
		Height:             utils.ConvertSliceToCommaSeparatedString(advFil.Height),
		Exercise:           utils.ConvertCustomTypeExerciseToString(advFil.Exercise),
		Religion:           utils.ConvertCustomTypeReligionToString(advFil.Religion),
		Education:          utils.ConvertCustomTypeEducationLevelToString(advFil.EducationLevel),
		Occupation:         utils.ConvertCustomTypeOccupationToString(advFil.Occupation),
		MaritalStatus:      utils.ConvertCustomTypeMaritalStatusToString(advFil.MaritalStatus),
		Drink:              utils.ConvertCustomTypeDrinkToString(advFil.Drink),
		Smoke:              utils.ConvertCustomTypeSmokeToString(advFil.Smoke),
		IncognitoMode:      advFil.IncognitoMode,
		StarSign:           utils.ConvertSliceOfStringToCommaSeparatedString(advFil.StarSign),
		PoliticsLikes:      utils.ConvertSliceOfStringToCommaSeparatedString(advFil.PoliticsLikes),
		HaveOrWantChildren: utils.ConvertSliceOfStringToCommaSeparatedString(advFil.HaveOrWantChildren),
		LookingFor:         utils.ConvertSliceOfStringToCommaSeparatedString(advFil.LookingFor),
	}

	//Update Advanced Filter in database
	_, err := u.advancedFilterDao.UpdateAdvancedFilter(advancedFilterDB)

	if err != nil {
		zapLogger.Logger.Error("error in updating advanced filters:", zap.Error(err))
		return advFil, err
	}

	userProfile, err := u.esIndex.GetUserProfile(int(profile.ID))

	if err != nil {
		zapLogger.Logger.Error("error in getting user profile from elasticsearch:", zap.Error(err))
		return advFil, err
	}

	advancedFilterES := elasticsearchPkg.AdvancedFilter{
		IsProfileVerified:  advFil.IsProfileVerified,
		IsOnline:           advFil.IsOnline,
		Height:             advFil.Height,
		Exercise:           advFil.Exercise,
		Religion:           advFil.Religion,
		EducationLevel:     advFil.EducationLevel,
		Occupation:         advFil.Occupation,
		MaritalStatus:      advFil.MaritalStatus,
		Drink:              advFil.Drink,
		Smoke:              advFil.Smoke,
		IncognitoMode:      advFil.IncognitoMode,
		StarSign:           advFil.StarSign,
		PoliticsLikes:      advFil.PoliticsLikes,
		HaveOrWantChildren: advFil.HaveOrWantChildren,
		LookingFor:         advFil.LookingFor,
	}

	userProfile.UserSearchProfile.AdvancedFilter = advancedFilterES
	jsonString, err := json.Marshal(userProfile)
	zapLogger.Logger.Debug(string(jsonString))
	if err != nil {
		return advFil, err
	}

	//update advanced filter in elasticsearch
	err = u.esIndex.UpdateSearchProfile(userProfile, jsonString)

	if err != nil {
		return advFil, err
	}
	return advFil, nil
}

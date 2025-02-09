package tests

import (
	"encoding/json"
	"github.com/SuperMatch/model"
	elasticsearchPkg "github.com/SuperMatch/model/elasticSearch"
	"github.com/SuperMatch/pkg/db/dao/mocks"
	esMocks "github.com/SuperMatch/pkg/elasticseach/mocks"
	mockService "github.com/SuperMatch/service/mocks"
	utils "github.com/SuperMatch/utilities"
	"github.com/golang/mock/gomock"
	"gorm.io/gorm"
	"strings"
	"testing"
	"time"
)

func userProfileESFunc() elasticsearchPkg.UserProfile {
	return elasticsearchPkg.UserProfile{
		Id:                1,
		UserId:            1,
		FirstName:         utils.ConvertStringToStringPointer("John"),
		LastName:          utils.ConvertStringToStringPointer("Doe"),
		IsPremium:         utils.ConvertBoolToBoolPointer(true),
		IsVerified:        utils.ConvertBoolToBoolPointer(true),
		DateOfBirth:       utils.ConvertStringToCustomDOBPointer(time.Now()),
		Gender:            utils.ConvertStringToCustomGenderPointer("male"),
		SexualOrientation: utils.ConvertStringToCustomSexualOrientationPointer("heterosexual"),
		Location:          []float64{11.55, 65.24},
		Education:         nil,
		Occupation:        utils.ConvertStringToCustomOccupationPointer("student"),
		MaritalStatus:     utils.ConvertStringToCustomMaritalStatusPointer("single"),
		Religion:          utils.ConvertStringToCustomReligionPointer("hindu"),
		Height:            utils.ConvertIntToIntPointer(165),
		Weight:            utils.ConvertIntToIntPointer(65),
		LookingFor:        utils.ConvertStringToCustomLookingForPointer("relationship"),
		Exercise:          utils.ConvertStringToCustomExercisePointer("everyday"),
		Drink:             utils.ConvertStringToCustomDrinkPointer("never"),
		Smoke:             utils.ConvertStringToCustomSmokePointer("never"),
		About:             utils.ConvertStringToStringPointer("Developer"),
	}
}

func userProfileDBFunc() model.UserProfile {
	return model.UserProfile{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		UserId:            1,
		FirstName:         utils.ConvertStringToStringPointer("John"),
		LastName:          utils.ConvertStringToStringPointer("Doe"),
		IsPremium:         true,
		IsVerified:        true,
		DateOfBirth:       utils.ConvertStringToCustomDOBPointer(time.Now()),
		Gender:            utils.ConvertStringToCustomGenderPointer("male"),
		SexualOrientation: utils.ConvertStringToCustomSexualOrientationPointer("heterosexual"),
		Latitude:          utils.ConvertFloatToFloatPointer(11.55),
		Longitude:         utils.ConvertFloatToFloatPointer(65.24),
		EducationLevel:    utils.ConvertStringToCustomEduLevelPointer("college"),
		College:           utils.ConvertStringToStringPointer("IIT"),
		Occupation:        utils.ConvertStringToCustomOccupationPointer("student"),
		MaritalStatus:     utils.ConvertStringToCustomMaritalStatusPointer("single"),
		Religion:          utils.ConvertStringToCustomReligionPointer("hindu"),
		Height:            utils.ConvertIntToIntPointer(165),
		Weight:            utils.ConvertIntToIntPointer(65),
		LookingFor:        utils.ConvertStringToCustomLookingForPointer("relationship"),
		Exercise:          utils.ConvertStringToCustomExercisePointer("everyday"),
		Drink:             utils.ConvertStringToCustomDrinkPointer("never"),
		Smoke:             utils.ConvertStringToCustomSmokePointer("never"),
		About:             utils.ConvertStringToStringPointer("Developer"),
	}
}

func userProfileMedia() model.UserMedia {
	return model.UserMedia{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		UserProfileId: 1,
		UserId:        1,
		URL:           "https://user-profile-supermatch.s3.ap-south-1.amazonaws.com/1/profile/2023-05-15T06%3A05%3A35.53984.png",
		OrderId:       1,
		ImageText:     "Random image text",
		Latitude:      11.55,
		Longitude:     65.24,
	}
}

func TestGetUserProfile(t *testing.T) {
	userProfile := userProfileESFunc()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockES := esMocks.NewMockElasticSearchIndexer(ctrl)
	mockES.EXPECT().GetUserProfile(gomock.Eq(1)).
		Return(userProfile, nil).
		AnyTimes()
}

func TestGetUserProfileFromDB(t *testing.T) {
	userProfile := userProfileDBFunc()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserProfile := mocks.NewMockUserProfileRepository(ctrl)
	mockUserProfile.EXPECT().FindByUserId(gomock.Any(), gomock.Eq(1)).
		Return(userProfile, nil).
		AnyTimes()
}

func TestUpdatePremiumProfile(t *testing.T) {
	userProfile := userProfileDBFunc()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserProfile := mocks.NewMockUserProfileRepository(ctrl)
	mockUserProfile.EXPECT().UpdateUserProfile(gomock.Any(), gomock.Eq(userProfile)).
		Return(userProfile, nil).
		AnyTimes()
}

func TestUpdateUserProfileFirst(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userProfileDB := userProfileDBFunc()
	mockUserProfile := mocks.NewMockUserProfileRepository(ctrl)
	mockUserProfile.EXPECT().UpdateUserProfile(gomock.Any(), gomock.Eq(userProfileDB)).
		Return(userProfileDB, nil).
		AnyTimes()

	userProfileES := userProfileESFunc()
	jsonString, err := json.Marshal(userProfileES)
	if err != nil {
		t.Error("error marshalling user profile for elasticsearch")
	}
	mockES := esMocks.NewMockElasticSearchIndexer(ctrl)
	mockES.EXPECT().IndexUserProfile(gomock.Eq(userProfileES), gomock.Eq(jsonString)).
		Return(nil).
		AnyTimes()
}

func TestCreateUserProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userProfileDB := userProfileDBFunc()
	mockUserProfile := mocks.NewMockUserProfileRepository(ctrl)
	mockUserProfile.EXPECT().CreateUserProfile(gomock.Any(), gomock.Eq(userProfileDB)).
		Return(userProfileDB, nil).
		AnyTimes()

	mp := make(map[string]interface{})
	mp["user_profile_id"] = 1
	mp["user_id"] = 1

	mockUserSearchProfile := mocks.NewMockUserSearchProfileRepository(ctrl)
	mockUserSearchProfile.EXPECT().CreateUserSearchProfile(gomock.Any(), gomock.Eq(mp)).
		Return(model.UserSearchProfile{
			UserId:        1,
			UserProfileId: 1,
		}, nil).AnyTimes()

	advancedFilter := model.AdvancedFilter{
		UserId:        1,
		UserProfileId: 1,
	}
	mockAdvFil := mocks.NewMockAdvancedFilterRepository(ctrl)
	mockAdvFil.EXPECT().CreateAdvancedFilter(gomock.Any(), gomock.Eq(advancedFilter)).
		Return(nil).
		AnyTimes()

	userProfileES := elasticsearchPkg.UserProfile{
		Id:     1,
		UserId: 1,
		UserSearchProfile: elasticsearchPkg.UserSearchProfile{
			Gender:         []elasticsearchPkg.Gender{elasticsearchPkg.MALE, elasticsearchPkg.FEMALE, elasticsearchPkg.BINARY},
			AdvancedFilter: elasticsearchPkg.AdvancedFilter{},
		},
	}
	jsonString, err := json.Marshal(userProfileES)
	if err != nil {
		t.Error("error marshalling user profile for elasticsearch")
	}

	mockES := esMocks.NewMockElasticSearchIndexer(ctrl)
	mockES.EXPECT().IndexUserProfile(userProfileES, gomock.Eq(jsonString)).
		Return(nil).
		AnyTimes()

}

func TestUpdateUserProfile(t *testing.T) {
	profileMap := make(map[string]interface{})
	profileMap["ID"] = 1
	profileMap["user_id"] = 1
	profileMap["first_name"] = utils.ConvertStringToStringPointer("John")
	profileMap["last_name"] = utils.ConvertStringToStringPointer("Doe")
	profileMap["date_of_birth"] = utils.ConvertStringToCustomDOBPointer(time.Now())
	profileMap["gender"] = utils.ConvertStringToCustomGenderPointer("male")
	profileMap["sexual_orientation"] = utils.ConvertStringToCustomSexualOrientationPointer("heterosexual")
	profileMap["latitude"] = utils.ConvertFloatToFloatPointer(11.55)
	profileMap["longitude"] = utils.ConvertFloatToFloatPointer(65.24)
	profileMap["education_level"] = utils.ConvertStringToCustomEduLevelPointer("college")
	profileMap["college"] = utils.ConvertStringToStringPointer("IIT")
	profileMap["occupation"] = utils.ConvertStringToCustomOccupationPointer("student")
	profileMap["marital_status"] = utils.ConvertStringToCustomMaritalStatusPointer("single")
	profileMap["religion"] = utils.ConvertStringToCustomReligionPointer("hindu")
	profileMap["height"] = utils.ConvertIntToIntPointer(165)
	profileMap["weight"] = utils.ConvertIntToIntPointer(65)
	profileMap["looking_for"] = utils.ConvertStringToCustomLookingForPointer("relationship")
	profileMap["exercise"] = utils.ConvertStringToCustomExercisePointer("everyday")
	profileMap["drink"] = utils.ConvertStringToCustomDrinkPointer("never")
	profileMap["smoke"] = utils.ConvertStringToCustomSmokePointer("never")
	profileMap["about"] = utils.ConvertStringToStringPointer("Developer")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserProfile := mocks.NewMockUserProfileRepository(ctrl)
	mockUserProfile.EXPECT().UpdateProfileByMap(gomock.Any(), gomock.Eq(profileMap)).
		Return(profileMap, nil).
		AnyTimes()

	userProfileES := userProfileESFunc()
	jsonString, err := json.Marshal(userProfileES)
	if err != nil {
		t.Error("error marshalling user profile for elasticsearch")
	}

	mockES := esMocks.NewMockElasticSearchIndexer(ctrl)
	mockES.EXPECT().UpdateSearchProfile(userProfileES, jsonString).
		Return(nil).
		AnyTimes()
}

func TestSearchProfile(t *testing.T) {
	userProfileES := userProfileESFunc()
	userProfiles := make([]elasticsearchPkg.UserProfile, 0)
	userProfiles = append(userProfiles, userProfileES)
	ctrl := gomock.NewController(t)
	ctrl.Finish()
	mockES := esMocks.NewMockElasticSearchIndexer(ctrl)
	mockES.EXPECT().SearchProfile(gomock.Any()).
		Return(userProfiles, nil).
		AnyTimes()
}

func TestSaveProfileMedia(t *testing.T) {
	userMedia := userProfileMedia()
	ctrl := gomock.NewController(t)
	ctrl.Finish()

	mockUserMedia := mocks.NewMockUserMediaRepository(ctrl)
	mockUserMedia.EXPECT().Insert(gomock.Any(), gomock.Eq(userMedia)).
		Return(nil).
		AnyTimes()
}

func TestGetProfileMediaByUserId(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctrl.Finish()

	userMedia := userProfileMedia()
	userMedias := make([]model.UserMedia, 0)
	userMedias = append(userMedias, userMedia)
	mockUserMedia := mocks.NewMockUserMediaRepository(ctrl)
	mockUserMedia.EXPECT().FindByUserId(gomock.Any(), gomock.Eq(1)).
		Return(userMedias, nil).
		AnyTimes()

	for idx := range userMedias {
		signedURL := ""
		key := strings.ReplaceAll(strings.TrimPrefix(userMedias[idx].URL, S3BucketPath), "%3A", ":")
		mockS3Service := mockService.NewMockS3ServiceInterface(ctrl)
		mockS3Service.EXPECT().SignS3FilesUrl(gomock.Eq(userProfileS3Bucket), gomock.Eq(key)).
			Return(signedURL, nil).
			AnyTimes()
		userMedias[idx].URL = signedURL
	}
}

func TestRemoveProfileMedia(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctrl.Finish()

	userMedia := userProfileMedia()
	mockUserMedia := mocks.NewMockUserMediaRepository(ctrl)
	mockUserMedia.EXPECT().FindByIdAndUserId(gomock.Any(), gomock.Eq(1), gomock.Eq(1)).
		Return(&userMedia, nil).
		AnyTimes()

	key := strings.ReplaceAll(strings.TrimPrefix(userMedia.URL, S3BucketPath), "%3A", ":")
	mockS3Service := mockService.NewMockS3ServiceInterface(ctrl)
	mockS3Service.EXPECT().DeleteFile(gomock.Eq(userProfileS3Bucket), gomock.Eq(key)).
		Return(nil).
		AnyTimes()

	mockUserMedia.EXPECT().DeleteById(gomock.Any(), gomock.Eq(1), gomock.Eq(1)).
		Return(nil).
		AnyTimes()
}

func TestGetInterestsService(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctrl.Finish()
	interests := make([]model.Interests, 0)
	mockInterests := mocks.NewMockInterestsDao(ctrl)
	mockInterests.EXPECT().GetInterestsList().Return(interests, nil)
}

func TestCreateUserInterestService(t *testing.T) {
	userInterest := model.UserInterests{
		UserID:         1,
		InterestID:     1,
		InterestValues: "Sleeping Well",
	}
	userInterests := make([]model.UserInterests, 0)
	userInterests = append(userInterests, userInterest)

	ctrl := gomock.NewController(t)
	ctrl.Finish()

	mockInterests := mocks.NewMockInterestsDao(ctrl)
	mockInterests.EXPECT().CreateUserInterestDao(gomock.Eq(userInterests)).Return(userInterests, nil).AnyTimes()
}

func TestGetUserInterestsService(t *testing.T) {
	userInterest := model.UserInterests{
		UserID:         1,
		InterestID:     1,
		InterestValues: "Sleeping Well",
	}
	userInterests := make([]model.UserInterests, 0)
	userInterests = append(userInterests, userInterest)

	ctrl := gomock.NewController(t)
	ctrl.Finish()

	mockInterests := mocks.NewMockInterestsDao(ctrl)
	mockInterests.EXPECT().GetUserInterestsDao(gomock.Eq(1)).Return(userInterests, nil).AnyTimes()
}

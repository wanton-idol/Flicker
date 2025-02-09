package tests

import (
	"encoding/json"
	"github.com/SuperMatch/model"
	elasticsearchPkg "github.com/SuperMatch/model/elasticSearch"
	"github.com/SuperMatch/pkg/db/dao/mocks"
	esMocks "github.com/SuperMatch/pkg/elasticseach/mocks"
	utils "github.com/SuperMatch/utilities"
	"github.com/golang/mock/gomock"
	"gorm.io/gorm"
	"testing"
	"time"
)

func userSearchProfileFunc() model.UserSearchProfile {
	return model.UserSearchProfile{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		UserId:        1,
		UserProfileId: 1,
		Gender:        "male,female",
		MinAge:        18,
		MaxAge:        30,
		Distance:      50,
		Snooze:        utils.ConvertStringToDateTime("2023-01-12 15:14:05"),
		HideMyName:    utils.ConvertBoolToBoolPointer(true),
	}
}
func TestCreateUserSearchProfile(t *testing.T) {
	var userSearchProfileMap = make(map[string]interface{})
	userSearchProfileMap["user_profile_id"] = 1
	userSearchProfileMap["user_id"] = 1
	userSearchProfileMap["min_age"] = 18
	userSearchProfileMap["max_age"] = 30
	userSearchProfileMap["distance"] = 50

	userSearchProfileMap["gender"] = "male,female"
	userSearchProfileMap["snooze"] = utils.ConvertStringToDateTime("2023-01-12 15:14:05")
	userSearchProfileMap["hide_my_name"] = utils.ConvertBoolToBoolPointer(true)
	userSearchProfile := userSearchProfileFunc()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserSearchProfileDao := mocks.NewMockUserSearchProfileRepository(ctrl)
	mockUserSearchProfileDao.EXPECT().CreateUserSearchProfile(gomock.Any(), gomock.Eq(userSearchProfileMap)).
		Return(userSearchProfile, nil).
		AnyTimes()

	searchProfileES := elasticsearchPkg.UserSearchProfile{
		MinAge:     18,
		MaxAge:     30,
		Distance:   50,
		Gender:     utils.ConvertStringToCustomGender("male,female"),
		Snooze:     utils.ConvertStringToDateTime("2023-01-12 15:14:05"),
		HideMyName: utils.ConvertBoolToBoolPointer(true),
	}

	userProfileES := elasticsearchPkg.UserProfile{
		UserSearchProfile: searchProfileES,
	}

	jsonString, err := json.Marshal(userProfileES)
	if err != nil {
		t.Error("error marshalling user profile for elasticsearch")
	}

	userProfileES.Id = 1
	mockES := esMocks.NewMockElasticSearchIndexer(ctrl)
	mockES.EXPECT().UpdateSearchProfile(gomock.Eq(userProfileES), gomock.Eq(jsonString)).
		Return(nil).AnyTimes()
}

func TestUpdateUserSearchProfile(t *testing.T) {
	var userSearchProfileMap = make(map[string]interface{})
	userSearchProfileMap["user_profile_id"] = 1
	userSearchProfileMap["user_id"] = 1
	userSearchProfileMap["min_age"] = 18
	userSearchProfileMap["max_age"] = 30
	userSearchProfileMap["distance"] = 50

	userSearchProfileMap["gender"] = "male,female"
	userSearchProfileMap["snooze"] = utils.ConvertStringToDateTime("2023-01-12 15:14:05")
	userSearchProfileMap["hide_my_name"] = true

	userSearchProfile := userSearchProfileFunc()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserSearchProfileDao := mocks.NewMockUserSearchProfileRepository(ctrl)
	mockUserSearchProfileDao.EXPECT().UpdateUserSearchProfile(gomock.Any(), gomock.Eq(userSearchProfileMap)).
		Return(userSearchProfile, nil).AnyTimes()

	searchProfileES := elasticsearchPkg.UserSearchProfile{
		MinAge:     18,
		MaxAge:     30,
		Distance:   50,
		Gender:     utils.ConvertStringToCustomGender("male,female"),
		Snooze:     utils.ConvertStringToDateTime("2023-01-12 15:14:05"),
		HideMyName: utils.ConvertBoolToBoolPointer(true),
	}

	userProfileES := elasticsearchPkg.UserProfile{
		UserSearchProfile: searchProfileES,
	}

	jsonString, err := json.Marshal(userProfileES)
	if err != nil {
		t.Error("error marshalling user profile for elasticsearch")
	}

	userProfileES.Id = 1
	mockES := esMocks.NewMockElasticSearchIndexer(ctrl)
	mockES.EXPECT().UpdateSearchProfile(gomock.Eq(userProfileES), gomock.Eq(jsonString)).
		Return(nil).AnyTimes()

}

func TestFindByProfileId(t *testing.T) {
	userSearchProfile := userSearchProfileFunc()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserSearchProfileDao := mocks.NewMockUserSearchProfileRepository(ctrl)
	mockUserSearchProfileDao.EXPECT().FindByProfileId(gomock.Any(), gomock.Eq(1)).
		Return(userSearchProfile, nil).
		AnyTimes()
}

func TestFindByUserId(t *testing.T) {
	userSearchProfile := userSearchProfileFunc()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserSearchProfileDao := mocks.NewMockUserSearchProfileRepository(ctrl)
	mockUserSearchProfileDao.EXPECT().FindByUserId(gomock.Any(), gomock.Eq(1)).
		Return(userSearchProfile, nil).
		AnyTimes()
}

func TestUpdateAdvancedFilters(t *testing.T) {
	advancedFilterDB := model.AdvancedFilter{
		UserId:            1,
		UserProfileId:     1,
		IsProfileVerified: utils.ConvertBoolToBoolPointer(true),
		IsOnline:          utils.ConvertBoolToBoolPointer(true),
		Height:            utils.ConvertSliceToCommaSeparatedString([]int64{130, 160}),
		Exercise:          "everyday, sometimes",
		Religion:          "hindu, jain",
		Education:         "college, graduate",
		Occupation:        "student, engineer",
		MaritalStatus:     "single, married",
		Drink:             "never, regular",
		Smoke:             "socially, never",
		IncognitoMode:     utils.ConvertBoolToBoolPointer(true),
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAdvFil := mocks.NewMockAdvancedFilterRepository(ctrl)
	mockAdvFil.EXPECT().UpdateAdvancedFilter(gomock.Eq(advancedFilterDB)).
		Return(model.AdvancedFilter{}, nil).
		AnyTimes()

	advancedFilterES := elasticsearchPkg.AdvancedFilter{
		IsProfileVerified: utils.ConvertBoolToBoolPointer(true),
		IsOnline:          utils.ConvertBoolToBoolPointer(true),
		Height:            []int64{130, 160},
		Exercise:          utils.ConvertStringToCustomExercise("everyday, sometimes"),
		Religion:          utils.ConvertStringToCustomReligion("hindu, jain"),
		EducationLevel:    utils.ConvertStringToCustomEducationLevel("college, graduate"),
		Occupation:        utils.ConvertStringToCustomOccupation("student, engineer"),
		MaritalStatus:     utils.ConvertStringToCustomMaritalStatus("single, married"),
		Drink:             utils.ConvertStringToCustomDrink("never, regular"),
		Smoke:             utils.ConvertStringToCustomSmoke("socially, never"),
		IncognitoMode:     utils.ConvertBoolToBoolPointer(true),
	}

	searchProfile := elasticsearchPkg.UserProfile{
		Id: 1,
		UserSearchProfile: elasticsearchPkg.UserSearchProfile{
			AdvancedFilter: advancedFilterES,
		},
	}

	jsonString, err := json.Marshal(searchProfile)
	if err != nil {
		t.Error("error marshalling search profile for elasticsearch")
	}

	mockES := esMocks.NewMockElasticSearchIndexer(ctrl)
	mockES.EXPECT().UpdateSearchProfile(gomock.Eq(searchProfile), gomock.Eq(jsonString)).
		Return(nil).
		AnyTimes()

}

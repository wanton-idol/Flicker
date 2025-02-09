package dao

import (
	"github.com/SuperMatch/model"
	"github.com/SuperMatch/pkg/db"
	"github.com/SuperMatch/zapLogger"
	"gorm.io/gorm"
)

type FiltersDao interface {
	GetExerciseFilterList() ([]model.ExerciseFilter, error)
	GetStarSignFilterList() ([]model.StarSignFilter, error)
	GetEducationFilterList() ([]model.EducationFilter, error)
	GetDrinkFilterList() ([]model.DrinkFilter, error)
	GetSmokeFilterList() ([]model.SmokeFilter, error)
	GetLookingForFilterList() ([]model.LookingForFilter, error)
	GetReligionFilterList() ([]model.ReligionFilter, error)
	GetPoliticsLikesFilterList() ([]model.PoliticsLikesFilter, error)
	GetChildrenFilterList() ([]model.ChildrenFilter, error)
}

type FiltersDaoImpl struct {
	Connection gorm.DB
}

func NewFiltersDaoImpl() *FiltersDaoImpl {
	return &FiltersDaoImpl{Connection: *db.GlobalOrm}
}

func (i *FiltersDaoImpl) GetExerciseFilterList() ([]model.ExerciseFilter, error) {
	var exerciseFilter []model.ExerciseFilter

	err := i.Connection.Table("exercise_filter").Find(&exerciseFilter)
	if err.Error != nil {
		zapLogger.Logger.Error("[GetExerciseFilterList] error in retrieving data from exercise_filter table")
		return exerciseFilter, err.Error
	}

	return exerciseFilter, nil
}

func (i *FiltersDaoImpl) GetStarSignFilterList() ([]model.StarSignFilter, error) {
	var starSignFilter []model.StarSignFilter

	err := i.Connection.Table("star_sign_filter").Find(&starSignFilter)
	if err.Error != nil {
		zapLogger.Logger.Error("[GetStarSignFilterList] error in retrieving data from star_sign_filter table")
		return starSignFilter, err.Error
	}

	return starSignFilter, nil
}

func (i *FiltersDaoImpl) GetEducationFilterList() ([]model.EducationFilter, error) {
	var educationFilter []model.EducationFilter

	err := i.Connection.Table("education_filter").Find(&educationFilter)
	if err.Error != nil {
		zapLogger.Logger.Error("[GetEducationFilterList] error in retrieving data from education_filter table")
		return educationFilter, err.Error
	}

	return educationFilter, nil
}

func (i *FiltersDaoImpl) GetDrinkFilterList() ([]model.DrinkFilter, error) {
	var drinkFilter []model.DrinkFilter

	err := i.Connection.Table("drink_filter").Find(&drinkFilter)
	if err.Error != nil {
		zapLogger.Logger.Error("[GetDrinkFilterList] error in retrieving data from drink_filter table")
		return drinkFilter, err.Error
	}

	return drinkFilter, nil
}

func (i *FiltersDaoImpl) GetSmokeFilterList() ([]model.SmokeFilter, error) {
	var smokeFilter []model.SmokeFilter

	err := i.Connection.Table("smoke_filter").Find(&smokeFilter)
	if err.Error != nil {
		zapLogger.Logger.Error("[GetSmokeFilterList] error in retrieving data from smoke_filter table")
		return smokeFilter, err.Error
	}

	return smokeFilter, nil
}

func (i *FiltersDaoImpl) GetLookingForFilterList() ([]model.LookingForFilter, error) {
	var lookingForFilter []model.LookingForFilter

	err := i.Connection.Table("looking_for_filter").Find(&lookingForFilter)
	if err.Error != nil {
		zapLogger.Logger.Error("[GetLookingForFilterList] error in retrieving data from looking_for_filter table")
		return lookingForFilter, err.Error
	}

	return lookingForFilter, nil
}

func (i *FiltersDaoImpl) GetReligionFilterList() ([]model.ReligionFilter, error) {
	var religionFilter []model.ReligionFilter

	err := i.Connection.Table("religion_filter").Find(&religionFilter)
	if err.Error != nil {
		zapLogger.Logger.Error("[GetReligionFilterList] error in retrieving data from religion_filter table")
		return religionFilter, err.Error
	}

	return religionFilter, nil
}

func (i *FiltersDaoImpl) GetPoliticsLikesFilterList() ([]model.PoliticsLikesFilter, error) {
	var politicsLikesFilter []model.PoliticsLikesFilter

	err := i.Connection.Table("politics_likes_filter").Find(&politicsLikesFilter)
	if err.Error != nil {
		zapLogger.Logger.Error("[GetPoliticsLikesFilterList] error in retrieving data from politics_likes_filter table")
		return politicsLikesFilter, err.Error
	}

	return politicsLikesFilter, nil
}

func (i *FiltersDaoImpl) GetChildrenFilterList() ([]model.ChildrenFilter, error) {
	var childrenFilter []model.ChildrenFilter

	err := i.Connection.Table("children_filter").Find(&childrenFilter)
	if err.Error != nil {
		zapLogger.Logger.Error("[GetChildrenFilterList] error in retrieving data from children_filter table")
		return childrenFilter, err.Error
	}

	return childrenFilter, nil
}

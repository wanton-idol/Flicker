package model

import (
	"gorm.io/gorm"
	"time"
)

type InterestsCategory struct {
	SelfCare               []string `json:"self_care"`
	Sports                 []string `json:"sports"`
	Creativity             []string `json:"creativity"`
	GoingOut               []string `json:"going_out"`
	FilmAndTV              []string `json:"film_and_tv"`
	StayingIn              []string `json:"staying_in"`
	Reading                []string `json:"reading"`
	Music                  []string `json:"music"`
	FoodAndDrink           []string `json:"food_and_drink"`
	Travelling             []string `json:"travelling"`
	Pets                   []string `json:"pets"`
	ValuesAndTraits        []string `json:"values_and_traits"`
	PlutoValuesAndAllyship []string `json:"pluto_values_and_allyship"`
}

type Interests struct {
	ID             int        `json:"id" gorm:"column:id;primaryKey"`
	InterestID     int        `json:"interest_id" gorm:"column:interest_id"`
	InterestValues string     `json:"interest_values" gorm:"column:interest_values"`
	DeletedAt      *time.Time `json:"deleted_at" gorm:"column:deleted_at;default:null"`
	Emoticon       string     `json:"emoticon" gorm:"column:emoticon"`
}

type UserInterests struct {
	gorm.Model
	UserID         int    `json:"user_id" gorm:"column:user_id"`
	InterestID     int    `json:"interest_id" gorm:"column:interest_id"`
	InterestValues string `json:"interest_values" gorm:"column:interest_values"`
}

type InterestsResponse struct {
	InterestValue string `json:"interest_value"`
	Emoticon      string `json:"emoticon"`
}

type InterestsListResponse struct {
	SelfCare               []InterestsResponse `json:"self_care"`
	Sports                 []InterestsResponse `json:"sports"`
	Creativity             []InterestsResponse `json:"creativity"`
	GoingOut               []InterestsResponse `json:"going_out"`
	FilmAndTV              []InterestsResponse `json:"film_and_tv"`
	StayingIn              []InterestsResponse `json:"staying_in"`
	Reading                []InterestsResponse `json:"reading"`
	Music                  []InterestsResponse `json:"music"`
	FoodAndDrink           []InterestsResponse `json:"food_and_drink"`
	Travelling             []InterestsResponse `json:"travelling"`
	Pets                   []InterestsResponse `json:"pets"`
	ValuesAndTraits        []InterestsResponse `json:"values_and_traits"`
	PlutoValuesAndAllyship []InterestsResponse `json:"pluto_values_and_allyship"`
}

type InterestData struct {
	InterestDetails []InterestDetails `json:"interest_details"`
}

type InterestDetails struct {
	InterestID     int    `json:"interest_id"`
	InterestValues string `json:"interest_values"`
}

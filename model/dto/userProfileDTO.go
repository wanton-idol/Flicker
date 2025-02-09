package dto

import (
	"time"

	elasticsearchPkg "github.com/SuperMatch/model/elasticSearch"
)

type UserProfile struct {
	Id                int                                 `json:"id,omitempty"`
	UserId            int                                 `json:"user_id,omitempty"`
	FirstName         *string                             `json:"first_name,omitempty"`
	LastName          *string                             `json:"last_name,omitempty"`
	IsVerified        *bool                               `json:"is_verified,omitempty"`
	IsPremium         *bool                               `json:"is_premium,omitempty"`
	DateOfBirth       *elasticsearchPkg.JsonBirthDate     `json:"date_of_birth,omitempty"`
	Gender            *elasticsearchPkg.Gender            `json:"gender,omitempty"`
	SexualOrientation *elasticsearchPkg.SexualOrientation `json:"sexual_orientation,omitempty"`
	Latitude          *float64                            `json:"latitude,omitempty"`
	Longitude         *float64                            `json:"longitude,omitempty"`
	Education         *elasticsearchPkg.Education         `json:"education,omitempty"`
	Occupation        *elasticsearchPkg.Occupation        `json:"occupation,omitempty"`
	MaritalStatus     *elasticsearchPkg.MaritalStatus     `json:"marital_status,omitempty"`
	Religion          *elasticsearchPkg.Religion          `json:"religion,omitempty"`
	Height            *int                                `json:"height,omitempty"`
	Weight            *int                                `json:"weight,omitempty"`
	LookingFor        *elasticsearchPkg.LookingFor        `json:"looking_for,omitempty"`
	Exercise          *elasticsearchPkg.Exercise          `json:"exercise,omitempty"`
	Drink             *elasticsearchPkg.Drink             `json:"drink,omitempty"`
	Smoke             *elasticsearchPkg.Smoke             `json:"smoke,omitempty"`
	About             *string                             `json:"about,omitempty"`
	Pronoun           *string                             `json:"pronoun,omitempty"`
	//Images            []UserImages                       `json:"images,omitempty"`
	//Nudges            []UserNudgeProfile                 `json:"questions,omitempty"`
	//SearchProfile     UserSearchProfile                  `json:"userSearchProfile,omitempty"`
}

type UserImages struct {
	Id    int    `json:"id,omitempty"`
	URL   string `json:"url,omitempty"`
	Order int    `json:"order,omitempty"`
}

type UserNudgeProfile struct {
	Id       int    `json:"id,omitempty"`
	Question string `json:"question,omitempty"`
	Answer   string `json:"answer,omitempty"`
	Order    int    `json:"order,omitempty"`
}

type UserSearchProfile struct {
	Id         int                       `json:"id,omitempty"`
	Gender     []elasticsearchPkg.Gender `json:"gender,omitempty"`
	MinAge     int                       `json:"min_age,omitempty"`
	MaxAge     int                       `json:"max_age,omitempty"`
	Distance   int                       `json:"distance,omitempty"`
	Language   []string                  `json:"language,omitempty"`
	Snooze     *time.Time                `json:"snooze,omitempty"`
	HideMyName *bool                     `json:"hide_my_name,omitempty"`
}

type AdvancedFilter struct {
	IsProfileVerified  *bool                             `json:"is_profile_verified"`
	IsOnline           *bool                             `json:"is_online"`
	Height             []int64                           `json:"height"`
	Exercise           []elasticsearchPkg.Exercise       `json:"exercise"`
	Religion           []elasticsearchPkg.Religion       `json:"religion"`
	EducationLevel     []elasticsearchPkg.EducationLevel `json:"education"`
	Occupation         []elasticsearchPkg.Occupation     `json:"occupation"`
	MaritalStatus      []elasticsearchPkg.MaritalStatus  `json:"marital_status"`
	Drink              []elasticsearchPkg.Drink          `json:"drink"`
	Smoke              []elasticsearchPkg.Smoke          `json:"smoke"`
	IncognitoMode      *bool                             `json:"incognito_mode"`
	StarSign           []string                          `json:"star_sign"`
	PoliticsLikes      []string                          `json:"politics_likes"`
	HaveOrWantChildren []string                          `json:"have_or_want_children"`
	LookingFor         []string                          `json:"looking_for"`
}

type Education struct {
	Level   string `json:"education_level,omitempty"`
	College string `json:"college,omitempty"`
}

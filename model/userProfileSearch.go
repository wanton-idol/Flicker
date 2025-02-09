package model

import (
	"database/sql/driver"
	"errors"
	elasticsearchPkg "github.com/SuperMatch/model/elasticSearch"
	"gorm.io/gorm"
	"strings"
	"time"
)

// UserProfileSearchDTO is a struct that contains the search criteria for the user profile

func (UserProfile) TableName() string {
	return "user_profile"
}

type UserProfile struct {
	gorm.Model
	UserId            int                                 `json:"user_id" gorm:"column:user_id"`
	FirstName         *string                             `json:"first_name" gorm:"column:first_name"`
	LastName          *string                             `json:"last_name" gorm:"column:last_name"`
	IsVerified        bool                                `json:"is_verified" gorm:"column:is_verified"`
	IsPremium         bool                                `json:"is_premium" gorm:"column:is_premium"`
	DateOfBirth       *elasticsearchPkg.JsonBirthDate     `json:"date_of_birth" gorm:"column:date_of_birth"`
	Gender            *elasticsearchPkg.Gender            `json:"gender" gorm:"column:gender"`
	SexualOrientation *elasticsearchPkg.SexualOrientation `json:"sexual_orientation" gorm:"column:sexual_orientation"`
	Latitude          *float64                            `json:"latitude" gorm:"column:latitude"`
	Longitude         *float64                            `json:"longitude" gorm:"column:longitude"`
	EducationLevel    *elasticsearchPkg.EducationLevel    `json:"education_level" gorm:"column:education_level"`
	College           *string                             `json:"college" gorm:"column:college"`
	Occupation        *elasticsearchPkg.Occupation        `json:"occupation" gorm:"column:occupation"`
	MaritalStatus     *elasticsearchPkg.MaritalStatus     `json:"marital_status" gorm:"column:marital_status"`
	Religion          *elasticsearchPkg.Religion          `json:"religion" gorm:"column:religion"`
	Height            *int                                `json:"height" gorm:"column:height"`
	Weight            *int                                `json:"weight" gorm:"column:weight"`
	LookingFor        *elasticsearchPkg.LookingFor        `json:"looking_for" gorm:"column:looking_for"`
	Exercise          *elasticsearchPkg.Exercise          `json:"exercise" gorm:"column:exercise"`
	Drink             *elasticsearchPkg.Drink             `json:"drink" gorm:"column:drink"`
	Smoke             *elasticsearchPkg.Smoke             `json:"smoke" gorm:"column:smoke"`
	About             *string                             `json:"about" gorm:"column:about"`
	Pronoun           *string                             `json:"pronoun" gorm:"column:pronoun"`

	//Images            []UserMedia        `json:"images" gorm:"foreignKey:user_profile_id;references:id"`
	//Nudges            []UserNudgeProfile `json:"questions" gorm:"foreignKey:user_profile_id;references:id"`
	// UserSearchProfile UserSearchProfile  `json:"userSearchProfile" gorm:"foreignKey:user_profile_id;references:id"`
}

func (UserNudgeProfile) TableName() string {
	return "user_nudge"
}

type UserNudgeProfile struct {
	gorm.Model
	UserProfileId *int   `json:"user_profile_id" gorm:"column:user_profile_id"`
	Question      string `json:"question" gorm:"column:question"`
	Answer        string `json:"answer" gorm:"column:answer"`
	Order         *int   `json:"order" gorm:"column:order"`
}

func (UserSearchProfile) TableName() string {
	return "user_search_profile"
}

type UserSearchProfile struct {
	gorm.Model
	UserId        int        `json:"user_id" gorm:"column:user_id"`
	UserProfileId int        `json:"user_profile_id" gorm:"column:user_profile_id"`
	Gender        string     `json:"gender" gorm:"column:gender"`
	MinAge        int        `json:"min_age" gorm:"column:min_age"`
	MaxAge        int        `json:"max_age" gorm:"column:max_age"`
	Distance      int        `json:"distance" gorm:"column:distance"`
	Snooze        *time.Time `json:"snooze" gorm:"column:snooze"`
	HideMyName    *bool      `json:"hide_my_name" gorm:"column:hide_my_name"`
	// Language      Language `json:"language" gorm:"column:language"`
	//AdvancedFilter AdvancedFilter `json:"advanced_filter" gorm:"foreignKey:UserSearchProfileId;references:ID"`
}

func (AdvancedFilter) TableName() string {
	return "advanced_filters"
}

type AdvancedFilter struct {
	gorm.Model
	UserId             int    `json:"user_id" gorm:"column:user_id"`
	UserProfileId      int    `json:"user_profile_id" gorm:"user_profile_id"`
	IsProfileVerified  *bool  `json:"is_profile_verified" gorm:"column:is_profile_verified"`
	IsOnline           *bool  `json:"is_online" gorm:"column:is_online"`
	Height             string `json:"height" gorm:"column:height"`
	Exercise           string `json:"exercise" gorm:"column:exercise"`
	Religion           string `json:"religion" gorm:"column:religion"`
	Education          string `json:"education" gorm:"column:education"`
	Occupation         string `json:"occupation" gorm:"column:occupation"`
	MaritalStatus      string `json:"marital_status" gorm:"column:marital_status"`
	Drink              string `json:"drink" gorm:"column:drink"`
	Smoke              string `json:"smoke" gorm:"column:smoke"`
	IncognitoMode      *bool  `json:"incognito_mode" gorm:"column:incognito_mode"`
	StarSign           string `json:"star_sign"`
	PoliticsLikes      string `json:"politics_likes"`
	HaveOrWantChildren string `json:"have_or_want_children"`
	LookingFor         string `json:"looking_for"`
}

type Gender []elasticsearchPkg.Gender

// convert string to array of custom datatype elasticsearchPkg.Gender
func (t *Gender) Scan(value interface{}) error {
	if value == nil {
		*t = []elasticsearchPkg.Gender{}
		return nil
	}
	str, ok := value.(string)
	if !ok {
		return errors.New("type assertion to []string failed")
	}
	if str == "" {
		*t = []elasticsearchPkg.Gender{}
		return nil
	}
	str = strings.ToLower(str)
	if str == "male" {
		*t = []elasticsearchPkg.Gender{elasticsearchPkg.MALE}
	} else if str == "female" {
		*t = []elasticsearchPkg.Gender{elasticsearchPkg.FEMALE}
	} else {
		return errors.New("type assertion to []string failed")
	}
	return nil
}

// Implement Valuer interface
func (t Gender) Value() (driver.Value, error) {
	if len(t) == 0 {
		return nil, nil
	}

	// convert to string by iterating 	through
	var str string
	for _, v := range t {
		str += string(v)
	}
	return str, nil
}

type Language []string

func (t *Language) Scan(value interface{}) error {
	val, _ := value.([]byte)
	// if !ok {
	// 	return errors.New(fmt.Sprint("wrong type", value))
	// }

	*t = Language(strings.Split(string(val), ","))

	return nil
}

// Implement Valuer interface
func (t Language) Value() (driver.Value, error) {
	//this check is here if you don't want to save an empty string
	if len(t) == 0 {
		return nil, nil
	}

	return []byte(strings.Join(t, ",")), nil
}

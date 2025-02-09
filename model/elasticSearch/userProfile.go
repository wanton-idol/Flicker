package elasticsearchPkg

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Education struct {
	Level   *EducationLevel `json:"education_level,omitempty"`
	College *string         `json:"college,omitempty"`
}

type EducationLevel string

const (
	HighSchool          EducationLevel = "high school"
	VocationalSchool    EducationLevel = "vocational school"
	InCollege           EducationLevel = "in college"
	UnderGraduateDegree EducationLevel = "undergraduate degree"
	InGradSchool        EducationLevel = "in grad school"
	GraduateDegree      EducationLevel = "graduate degree"
)

func (o *EducationLevel) UnmarshalJSON(b []byte) error {
	type OT EducationLevel
	var r *OT = (*OT)(o)
	err := json.Unmarshal(b, &r)
	if err != nil {
		return err
	}

	switch *o {
	case HighSchool, VocationalSchool, InCollege, UnderGraduateDegree, InGradSchool, GraduateDegree:
		return nil
	}
	return errors.New("invalid education level type")
}

type Occupation string

const (
	Student      Occupation = "student"
	Doctor       Occupation = "doctor"
	Engineer     Occupation = "engineer"
	Entrepreneur Occupation = "entrepreneur"
	Other        Occupation = "other"
)

func (o *Occupation) UnmarshalJSON(b []byte) error {
	// Unmarshal into the pointer receiver, which is a pointer to the enum
	type OT Occupation
	var r *OT = (*OT)(o)
	err := json.Unmarshal(b, &r)
	if err != nil {
		return err
	}

	switch *o {
	case Student, Doctor, Engineer, Entrepreneur, Other:
		return nil
	}
	return errors.New("invalid occupation type")
}

type MaritalStatus string

const (
	Single   MaritalStatus = "single"
	Married  MaritalStatus = "married"
	Divorced MaritalStatus = "divorced"
	Widowed  MaritalStatus = "widowed"
)

func (o *MaritalStatus) UnmarshalJSON(b []byte) error {
	// Unmarshal into the pointer receiver, which is a pointer to the enum
	type OT MaritalStatus
	var r *OT = (*OT)(o)
	err := json.Unmarshal(b, &r)
	if err != nil {
		return err
	}
	switch *o {
	case Single, Married, Divorced, Widowed:
		return nil
	}
	return errors.New("invalid marital status type")
}

type Religion string

const (
	Agnostic       Religion = "agnostic"
	Atheist        Religion = "atheist"
	Catholic       Religion = "catholic"
	Hindu          Religion = "hindu"
	Muslim         Religion = "muslim"
	Christian      Religion = "christian"
	Jain           Religion = "jain"
	Sikh           Religion = "sikh"
	Buddhist       Religion = "buddhist"
	Jewish         Religion = "jewish"
	Mormon         Religion = "mormon"
	LatterDaySaint Religion = "latter-day saint"
	Zoroastrian    Religion = "zoroastrian"
	Spiritual      Religion = "spiritual"
)

func (o *Religion) UnmarshalJSON(b []byte) error {
	// Unmarshal into the pointer receiver, which is a pointer to the enum
	type OT Religion
	var r *OT = (*OT)(o)
	err := json.Unmarshal(b, &r)
	if err != nil {
		return err
	}

	switch *o {
	case Agnostic, Atheist, Catholic, Hindu, Muslim, Christian, Jain, Sikh, Buddhist, Jewish, Mormon, LatterDaySaint, Zoroastrian, Spiritual:
		return nil
	}
	return errors.New("invalid religion type")
}

type Gender string

const (
	MALE   Gender = "male"
	FEMALE Gender = "female"
	BINARY Gender = "binary"
)

func (o *Gender) UnmarshalJSON(b []byte) error {
	type OT Gender
	var r *OT = (*OT)(o)
	err := json.Unmarshal(b, &r)
	if err != nil {
		return err
	}

	switch *o {
	case MALE, FEMALE, BINARY:
		return nil
	}
	return errors.New("invalid gender type")
}

type SexualOrientation string

const (
	Heterosexual SexualOrientation = "heterosexual"
	Homosexual   SexualOrientation = "homosexual"
	Bisexual     SexualOrientation = "bisexual"
	Asexual      SexualOrientation = "asexual"
	Pansexual    SexualOrientation = "pansexual"
	Demisexual   SexualOrientation = "demisexual"
)

func (o *SexualOrientation) UnmarshalJSON(b []byte) error {
	type OT SexualOrientation
	var r *OT = (*OT)(o)
	err := json.Unmarshal(b, &r)
	if err != nil {
		return err
	}

	switch *o {
	case Heterosexual, Homosexual, Bisexual, Asexual, Pansexual, Demisexual:
		return nil
	}
	return errors.New("invalid SexualOrientation type")
}

type LookingFor string

const (
	Relationship LookingFor = "relationship"
	Marriage     LookingFor = "marriage"
	Casual       LookingFor = "casual"
	NotKnownYet  LookingFor = "NotKnownYet"
)

func (o *LookingFor) UnmarshalJSON(b []byte) error {
	type OT LookingFor
	var r *OT = (*OT)(o)
	err := json.Unmarshal(b, &r)
	if err != nil {
		return err
	}

	switch *o {
	case Relationship, Marriage, Casual, NotKnownYet:
		return nil
	}
	return errors.New("invalid SexualOrientation type")
}

type JsonBirthDate time.Time

func (j JsonBirthDate) Value() (driver.Value, error) {
	t := time.Time(j)

	formatted := t.Format("2006-01-02")

	return []byte(formatted), nil
}

func (t *JsonBirthDate) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = JsonBirthDate(value)
		return nil
	}
	return fmt.Errorf("failed to scan JsonBirthDate")
}

func (j *JsonBirthDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*j = JsonBirthDate(t)
	return nil
}

func (j *JsonBirthDate) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", time.Time(*j).Format("2006-01-02"))), nil
}

type Exercise string

const (
	Exercise_Active      Exercise = "active"
	Exercise_Sometimes   Exercise = "sometimes"
	Exercise_AlmostNever Exercise = "almost never"
)

func (o *Exercise) UnmarshalJSON(b []byte) error {
	type OT Exercise
	var r *OT = (*OT)(o)
	err := json.Unmarshal(b, &r)
	if err != nil {
		return err
	}

	switch *o {
	case Exercise_Active, Exercise_Sometimes, Exercise_AlmostNever:
		return nil
	}
	return errors.New("invalid exercise type")
}

type Smoke string

const (
	Smoke_Socially  Smoke = "socially"
	Smoke_Never     Smoke = "never"
	Smoke_Regularly Smoke = "regularly"
)

func (o *Smoke) UnmarshalJSON(b []byte) error {
	type OT Smoke
	var r *OT = (*OT)(o)
	err := json.Unmarshal(b, &r)
	if err != nil {
		return err
	}

	switch *o {
	case Smoke_Socially, Smoke_Never, Smoke_Regularly:
		return nil
	}
	return errors.New("invalid smoke type")
}

type Drink string

const (
	Drink_Frequently Drink = "frequently"
	Drink_Socially   Drink = "socially"
	Drink_Rarely     Drink = "rarely"
	Drink_Never      Drink = "never"
	Drink_Sober      Drink = "Sober"
)

func (o *Drink) UnmarshalJSON(b []byte) error {
	type OT Drink
	var r *OT = (*OT)(o)
	err := json.Unmarshal(b, &r)
	if err != nil {
		return err
	}

	switch *o {
	case Drink_Frequently, Drink_Socially, Drink_Rarely, Drink_Never, Drink_Sober:
		return nil
	}
	return errors.New("invalid drink type")
}

type UserProfile struct {
	Id                 int            `json:"id,omitempty"`
	UserId             int            `json:"user_id,omitempty"`
	FirstName          *string        `json:"first_name,omitempty"`
	LastName           *string        `json:"last_name,omitempty"`
	IsVerified         *bool          `json:"is_verified,omitempty"`
	IsPremium          *bool          `json:"is_premium,omitempty"`
	DateOfBirth        *JsonBirthDate `json:"date_of_birth,omitempty"`
	*Gender            `json:"gender,omitempty"`
	*SexualOrientation `json:"sexual_orientation,omitempty"`
	Location           []float64 `json:"location,omitempty"`
	*Education         `json:"education,omitempty"`
	*Occupation        `json:"occupation,omitempty"`
	*MaritalStatus     `json:"marital_status,omitempty"`
	*Religion          `json:"religion,omitempty"`
	Height             *int `json:"height,omitempty"`
	Weight             *int `json:"weight,omitempty"`
	*LookingFor        `json:"looking_for,omitempty"`
	*Exercise          `json:"exercise,omitempty"`
	*Drink             `json:"drink,omitempty"`
	*Smoke             `json:"smoke,omitempty"`
	About              *string            `json:"about,omitempty"`
	Images             []Image            `json:"images,omitempty"`
	Nudges             []UserNudgeProfile `json:"questions,omitempty"`
	Pronoun            *string            `json:"pronoun,omitempty"`
	UserSearchProfile  `json:"userSearchProfile,omitempty"`
}

type AdvancedFilter struct {
	IsProfileVerified  *bool            `json:"is_profile_verified,omitempty"`
	IsOnline           *bool            `json:"is_online,omitempty"`
	Height             []int64          `json:"height,omitempty"`
	Exercise           []Exercise       `json:"exercise,omitempty"`
	Religion           []Religion       `json:"religion,omitempty"`
	EducationLevel     []EducationLevel `json:"education,omitempty"`
	Occupation         []Occupation     `json:"occupation,omitempty"`
	MaritalStatus      []MaritalStatus  `json:"marital_status,omitempty"`
	Drink              []Drink          `json:"drink,omitempty"`
	Smoke              []Smoke          `json:"smoke,omitempty"`
	IncognitoMode      *bool            `json:"incognito_mode,omitempty"`
	StarSign           []string         `json:"star_sign,omitempty"`
	PoliticsLikes      []string         `json:"politics_likes,omitempty"`
	HaveOrWantChildren []string         `json:"have_or_want_children,omitempty"`
	LookingFor         []string         `json:"looking_for,omitempty"`
}

type UserSearchProfile struct {
	Gender         []Gender   `json:"gender,omitempty"`
	MinAge         int        `json:"min_age,omitempty"`
	MaxAge         int        `json:"max_age,omitempty"`
	Distance       int        `json:"distance,omitempty" default:"10"`
	Language       []string   `json:"language,omitempty"`
	Snooze         *time.Time `json:"snooze,omitempty"`
	HideMyName     *bool      `json:"hide_my_name,omitempty"`
	AdvancedFilter `json:"advanced_filter,omitempty"`
}

type UserNudgeProfile struct {
	Id       int    `json:"id,omitempty"`
	Question string `json:"question,omitempty"`
	Answer   string `json:"answer,omitempty"`
	Order    int    `json:"order,omitempty"`
}

type Image struct {
	Id    int    `json:"id,omitempty"`
	URL   string `json:"url,omitempty"`
	Order int    `json:"order,omitempty"`
}

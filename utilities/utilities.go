package utilities

import (
	"bytes"
	"fmt"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/SuperMatch/model"

	"github.com/SuperMatch/model/dto"
	elasticsearchPkg "github.com/SuperMatch/model/elasticSearch"
	"github.com/SuperMatch/zapLogger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ConvertIntToString(val int) string {
	return fmt.Sprintf("%v", val)
}

// CreateArrayPointer create pointer array of float
func CreateArrayPointer(x, y *float64) []float64 {
	if x == nil || y == nil {
		return nil
	}
	return []float64{*x, *y}
}

func ConvertCustomTypeGenderToString(gender []elasticsearchPkg.Gender) (string, error) {
	if len(gender) == 0 {
		return "", nil
	}
	if len(gender) == 1 {
		return string(gender[0]), nil
	}
	if len(gender) > 1 {
		genderString := ""
		for i := range gender {
			genderString += string(gender[i])
			if i < len(gender)-1 {
				genderString += ", "
			}
		}
		return genderString, nil
	}
	return "", nil
}

func ConvertStringToCustomGender(gender string) []elasticsearchPkg.Gender {

	if len(gender) == 0 {
		return nil
	}
	genderArray := strings.Split(gender, ", ")
	var genderSlice []elasticsearchPkg.Gender
	for i := range genderArray {
		genderSlice = append(genderSlice, elasticsearchPkg.Gender(genderArray[i]))
	}
	return genderSlice
}

// ConvertSliceToCommaSeparatedString Covert slice to comma separated string
func ConvertSliceToCommaSeparatedString(arr []int64) string {

	var buffer bytes.Buffer
	for i := range arr {
		buffer.WriteString(strconv.FormatInt(arr[i], 10))
		if i < len(arr)-1 {
			buffer.WriteString(", ")
		}
	}
	return buffer.String()
}

// ConvertCustomTypeExerciseToString Convert Custom Type Exercise to Comma separated String
func ConvertCustomTypeExerciseToString(exercises []elasticsearchPkg.Exercise) string {
	if len(exercises) == 0 {
		return ""
	}
	if len(exercises) == 1 {
		return string(exercises[0])
	}
	if len(exercises) > 1 {
		exerciseString := ""
		for i := range exercises {
			exerciseString += string(exercises[i])
			if i < len(exercises)-1 {
				exerciseString += ", "
			}
		}
		return exerciseString
	}
	return ""
}

// ConvertStringToCustomExercise Convert String to Custom Type Exercise
func ConvertStringToCustomExercise(exercise string) []elasticsearchPkg.Exercise {

	if len(exercise) == 0 {
		return nil
	}
	exerciseArray := strings.Split(exercise, ", ")
	var exerciseSlice []elasticsearchPkg.Exercise
	for i := range exerciseArray {
		exerciseSlice = append(exerciseSlice, elasticsearchPkg.Exercise(exerciseArray[i]))
	}
	return exerciseSlice
}

// ConvertCustomTypeReligionToString Convert Custom Type Religion to Comma separated String
func ConvertCustomTypeReligionToString(religions []elasticsearchPkg.Religion) string {
	if len(religions) == 0 {
		return ""
	}
	if len(religions) == 1 {
		return string(religions[0])
	}
	if len(religions) > 1 {
		religionString := ""
		for i := range religions {
			religionString += string(religions[i])
			if i < len(religions)-1 {
				religionString += ", "
			}
		}
		return religionString
	}
	return ""
}

// ConvertStringToCustomReligion Convert String to Custom Type Religion
func ConvertStringToCustomReligion(religion string) []elasticsearchPkg.Religion {

	if len(religion) == 0 {
		return nil
	}
	religionArray := strings.Split(religion, ", ")
	var religionSlice []elasticsearchPkg.Religion
	for i := range religionArray {
		religionSlice = append(religionSlice, elasticsearchPkg.Religion(religionArray[i]))
	}
	return religionSlice
}

// ConvertCustomTypeEducationLevelToString Convert Custom Type EducationLevel to Comma separated String
func ConvertCustomTypeEducationLevelToString(educations []elasticsearchPkg.EducationLevel) string {
	if len(educations) == 0 {
		return ""
	}
	if len(educations) == 1 {
		return string(educations[0])
	}
	if len(educations) > 1 {
		educationString := ""
		for i := range educations {
			educationString += string(educations[i])
			if i < len(educations)-1 {
				educationString += ", "
			}
		}
		return educationString
	}
	return ""
}

// ConvertStringToCustomEducationLevel Convert String to Custom Type EducationLevel
func ConvertStringToCustomEducationLevel(education string) []elasticsearchPkg.EducationLevel {

	if len(education) == 0 {
		return nil
	}
	educationArray := strings.Split(education, ", ")
	var educationSlice []elasticsearchPkg.EducationLevel
	for i := range educationArray {
		educationSlice = append(educationSlice, elasticsearchPkg.EducationLevel(educationArray[i]))
	}
	return educationSlice
}

// ConvertCustomTypeMaritalStatusToString Convert Custom Type MaritalStatus to Comma separated String
func ConvertCustomTypeMaritalStatusToString(maritalStatus []elasticsearchPkg.MaritalStatus) string {
	if len(maritalStatus) == 0 {
		return ""
	}
	if len(maritalStatus) == 1 {
		return string(maritalStatus[0])
	}
	if len(maritalStatus) > 1 {
		maritalStatusString := ""
		for i := range maritalStatus {
			maritalStatusString += string(maritalStatus[i])
			if i < len(maritalStatus)-1 {
				maritalStatusString += ", "
			}
		}
		return maritalStatusString
	}
	return ""
}

// ConvertStringToCustomMaritalStatus Convert String to Custom Type MaritalStatus
func ConvertStringToCustomMaritalStatus(maritalStatus string) []elasticsearchPkg.MaritalStatus {

	if len(maritalStatus) == 0 {
		return nil
	}
	maritalStatusArray := strings.Split(maritalStatus, ", ")
	var maritalStatusSlice []elasticsearchPkg.MaritalStatus
	for i := range maritalStatusArray {
		maritalStatusSlice = append(maritalStatusSlice, elasticsearchPkg.MaritalStatus(maritalStatusArray[i]))
	}
	return maritalStatusSlice
}

// ConvertCustomTypeDrinkToString Convert Custom Type Drink to Comma separated String
func ConvertCustomTypeDrinkToString(drink []elasticsearchPkg.Drink) string {
	if len(drink) == 0 {
		return ""
	}
	if len(drink) == 1 {
		return string(drink[0])
	}
	if len(drink) > 1 {
		drinkString := ""
		for i := range drink {
			drinkString += string(drink[i])
			if i < len(drink)-1 {
				drinkString += ", "
			}
		}
		return drinkString
	}
	return ""
}

// ConvertStringToCustomDrink Convert String to Custom Type Drink
func ConvertStringToCustomDrink(drink string) []elasticsearchPkg.Drink {

	if len(drink) == 0 {
		return nil
	}
	drinkArray := strings.Split(drink, ", ")
	var drinkSlice []elasticsearchPkg.Drink
	for i := range drinkArray {
		drinkSlice = append(drinkSlice, elasticsearchPkg.Drink(drinkArray[i]))
	}
	return drinkSlice
}

// ConvertCustomTypeSmokeToString Convert Custom Type Smoke to Comma separated String
func ConvertCustomTypeSmokeToString(smoke []elasticsearchPkg.Smoke) string {
	if len(smoke) == 0 {
		return ""
	}
	if len(smoke) == 1 {
		return string(smoke[0])
	}
	if len(smoke) > 1 {
		smokeString := ""
		for i := range smoke {
			smokeString += string(smoke[i])
			if i < len(smoke)-1 {
				smokeString += ", "
			}
		}
		return smokeString
	}
	return ""
}

// ConvertStringToCustomSmoke Convert String to Custom Type Smoke
func ConvertStringToCustomSmoke(smoke string) []elasticsearchPkg.Smoke {

	if len(smoke) == 0 {
		return nil
	}
	smokeArray := strings.Split(smoke, ", ")
	var smokeSlice []elasticsearchPkg.Smoke
	for i := range smokeArray {
		smokeSlice = append(smokeSlice, elasticsearchPkg.Smoke(smokeArray[i]))
	}
	return smokeSlice
}

// ConvertCustomTypeOccupationToString Convert Custom Type Occupation to Comma separated String
func ConvertCustomTypeOccupationToString(occupation []elasticsearchPkg.Occupation) string {
	if len(occupation) == 0 {
		return ""
	}
	if len(occupation) == 1 {
		return string(occupation[0])
	}
	if len(occupation) > 1 {
		occupationString := ""
		for i := range occupation {
			occupationString += string(occupation[i])
			if i < len(occupation)-1 {
				occupationString += ", "
			}
		}
		return occupationString
	}
	return ""
}

// ConvertStringToCustomOccupation Convert String to Custom Type Occupation
func ConvertStringToCustomOccupation(occupation string) []elasticsearchPkg.Occupation {

	if len(occupation) == 0 {
		return nil
	}
	occupationArray := strings.Split(occupation, ", ")
	var occupationSlice []elasticsearchPkg.Occupation
	for i := range occupationArray {
		occupationSlice = append(occupationSlice, elasticsearchPkg.Occupation(occupationArray[i]))
	}
	return occupationSlice
}

func RootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Dir(d)
}
func ConvertStringToDateTime(s string) *time.Time {
	DateTime := "2006-01-02 15:04:05"
	t, err := time.Parse(DateTime, s)
	if err != nil {
		zapLogger.Logger.Error("error while converting string to date:", zap.Error(err))
	}

	return &t
}

func ConvertStringToStringPointer(s string) *string {
	return &s
}

func ConvertBoolToBoolPointer(b bool) *bool {
	return &b
}

func ConvertStringToCustomGenderPointer(s string) *elasticsearchPkg.Gender {
	gender := elasticsearchPkg.Gender(s)
	return &gender
}

func ConvertStringToCustomSexualOrientationPointer(s string) *elasticsearchPkg.SexualOrientation {
	sexualOrientation := elasticsearchPkg.SexualOrientation(s)
	return &sexualOrientation
}

func ConvertStringToCustomOccupationPointer(s string) *elasticsearchPkg.Occupation {
	occupation := elasticsearchPkg.Occupation(s)
	return &occupation
}

func ConvertStringToCustomMaritalStatusPointer(s string) *elasticsearchPkg.MaritalStatus {
	maritalStatus := elasticsearchPkg.MaritalStatus(s)
	return &maritalStatus
}

func ConvertStringToCustomReligionPointer(s string) *elasticsearchPkg.Religion {
	religion := elasticsearchPkg.Religion(s)
	return &religion
}

func ConvertIntToIntPointer(i int) *int {
	return &i
}

func ConvertStringToCustomLookingForPointer(s string) *elasticsearchPkg.LookingFor {
	lookingFor := elasticsearchPkg.LookingFor(s)
	return &lookingFor
}

func ConvertStringToCustomExercisePointer(s string) *elasticsearchPkg.Exercise {
	exercise := elasticsearchPkg.Exercise(s)
	return &exercise
}

func ConvertStringToCustomDrinkPointer(s string) *elasticsearchPkg.Drink {
	drink := elasticsearchPkg.Drink(s)
	return &drink
}

func ConvertStringToCustomSmokePointer(s string) *elasticsearchPkg.Smoke {
	smoke := elasticsearchPkg.Smoke(s)
	return &smoke
}

func ConvertStringToCustomDOBPointer(s time.Time) *elasticsearchPkg.JsonBirthDate {
	DOB := elasticsearchPkg.JsonBirthDate(s)
	return &DOB
}

func ConvertStringToCustomEduLevelPointer(s string) *elasticsearchPkg.EducationLevel {
	eduLevel := elasticsearchPkg.EducationLevel(s)
	return &eduLevel
}

func ConvertFloatToFloatPointer(f float64) *float64 {
	return &f
}

func ReadPaginationDataFromRequest(c *gin.Context) (model.Pagination, error) {

	queryParams := c.Request.URL.Query()
	pageNo := queryParams["page_no"]
	pageSize := queryParams["page_size"]
	sort := queryParams["sort"]
	sortBy := queryParams["sort_by"]

	page := model.Pagination{
		PageSize:   10,
		PageNumber: 1,
		Sort:       "asc",
		SortBy:     "created_at",
	}

	if len(pageNo) == 1 {
		pageNo, err := strconv.Atoi(pageNo[0])
		if err != nil {
			zapLogger.Logger.Error("error while converting string to int:", zap.Error(err))
			return page, err
		}
		page.PageNumber = pageNo
	}

	if len(pageSize) == 1 {
		pageLimit, err := strconv.Atoi(pageSize[0])
		if err != nil {
			zapLogger.Logger.Error("error while converting string to int:", zap.Error(err))
			return page, err
		}
		page.PageSize = pageLimit
	}

	if sort != nil && len(sort[0]) > 0 {
		page.Sort = sort[0]
	}
	if sortBy != nil && len(sortBy[0]) > 0 {
		page.SortBy = sortBy[0]
	}

	return page, nil
}

func ParseSearchEventFilters(filters *dto.EventFilterDTO) {

	if filters.Distance == 0 {
		filters.Distance = 2
	}

	if filters.StartDate == nil {
		time := time.Now()
		filters.StartDate = &time
	}

	if filters.EndDate == nil {
		time := time.Now().AddDate(0, 0, 7)
		filters.EndDate = &time
	}
}

func ConvertSliceOfStringToCommaSeparatedString(arr []string) string {
	str := strings.Join(arr[:], ", ")
	return str
}

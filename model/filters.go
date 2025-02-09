package model

type ExerciseFilter struct {
	DoTheyExercise string `json:"do_they_exercise"`
}

type StarSignFilter struct {
	StarSign string `json:"star_sign"`
}

type EducationFilter struct {
	Education string `json:"education"`
}

type DrinkFilter struct {
	DoTheyDrink string `json:"do_they_drink"`
}

type SmokeFilter struct {
	DoTheySmoke string `json:"do_they_smoke"`
}

type LookingForFilter struct {
	LookingFor string `json:"looking_for"`
}

type ReligionFilter struct {
	Religion string `json:"religion"`
}

type PoliticsLikesFilter struct {
	PoliticsLikes string `json:"politics_likes"`
}

type ChildrenFilter struct {
	HaveOrWantChildren string `json:"have_or_want_children"`
}

type Filters struct {
	Exercise      []string `json:"exercise"`
	StarSign      []string `json:"star_sign"`
	Education     []string `json:"education"`
	Drink         []string `json:"drink"`
	Smoke         []string `json:"smoke"`
	LookingFor    []string `json:"looking_for"`
	Religion      []string `json:"religion"`
	PoliticsLikes []string `json:"politics_likes"`
	Children      []string `json:"children"`
}

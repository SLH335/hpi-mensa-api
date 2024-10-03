package types

import "time"

type Format string

const (
	FormatJSON Format = "json"
	FormatHTML Format = "html"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

type Location struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Language struct {
	Id    int
	Short string
}

var (
	German Language = Language{
		Id:    1,
		Short: "de",
	}
	English Language = Language{
		Id:    2,
		Short: "en",
	}
)

type Model string
type AttributeModel Model

const (
	AdditivesModel    AttributeModel = "additives"
	AllergensModel    AttributeModel = "allergens"
	FeaturesModel     AttributeModel = "features"
	LocationsModel    Model          = "location"
	MealCategoryModel Model          = "mealCategory"
	MenuModel         Model          = "menu"
	OutletsModel      Model          = "outlet"
)

type MealCategory struct {
	Id       int       `json:"id"`
	Name     string    `json:"name"`
	NameDe   string    `json:",omitempty"`
	NameEn   string    `json:",omitempty"`
	Location *Location `json:"location,omitempty"`
}

func (category MealCategory) Translated(language Language) (translated MealCategory) {
	translated = category
	if language == German {
		translated.Name = category.NameDe
	} else {
		translated.Name = category.NameEn
	}
	translated.NameDe = ""
	translated.NameEn = ""
	return translated
}

type Meal struct {
	Id           int             `json:"id"`
	Name         string          `json:"name"`
	NameDe       string          `json:",omitempty"`
	NameEn       string          `json:",omitempty"`
	Category     MealCategory    `json:"category"`
	Date         time.Time       `json:"date"`
	StudentPrice float64         `json:"studentPrice"`
	GuestPrice   float64         `json:"guestPrice"`
	Nutrition    Nutrition       `json:"nutrition"`
	Additives    []MealAttribute `json:"additives"`
	Allergens    []MealAttribute `json:"allergens"`
	Features     []MealAttribute `json:"features"`
	CO2          MealCO2         `json:"co2"`
	Location     Location        `json:"location"`
}

type MealCO2 struct {
	Grams  int    `json:"grams"`
	Rating string `json:"rating"`
}

func (meal Meal) Translated(language Language) (translated Meal) {
	translated = meal
	if language == German {
		translated.Name = meal.NameDe
	} else {
		translated.Name = meal.NameEn
	}
	translated.NameDe = ""
	translated.NameEn = ""

	translated.Category = meal.Category.Translated(language)

	translated.Additives = []MealAttribute{}
	for _, attribute := range meal.Additives {
		translated.Additives = append(translated.Additives, attribute.Translated(language))
	}
	translated.Allergens = []MealAttribute{}
	for _, attribute := range meal.Allergens {
		translated.Allergens = append(translated.Allergens, attribute.Translated(language))
	}
	translated.Features = []MealAttribute{}
	for _, attribute := range meal.Features {
		translated.Features = append(translated.Features, attribute.Translated(language))
	}
	return translated
}

type Nutrition struct {
	Kj            int     `json:"kj"`
	Kcal          int     `json:"kcal"`
	Fat           float64 `json:"fat"`
	SaturatedFat  float64 `json:"saturatedFat"`
	Carbohydrates float64 `json:"carbohydrates"`
	Sugar         float64 `json:"sugar"`
	Protein       float64 `json:"protein"`
	Salt          float64 `json:"salt"`
}

type MealAttribute struct {
	Id       int               `json:"id"`
	Type     MealAttributeType `json:"type,omitempty"`
	Short    string            `json:"short"`
	Name     string            `json:"name"`
	NameDe   string            `json:"nameDe,omitempty"`
	NameEn   string            `json:"nameEn,omitempty"`
	Location *Location         `json:"location,omitempty"`
}

func (attribute MealAttribute) Translated(language Language) (translated MealAttribute) {
	translated = attribute
	if language == German {
		translated.Name = attribute.NameDe
	} else {
		translated.Name = attribute.NameEn
	}
	translated.NameDe = ""
	translated.NameEn = ""
	translated.Type = ""
	translated.Location = nil
	return translated
}

type MealAttributeType string

const (
	AdditiveAttribute MealAttributeType = "additive"
	AllergenAttribute MealAttributeType = "allergen"
	FeatureAttribute  MealAttributeType = "feature"
	AllAttributes     MealAttributeType = ""
)

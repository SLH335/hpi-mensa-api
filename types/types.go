package types

import "time"

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

type Location struct {
	Id   int
	Name string
}

type Language int

const (
	German  Language = 1
	English Language = 2
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

var UnknownOffer MealCategory = MealCategory{
	Id:     0,
	NameDe: "Unbekanntes Angebot",
	NameEn: "Unknown Offer",
}
var Offer1 MealCategory = MealCategory{
	Id:     149,
	NameDe: "Angebot 1",
	NameEn: "Offer 1",
}
var Offer2 MealCategory = MealCategory{
	Id:     150,
	NameDe: "Angebot 2",
	NameEn: "Offer 2",
}
var Offer3 MealCategory = MealCategory{
	Id:     151,
	NameDe: "Angebot 3",
	NameEn: "Offer 3",
}
var Offer4 MealCategory = MealCategory{
	Id:     152,
	NameDe: "Angebot 4",
	NameEn: "Offer 4",
}
var DailySpecial MealCategory = MealCategory{
	Id:     118,
	NameDe: "Tagesangebot",
	NameEn: "Daily Special",
}
var SaladBar MealCategory = MealCategory{
	Id:     112,
	NameDe: "Salattheke",
	NameEn: "Salad bar",
}
var Dessert1 MealCategory = MealCategory{
	Id:     119,
	NameDe: "Dessert 1",
	NameEn: "Dessert 1",
}
var Dessert2 MealCategory = MealCategory{
	Id:     120,
	NameDe: "Dessert 1",
	NameEn: "Dessert 2",
}

func GetMealCategoryFromId(id int) MealCategory {
	switch id {
	case Offer1.Id:
		return Offer1
	case Offer2.Id:
		return Offer2
	case Offer3.Id:
		return Offer3
	case Offer4.Id:
		return Offer4
	case DailySpecial.Id:
		return DailySpecial
	case SaladBar.Id:
		return SaladBar
	case Dessert1.Id:
		return Dessert1
	case Dessert2.Id:
		return Dessert2
	default:
		return UnknownOffer
	}
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
	Location     Location        `json:"location"`
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
	return translated
}

type MealAttributeType string

const (
	AdditiveAttribute MealAttributeType = "additive"
	AllergenAttribute MealAttributeType = "allergen"
	FeatureAttribute  MealAttributeType = "feature"
	AllAttributes     MealAttributeType = ""
)

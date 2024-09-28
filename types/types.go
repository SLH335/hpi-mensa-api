package types

import "time"

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
	Id     int    `json:"id"`
	NameDe string `json:"nameDe"`
	NameEn string `json:"nameEn"`
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
	NameDe       string          `json:"nameDe"`
	NameEn       string          `json:"nameEn"`
	Category     MealCategory    `json:"category"`
	StudentPrice float64         `json:"studentPrice"`
	GuestPrice   float64         `json:"guestPrice"`
	Nutrition    Nutrition       `json:"nutrition"`
	Additives    []MealAttribute `json:"additives"`
	Allergens    []MealAttribute `json:"allergens"`
	Features     []MealAttribute `json:"features"`
	Location     Location        `json:"location"`
	Date         time.Time       `json:"date"`
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
	Id     int               `json:"id"`
	Type   MealAttributeType `json:"type"`
	Short  string            `json:"short"`
	NameDe string            `json:"nameDe"`
	NameEn string            `json:"nameEn"`
}

type MealAttributeType string

const (
	AdditiveAttribute MealAttributeType = "additive"
	AllergenAttribute MealAttributeType = "allergen"
	FeatureAttribute  MealAttributeType = "feature"
	AllAttributes     MealAttributeType = ""
)

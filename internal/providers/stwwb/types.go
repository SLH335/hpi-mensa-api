package stwwb

import (
	"hpi-mensa/internal/domain/meal"
	"time"
)

type Language int
const (
	LangGerman  Language = 1
	LangEnglish Language = 2
)

type Model string
const (
	AdditivesModel    Model = "additives"
	AllergensModel    Model = "allergens"
	FeaturesModel     Model = "features"
	LocationsModel    Model = "location"
	MealCategoryModel Model = "mealCategory"
	MenuModel         Model = "menu"
	OutletsModel      Model = "outlet"
)

type Meal struct {
	ID             int
	Name           meal.LangString
	Category       MealCategory
	Date           time.Time
	IsActive       bool
	PriceStudent   float64
	PriceEmployee  float64
	PriceGuest     float64
	Contingent     int
	Nutrition      MealNutrition
	Allergens      []MealAttribute
	Additives      []MealAttribute
	Features       []MealAttribute
	Sustainability MealSustainability
	Location       Location
	Outlet         Outlet
}

type Location struct {
	ID   int
	Name string
}

type Outlet struct {
	ID           int
	Name         string
	Location     Location
	OpeningHours OutletOpeningHours
	Address      OutletAddress
	ContactEmail string
	ImageLink    string
}

type OutletOpeningHours struct {
	Mo OutletOpeningTime
	Tu OutletOpeningTime
	We OutletOpeningTime
	Th OutletOpeningTime
	Fr OutletOpeningTime
	Sa OutletOpeningTime
	Su OutletOpeningTime
}

type OutletOpeningTime struct {
	Text    string
	Opening string
	Closing string
}

type OutletAddress struct {
	ID          int
	Street      string
	City        string
	PostalCode  string
	Country     string
	CountryCode string
	Coords      Coords
}

type Coords struct {
	Lat float64
	Lon float64
}

type MealCategory struct {
	ID       int
	Name     meal.LangString
	Location Location
}

type MealAttribute struct {
	ID    int
	Name  meal.LangString
	Short string
}

type MealAttributeType string
const (
	AllergenAttribute MealAttributeType = "allergens"
	AdditiveAttribute MealAttributeType = "additives"
	FeatureAttribute  MealAttributeType = "features"
)

type MealNutrition struct {
	Kj           int
	Kcal         int
	Fat          float64
	SaturatedFat float64
	Carbs        float64
	Sugar        float64
	Protein      float64
	Salt         float64
}

type MealSustainability struct {
	CO2 MealCO2
}

type MealCO2 struct {
	Grams  int
	Rating string
}

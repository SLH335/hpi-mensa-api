package common

import (
	"hpi-mensa/internal/mealdata/util"
	"time"
)

type Provider interface {
	Slug() string
	Name() util.LangString
	Init() error
	GetLocations() ([]Location, error)
	GetMeals(Location) ([]Meal, error)
}

type Meal struct {
	ID        string
	Name      util.LangString
	Category  util.LangString
	Date      time.Time
	Prices    []PriceCategory
	Additives []MealAttribute
	Allergens []MealAttribute
	Features  []MealAttribute
	Location  Location
	Provider  Provider
}

type Location struct {
	Slug     string
	Name     util.LangString
	Provider Provider
}

type PriceCategory struct {
	Name  util.LangString
	Price float64
}

type MealAttribute struct {
	Slug  string
	Name  util.LangString
	Short util.LangString
}

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
	ID        string          `json:"id"`
	Name      util.LangString `json:"name"`
	Category  util.LangString `json:"category"`
	Date      time.Time       `json:"date"`
	Prices    []PriceCategory `json:"price"`
	Additives []MealAttribute `json:"additives"`
	Allergens []MealAttribute `json:"allergens"`
	Features  []MealAttribute `json:"features"`
	Location  Location        `json:"location"`
	Provider  Provider        `json:"-"`
}

type Location struct {
	Slug     string          `json:"slug"`
	Name     util.LangString `json:"name"`
	Provider Provider        `json:"-"`
}

type PriceCategory struct {
	Type  util.LangString `json:"type"`
	Amount float64        `json:"amount"`
}

type MealAttribute struct {
	Slug  string          `json:"slug"`
	Name  util.LangString `json:"name"`
	Short util.LangString `json:"short"`
}

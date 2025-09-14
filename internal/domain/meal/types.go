package meal

import (
	"time"
)

type Provider interface {
	Slug() string
	Name() LangString
	Init() error
	GetLocations() ([]Location, error)
	GetMenus(Location) ([]Menu, error)
}

type Menu struct {
	Slug     string    `json:"slug"`
	Date     time.Time `json:"date"`
	Meals    []Meal    `json:"meals,nilasempty"`
	Location Location  `json:"location"`
	Provider Provider  `json:"-"`
}

type Meal struct {
	ID        string          `json:"id"`
	Name      LangString `json:"name"`
	Category  LangString `json:"category"`
	Date      time.Time       `json:"date"`
	Prices    []PriceCategory `json:"price"`
	Allergens []MealAttribute `json:"allergens"`
	Additives []MealAttribute `json:"additives"`
	Features  []MealAttribute `json:"features"`
	Location  Location        `json:"location"`
	Provider  Provider        `json:"-"`
}

type Location struct {
	Slug     string          `json:"slug"`
	Name     LangString `json:"name"`
	Provider Provider        `json:"-"`
}

type PriceCategory struct {
	Type  LangString `json:"type"`
	Amount float64        `json:"amount"`
}

type MealAttribute struct {
	Slug  string            `json:"slug"`
	Type  MealAttributeType `json:"type"`
	Name  LangString        `json:"name"`
	Short LangString        `json:"short"`
}

type MealAttributeType string
const (
	AllergenType MealAttributeType = "allergen"
	AdditiveType MealAttributeType = "additive"
	FeatureType MealAttributeType = "feature"
)

func (t MealAttributeType) IsValid() bool {
	return t == AllergenType || t == AdditiveType || t == FeatureType
}

type LangString struct {
	De string `json:"de"`
	En string `json:"en"`
}

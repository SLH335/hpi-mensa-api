package services

import (
	"time"

	"github.com/slh335/hpi-mensa-api/database"
	"github.com/slh335/hpi-mensa-api/services/mensadata"
	. "github.com/slh335/hpi-mensa-api/types"
)

type MealService struct {
	DbService        *database.MealDBService
	AttributeService *MealAttributeService
	CategoryService  *MealCategoryService
}

func (s *MealService) Get(location Location, lang Language, date time.Time) (meals []Meal, err error) {
	meals, err = s.DbService.Get(location, date)
	if err != nil {
		return []Meal{}, err
	}
	if len(meals) == 0 {
		meals, err = mensadata.GetMeals(location, lang, date)
		if err != nil {
			return []Meal{}, err
		}
		if len(meals) == 0 {
			return []Meal{}, nil
		}

		// ensure category and attribute data is present in database
		s.CategoryService.Get(location, lang)
		s.AttributeService.Get(AdditiveAttribute, location, lang)
		s.AttributeService.Get(AllergenAttribute, location, lang)
		s.AttributeService.Get(FeatureAttribute, location, lang)

		err = s.DbService.Add(meals)
		if err != nil {
			return []Meal{}, err
		}

		meals, err = s.DbService.Get(location, date)
		if err != nil {
			return []Meal{}, err
		}
	}
	for i := range meals {
		meals[i] = meals[i].Translated(lang)
	}
	return meals, nil
}

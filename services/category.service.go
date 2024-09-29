package services

import (
	"github.com/slh335/hpi-mensa-api/database"
	"github.com/slh335/hpi-mensa-api/services/mensadata"
	. "github.com/slh335/hpi-mensa-api/types"
)

type MealCategoryService struct {
	DbService *database.MealCategoryDBService
}

func (s *MealCategoryService) Get(location Location, lang Language) (categories []MealCategory, err error) {
	categories, err = s.DbService.Get(location)
	if err != nil {
		return []MealCategory{}, err
	}
	if len(categories) == 0 {
		categoriesDe, err := mensadata.GetMealCategories(location, German)
		if err != nil {
			return []MealCategory{}, err
		}
		categoriesEn, err := mensadata.GetMealCategories(location, English)
		if err != nil {
			return []MealCategory{}, err
		}
		for _, categoryDe := range categoriesDe {
			for _, categoryEn := range categoriesEn {
				if categoryDe.Id == categoryEn.Id {
					category := categoryEn
					category.NameDe = categoryDe.Name
					category.NameEn = categoryEn.Name
					category.Name = ""
					categories = append(categories, category)
					break
				}
			}
		}
		err = s.DbService.Add(categories)
		if err != nil {
			return []MealCategory{}, err
		}
	}
	return categories, nil
}

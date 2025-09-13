package stwwb

import (
	"errors"
	"fmt"
	"hpi-mensa/internal/mealdata/common/types"
	"hpi-mensa/internal/mealdata/util"
	"slices"

	"github.com/rs/zerolog/log"
)

// Store static data in memory
var locations []Location = []Location{}
var outlets map[int]Outlet = map[int]Outlet{} // outlets by outlet id
var categories map[int]map[int]MealCategory = map[int]map[int]MealCategory{} // meal categories per location id
var allergens map[int]map[int]MealAttribute = map[int]map[int]MealAttribute{} // allergens per location id
var additives map[int]map[int]MealAttribute = map[int]map[int]MealAttribute{} // additives per location id
var features map[int]map[int]MealAttribute = map[int]map[int]MealAttribute{} // features per location id

type STWWBProvider struct {}

var Provider STWWBProvider = STWWBProvider{}

func (p STWWBProvider) Slug() (slug string) {
	return "stwwb"
}

func (p STWWBProvider) Name() (name util.LangString) {
	return util.LangString{
		De: "Studierendenwerk West:Brandenburg",
		En: "Studierendenwerk West:Brandenburg",
	}
}

func (p STWWBProvider) Init() (err error) {
	locations, err = getLocations()
	if err != nil {
		return fmt.Errorf("get locations: %w", err)
	}

	outlets, err = getOutlets()
	if err != nil {
		return fmt.Errorf("get outlets: %w", err)
	}

	return nil
}

func (p STWWBProvider) GetLocations() (locs []common.Location, err error) {
	slugs := []string{}
	for _, location := range locations {
		slug := location.Slug()
		// Make sure location slugs are unique
		if slices.Contains(slugs, slug) {
			return []common.Location{}, errors.New("aborting due to duplicate location slug")
		}
		slugs = append(slugs, slug)

		locs = append(locs, common.Location{
			Slug: slug,
			Name: util.LangString{
				De: location.Name,
				En: location.Name,
			},
			Provider: Provider,
		})
	}
	return locs, nil
}

func (p STWWBProvider) GetMeals(location common.Location) (meals []common.Meal, err error) {
	// Find provider location from common location
	providerLocation := Location{ID: -1}
	for _, loc := range locations {
		if location.Slug == loc.Slug() {
			providerLocation = loc
		}
	}
	if providerLocation.ID == -1 {
		return []common.Meal{}, errors.New("location not found for provider")
	}

	// Fetch meal categories and attributes if not yet loaded
	if len(categories[providerLocation.ID]) == 0 {
		categories[providerLocation.ID], err = getMealCategories(providerLocation)
		if err != nil {
			return []common.Meal{}, fmt.Errorf("get meal categories: %w", err)
		}
		log.Debug().
			Str("provider", Provider.Slug()).
			Int("count", len(categories[providerLocation.ID])).
			Msg("Loaded meal categories")
	}
	if len(allergens[providerLocation.ID]) == 0 {
		allergens[providerLocation.ID], err = getMealAttributes(providerLocation, AllergenAttribute)
		if err != nil {
			return []common.Meal{}, fmt.Errorf("get allergens: %w", err)
		}
		log.Debug().
			Str("provider", Provider.Slug()).
			Int("count", len(allergens)).
			Msg("Loaded allergens")
	}
	if len(additives[providerLocation.ID]) == 0 {
		additives[providerLocation.ID], err = getMealAttributes(providerLocation, AdditiveAttribute)
		if err != nil {
			return []common.Meal{}, fmt.Errorf("get additives: %w", err)
		}
		log.Debug().
			Str("provider", Provider.Slug()).
			Int("count", len(additives[providerLocation.ID])).
			Msg("Loaded additives")
	}
	if len(features[providerLocation.ID]) == 0 {
		features[providerLocation.ID], err = getMealAttributes(providerLocation, FeatureAttribute)
		if err != nil {
			return []common.Meal{}, fmt.Errorf("get features: %w", err)
		}
		log.Debug().
			Str("provider", Provider.Slug()).
			Int("count", len(features[providerLocation.ID])).
			Msg("Loaded meal features")
	}

	menu, err := getMenu(providerLocation)
	if err != nil {
		return []common.Meal{}, fmt.Errorf("get menu: %w", err)
	}
	log.Debug().
		Str("provider", Provider.Slug()).
		Int("count", len(menu)).
		Msg("Loaded provider meals")

	for _, meal := range menu {
		meals = append(meals, common.Meal{
			ID: fmt.Sprintf("%s-%d", Provider.Slug(), meal.ID),
			Name: meal.Name,
			Category: meal.Category.Name,
			Date: meal.Date,
			Prices: convertPriceCategories(meal),
			Allergens: convertMealAttributes(meal.Allergens),
			Additives: convertMealAttributes(meal.Additives),
			Features: convertMealAttributes(meal.Features),
			Location: location,
			Provider: Provider,
		})
	}

	return meals, nil
}

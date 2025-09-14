package stwwb

import (
	"errors"
	"fmt"
	"hpi-mensa/internal/mealdata/common/types"
	"hpi-mensa/internal/mealdata/util"
	"slices"
	"time"

	"github.com/rs/zerolog/log"
)

type STWWBProvider struct {
	// Store static data in memory
	locations  []Location
	outlets    map[int]Outlet
	categories map[int]map[int]MealCategory // meal categories per location id
	allergens  map[int]map[int]MealAttribute // allergens per location id
	additives  map[int]map[int]MealAttribute // additives per location id
	features   map[int]map[int]MealAttribute // features per location id
}

var Provider = &STWWBProvider{
	locations:  []Location{},
	outlets:    map[int]Outlet{},
	categories: map[int]map[int]MealCategory{},
	allergens:  map[int]map[int]MealAttribute{},
	additives:  map[int]map[int]MealAttribute{},
	features:   map[int]map[int]MealAttribute{},
}

func (p *STWWBProvider) Slug() (slug string) {
	return "stwwb"
}

func (p *STWWBProvider) Name() (name util.LangString) {
	return util.LangString{
		De: "Studierendenwerk West:Brandenburg",
		En: "Studierendenwerk West:Brandenburg",
	}
}

func (p *STWWBProvider) Init() (err error) {
	p.locations, err = getLocations()
	if err != nil {
		return fmt.Errorf("get locations: %w", err)
	}

	p.outlets, err = getOutlets()
	if err != nil {
		return fmt.Errorf("get outlets: %w", err)
	}

	return nil
}

func (p *STWWBProvider) GetLocations() (locations []common.Location, err error) {
	slugs := []string{}
	for _, location := range p.locations {
		slug := location.Slug()
		// Make sure location slugs are unique
		if slices.Contains(slugs, slug) {
			return []common.Location{}, errors.New("aborting due to duplicate location slug")
		}
		slugs = append(slugs, slug)

		locations = append(locations, common.Location{
			Slug: slug,
			Name: util.LangString{
				De: location.Name,
				En: location.Name,
			},
			Provider: Provider,
		})
	}
	return locations, nil
}

// Get all menus that are currently available via the STWWB API
func (p *STWWBProvider) GetMenus(location common.Location) (menus []common.Menu, err error) {
	// Find provider location from common location
	providerLocation, err := p.getProviderLocation(location)
	if err != nil {
		return []common.Menu{}, fmt.Errorf("provider location: %w", err)
	}

	// Fetch meal categories and attributes if not yet loaded
	err = p.fetchMealMetadata(location)
	if err != nil {
		return []common.Menu{}, fmt.Errorf("meal metadata: %w", err)
	}

	meals, start, end, err := getMeals(providerLocation)
	if err != nil {
		return []common.Menu{}, fmt.Errorf("get menu: %w", err)
	}
	log.Debug().
		Str("provider", Provider.Slug()).
		Int("count", len(meals)).
		Msg("Loaded provider meals")

	// Generate a menu for every day that the plan is valid for
	date := start
	for date.Before(end) {
		menu := common.Menu{
			Slug: p.Slug() + "-" + date.Format("2006-01-02"),
			Date: date,
			Provider: p,
		}

		// Add all meals from current date to menu
		for _, meal := range meals {
			if !util.SameDay(date, meal.Date) {
				continue
			}
			menu.Meals = append(menu.Meals, common.Meal{
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

		menus = append(menus, menu)

		date = date.Add(24 * time.Hour)
	}

	return menus, nil
}

func (p *STWWBProvider) getProviderLocation(location common.Location) (providerLocation Location, err error) {
	providerLocation = Location{ID: -1}
	for _, loc := range p.locations {
		if location.Slug == loc.Slug() {
			providerLocation = loc
		}
	}
	if providerLocation.ID == -1 {
		return Location{}, errors.New("location not found for provider")
	}
	return providerLocation, nil
}

func (p *STWWBProvider) fetchMealMetadata(location common.Location) (err error) {
	providerLocation, err := p.getProviderLocation(location)
	if err != nil {
		return fmt.Errorf("provider location: %w", err)
	}

	if len(p.categories[providerLocation.ID]) == 0 {
		p.categories[providerLocation.ID], err = getMealCategories(providerLocation)
		if err != nil {
			return fmt.Errorf("get meal categories: %w", err)
		}
		log.Debug().
			Str("provider", Provider.Slug()).
			Str("location", location.Slug).
			Int("count", len(p.categories[providerLocation.ID])).
			Msg("Loaded meal categories")
	}
	if len(p.allergens[providerLocation.ID]) == 0 {
		p.allergens[providerLocation.ID], err = getMealAttributes(providerLocation, AllergenAttribute)
		if err != nil {
			return fmt.Errorf("get allergens: %w", err)
		}
		log.Debug().
			Str("provider", Provider.Slug()).
			Str("location", location.Slug).
			Int("count", len(p.allergens[providerLocation.ID])).
			Msg("Loaded allergens")
	}
	if len(p.additives[providerLocation.ID]) == 0 {
		p.additives[providerLocation.ID], err = getMealAttributes(providerLocation, AdditiveAttribute)
		if err != nil {
			return fmt.Errorf("get additives: %w", err)
		}
		log.Debug().
			Str("provider", Provider.Slug()).
			Str("location", location.Slug).
			Int("count", len(p.additives[providerLocation.ID])).
			Msg("Loaded additives")
	}
	if len(p.features[providerLocation.ID]) == 0 {
		p.features[providerLocation.ID], err = getMealAttributes(providerLocation, FeatureAttribute)
		if err != nil {
			return fmt.Errorf("get features: %w", err)
		}
		log.Debug().
			Str("provider", Provider.Slug()).
			Str("location", location.Slug).
			Int("count", len(p.features[providerLocation.ID])).
			Msg("Loaded meal features")
	}

	return nil
}

package stwwb

import (
	"errors"
	"fmt"
	"hpi-mensa/internal/domain/meal"
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

func (p *STWWBProvider) Name() (name meal.LangString) {
	return meal.LangString{
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

func (p *STWWBProvider) GetLocations() (locations []meal.Location, err error) {
	slugs := []string{}
	for _, location := range p.locations {
		slug := location.Slug()
		// Make sure location slugs are unique
		if slices.Contains(slugs, slug) {
			return []meal.Location{}, errors.New("aborting due to duplicate location slug")
		}
		slugs = append(slugs, slug)

		locations = append(locations, meal.Location{
			Slug: slug,
			Name: meal.LangString{
				De: location.Name,
				En: location.Name,
			},
			Provider: Provider,
		})
	}
	return locations, nil
}

// Get all menus that are currently available via the STWWB API
func (p *STWWBProvider) GetMenus(location meal.Location) (menus []meal.Menu, err error) {
	// Find provider location from common location
	providerLocation, err := p.getProviderLocation(location)
	if err != nil {
		return []meal.Menu{}, fmt.Errorf("provider location: %w", err)
	}

	// Fetch meal categories and attributes if not yet loaded
	err = p.fetchMealMetadata(location)
	if err != nil {
		return []meal.Menu{}, fmt.Errorf("meal metadata: %w", err)
	}

	meals, start, end, err := getMeals(providerLocation)
	if err != nil {
		return []meal.Menu{}, fmt.Errorf("get menu: %w", err)
	}
	log.Debug().
		Str("provider", Provider.Slug()).
		Int("count", len(meals)).
		Msg("Loaded provider meals")

	// Generate a menu for every day that the plan is valid for
	date := start
	for date.Before(end) {
		menu := meal.Menu{
			Slug: p.Slug() + "-" + date.Format("2006-01-02"),
			Date: date,
			Provider: p,
		}

		// Add all meals from current date to menu
		for _, m := range meals {
			if !meal.SameDay(date, m.Date) {
				continue
			}
			menu.Meals = append(menu.Meals, meal.Meal{
				ID: fmt.Sprintf("%s-%d", Provider.Slug(), m.ID),
				Name: m.Name,
				Category: m.Category.Name,
				Date: m.Date,
				Prices: convertPriceCategories(m),
				Allergens: convertMealAttributes(m.Allergens),
				Additives: convertMealAttributes(m.Additives),
				Features: convertMealAttributes(m.Features),
				Location: location,
				Provider: Provider,
			})
		}

		menus = append(menus, menu)

		date = date.Add(24 * time.Hour)
	}

	return menus, nil
}

func (p *STWWBProvider) getProviderLocation(location meal.Location) (providerLocation Location, err error) {
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

func (p *STWWBProvider) fetchMealMetadata(location meal.Location) (err error) {
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

package mealdata

import (
	"fmt"
	"hpi-mensa/internal/mealdata/common/types"
	"hpi-mensa/internal/mealdata/providers/stwwb"

	"github.com/rs/zerolog/log"
)

// Active meal data providers
var providers []common.Provider = []common.Provider{
	stwwb.Provider, // Studierendenwerk West:Brandenburg
}

func Init() (err error) {
	for _, provider := range providers {
		err = provider.Init()
		if err != nil {
			return fmt.Errorf("%s init: %w", provider.Slug(), err)
		}
		log.Debug().Str("provider", provider.Slug()).Msg("Initialized provider")
	}

	return nil
}

func GetLocations() (locations []common.Location, err error) {
	for _, provider := range providers {
		log.Debug().Str("provider", provider.Slug()).Msg("Loading provider locations")
		providerLocations, err := provider.GetLocations()
		if err != nil {
			return []common.Location{}, fmt.Errorf("%s locations: %w", provider.Slug(), err)
		}
		locations = append(locations, providerLocations...)
		log.Debug().
			Str("provider", provider.Slug()).
			Int("count", len(locations)).
			Msg("Loaded provider locations successfully")
	}

	return locations, nil
}

func GetMeals(location common.Location) (meals []common.Meal, err error) {
	meals, err = location.Provider.GetMeals(location)
	if err != nil {
		return []common.Meal{}, fmt.Errorf("%s meals: %w", location.Provider.Slug(), err)
	}

	return meals, nil
}

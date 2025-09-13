package mealdata

import (
	"fmt"
	"hpi-mensa/internal/mealdata/common/types"
	"hpi-mensa/internal/mealdata/providers/stwwb"
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
	}

	return nil
}

func GetLocations() (locations []common.Location, err error) {
	for _, provider := range providers {
		providerLocations, err := provider.GetLocations()
		if err != nil {
			return []common.Location{}, fmt.Errorf("%s locations: %w", provider.Slug(), err)
		}
		locations = append(locations, providerLocations...)
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

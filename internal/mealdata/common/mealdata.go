package mealdata

import (
	"errors"
	"fmt"
	"hpi-mensa/internal/mealdata/common/types"
	"hpi-mensa/internal/mealdata/providers/stwwb"
	"hpi-mensa/internal/mealdata/util"
	"time"

	"github.com/rs/zerolog/log"
)

// Active meal data providers
var providers []common.Provider = []common.Provider{
	stwwb.Provider, // Studierendenwerk West:Brandenburg
}

// Number of active providers
func ProviderCount() (count int) {
	return len(providers)
}

// Initialize meal data providers
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

// Get all locations from all providers
func GetLocations() (locations []common.Location, err error) {
	for _, provider := range providers {
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

// Fetch all currently accessible menus for a given location
func FetchMenus(location common.Location) (menus []common.Menu, err error) {
	menus, err = location.Provider.GetMenus(location)
	if err != nil {
		return []common.Menu{}, fmt.Errorf("%s menus: %w", location.Provider.Slug(), err)
	}

	return menus, nil
}

// Get meals for a given day (TODO: from the DB and fetches new meals if no menu is stored yet)
func GetMenu(location common.Location, date time.Time) (menu common.Menu, err error) {
	menus, err := FetchMenus(location)
	if err != nil {
		return common.Menu{}, fmt.Errorf("fetch menu: %w", err)
	}

	// Return menu from specified date
	for _, menu := range menus {
		log.Debug().
			Int("meals", len(menu.Meals)).
			Str("date", menu.Date.Format("2006-01-02")).
			Msg("Loaded menu")
		if util.SameDay(date, menu.Date) {
			return menu, nil
		}
	}
	return common.Menu{}, errors.New(fmt.Sprintf("no menu available for date %s", date.Format("2006-01-02")))
}

package services

import (
	"github.com/slh335/hpi-mensa-api/database"
	"github.com/slh335/hpi-mensa-api/services/mensadata"
	. "github.com/slh335/hpi-mensa-api/types"
)

type LocationService struct {
	DbService *database.LocationDBService
}

func (s *LocationService) Get() (locations []Location, err error) {
	locations, err = s.DbService.Get()
	if err != nil {
		return []Location{}, err
	}
	if len(locations) == 0 {
		locations, err = mensadata.GetLocations()
		if err != nil {
			return []Location{}, err
		}
		locations = append(locations, Location{
			Id:   999,
			Name: "Campus Kitchen One",
		})
		err = s.DbService.Add(locations)
		if err != nil {
			return []Location{}, err
		}
	}

	return locations, nil
}

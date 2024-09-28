package database

import (
	"database/sql"
	"fmt"

	. "github.com/slh335/hpi-mensa-api/types"
)

type LocationService struct {
	DB *sql.DB
}

func (s *LocationService) GetAll() (locations []Location, err error) {
	stmt := "SELECT * FROM locations"
	rows, err := s.DB.Query(stmt)
	if err != nil {
		return []Location{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var location Location
		err = rows.Scan(&location.Id, &location.Name)
		if err != nil {
			return []Location{}, err
		}
		locations = append(locations, location)
	}

	return locations, nil
}

func (s *LocationService) Add(locations []Location) (err error) {
	if len(locations) == 0 {
		return fmt.Errorf("error: no locations were provided")
	}

	stmt := "INSERT INTO locations (id, name) VALUES "
	args := []any{}
	for i, location := range locations {
		if i != 0 {
			stmt += ", "
		}
		stmt += "(?, ?)"
		args = append(args, location.Id, location.Name)
	}

	_, err = s.DB.Exec(stmt, args...)
	return err
}

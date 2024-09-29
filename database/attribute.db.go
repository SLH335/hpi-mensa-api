package database

import (
	"database/sql"

	. "github.com/slh335/hpi-mensa-api/types"
)

type MealAttributeDBService struct {
	DB *sql.DB
}

func (s *MealAttributeDBService) GetAll(location Location) (attributes []MealAttribute, err error) {
	return s.Get(AllAttributes, location)
}

func (s *MealAttributeDBService) Get(attributeType MealAttributeType, location Location) (attributes []MealAttribute, err error) {
	stmt := "SELECT * FROM attributes WHERE location_id=?"
	var rows *sql.Rows
	if attributeType != AllAttributes {
		stmt += " AND type=?"
		rows, err = s.DB.Query(stmt, location.Id, attributeType)
	} else {
		rows, err = s.DB.Query(stmt, location.Id)
	}
	if err != nil {
		return []MealAttribute{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var attribute MealAttribute
		err = rows.Scan(&attribute.Id, &attribute.Type, &attribute.Short, &attribute.NameDe, &attribute.NameEn, &attribute.Location.Id)
		if err != nil {
			return []MealAttribute{}, err
		}
		attributes = append(attributes, attribute)
	}

	return attributes, nil
}

func (s *MealAttributeDBService) Add(attributes []MealAttribute) (err error) {
	stmt := "INSERT INTO attributes (id, type, short, name_de, name_en, location_id) VALUES "
	args := []any{}
	for i, attribute := range attributes {
		if i != 0 {
			stmt += ", "
		}
		stmt += "(?, ?, ?, ?, ?, ?)"
		args = append(args, attribute.Id, attribute.Type, attribute.Short, attribute.NameDe, attribute.NameEn, attribute.Location.Id)
	}

	_, err = s.DB.Exec(stmt, args...)
	return err
}

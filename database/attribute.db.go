package database

import (
	"database/sql"

	. "github.com/slh335/hpi-mensa-api/types"
)

type MealAttributeService struct {
	DB *sql.DB
}

func (s *MealAttributeService) GetAll() (features []MealAttribute, err error) {
	return s.Get(AllAttributes)
}

func (s *MealAttributeService) Get(featureType MealAttributeType) (features []MealAttribute, err error) {
	stmt := "SELECT * FROM features"
	var rows *sql.Rows
	if featureType != AllAttributes {
		stmt += " WHERE type=?"
		rows, err = s.DB.Query(stmt, featureType)
	} else {
		rows, err = s.DB.Query(stmt)
	}
	if err != nil {
		return []MealAttribute{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var feature MealAttribute
		err = rows.Scan(&feature.Id, &feature.Type, &feature.Short, &feature.NameDe, &feature.NameEn)
		if err != nil {
			return []MealAttribute{}, err
		}
		features = append(features, feature)
	}

	return features, nil
}

func (s *MealAttributeService) Add(features []MealAttribute) (err error) {
	stmt := "INSERT INTO features (id, type, short, name_de, name_en) VALUES "
	args := []any{}
	for i, feature := range features {
		if i != 0 {
			stmt += ", "
		}
		stmt += "(?, ?, ?, ?, ?)"
		args = append(args, feature.Id, feature.Type, feature.Short, feature.NameDe, feature.NameEn)
	}

	_, err = s.DB.Exec(stmt, args...)
	return err
}

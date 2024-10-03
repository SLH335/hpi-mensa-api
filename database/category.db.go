package database

import (
	"database/sql"

	"github.com/cockroachdb/errors"
	. "github.com/slh335/hpi-mensa-api/types"
)

type MealCategoryDBService struct {
	DB *sql.DB
}

func (s *MealCategoryDBService) Get(location Location) (categories []MealCategory, err error) {
	stmt := "SELECT * FROM categories WHERE categories.location_id=?"
	rows, err := s.DB.Query(stmt, location.Id)
	if err != nil {
		return []MealCategory{}, errors.Wrap(err, "meal category db get")
	}
	defer rows.Close()

	for rows.Next() {
		var category MealCategory
		category.Location = &Location{}
		err = rows.Scan(&category.Id, &category.NameDe, &category.NameEn, &category.Location.Id)
		if err != nil {
			return []MealCategory{}, errors.Wrap(err, "meal category db get")
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (s *MealCategoryDBService) Add(categories []MealCategory) (err error) {
	stmt := "INSERT INTO categories (id, name_de, name_en, location_id) VALUES "
	args := []any{}
	for i, category := range categories {
		if i != 0 {
			stmt += ", "
		}
		stmt += "(?, ?, ?, ?)"
		args = append(args, category.Id, category.NameDe, category.NameEn, category.Location.Id)
	}

	_, err = s.DB.Exec(stmt, args...)
	return errors.Wrap(err, "meal category db add")
}

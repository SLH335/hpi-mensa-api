package database

import (
	"database/sql"

	. "github.com/slh335/hpi-mensa-api/types"
)

type MealCategoryService struct {
	DB *sql.DB
}

func (s *MealCategoryService) Get() (categories []MealCategory, err error) {
	stmt := "SELECT * FROM categories"
	rows, err := s.DB.Query(stmt)
	if err != nil {
		return []MealCategory{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var category MealCategory
		err = rows.Scan(&category.Id, &category.NameDe, &category.NameEn)
		if err != nil {
			return []MealCategory{}, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (s *MealCategoryService) Add(categories []MealCategory) (err error) {
	stmt := "INSERT INTO categories (id, name_de, name_en) VALUES "
	args := []any{}
	for i, category := range categories {
		if i != 0 {
			stmt += ", "
		}
		stmt += "(?, ?, ?)"
		args = append(args, category.Id, category.NameDe, category.NameEn)
	}

	_, err = s.DB.Exec(stmt, args...)
	return err
}

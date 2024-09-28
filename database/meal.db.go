package database

import (
	"database/sql"
	"fmt"
	"time"

	. "github.com/slh335/hpi-mensa-api/types"
)

type MealService struct {
	DB *sql.DB
}

func (s *MealService) Get(location Location, day time.Time) (meals []Meal, err error) {
	// get basic meal data
	stmt := `SELECT meals.id, meals.name_de, meals.name_en, meals.price_student, meals.price_guest,
			categories.id, categories.name_de, categories.name_en, locations.id, locations.name,
			nutrition.kj, nutrition.kcal, nutrition.fat, nutrition.saturated_fat,
			nutrition.carbohydrates, nutrition.sugar, nutrition.salt FROM meals
		INNER JOIN categories ON categories.id=meals.category_id
		INNER JOIN locations ON locations.id=meals.location_id
		INNER JOIN nutrition ON nutrition.meal_id=meals.id
		WHERE meals.location_id=? AND meals.date=?`
	dateStr := day.Format("2006-01-02")
	rows, err := s.DB.Query(stmt, location.Id, dateStr)
	if err != nil {
		return []Meal{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var meal Meal
		err = rows.Scan(
			&meal.Id, &meal.NameDe, &meal.NameEn, &meal.StudentPrice, &meal.GuestPrice,
			&meal.Category.Id, &meal.Category.NameDe, &meal.Category.NameEn, &meal.Location.Id,
			&meal.Location.Name, &meal.Nutrition.Kj, &meal.Nutrition.Kcal, &meal.Nutrition.Fat,
			&meal.Nutrition.SaturatedFat, &meal.Nutrition.Carbohydrates, &meal.Nutrition.Sugar,
			&meal.Nutrition.Salt,
		)
		if err != nil {
			return []Meal{}, err
		}
		meal.Date, _ = time.Parse("2006-01-02", dateStr)
		meals = append(meals, meal)
	}

	// get meal additives, allergens and features
	stmt = `SELECT meal_attributes.meal_id, attributes.* FROM meal_attributes
		INNER JOIN attributes ON meal_attributes.attribute_id=attributes.id WHERE `
	args := []any{}
	for i, meal := range meals {
		if i != 0 {
			stmt += " OR "
		}
		stmt += "meal_attributes.meal_id=?"
		args = append(args, meal.Id)
	}

	rows, err = s.DB.Query(stmt, args...)
	if err != nil {
		return []Meal{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var attribute MealAttribute
		var mealId int
		err = rows.Scan(&mealId, &attribute.Id, &attribute.Type, &attribute.Short,
			&attribute.NameDe, &attribute.NameEn)
		if err != nil {
			return []Meal{}, err
		}
		for _, meal := range meals {
			if meal.Id == mealId {
				switch attribute.Type {
				case AdditiveAttribute:
					meal.Additives = append(meal.Additives, attribute)
				case AllergenAttribute:
					meal.Allergens = append(meal.Allergens, attribute)
				case FeatureAttribute:
					meal.Features = append(meal.Features, attribute)
				}
				break
			}
		}
	}

	return meals, nil
}

func (s *MealService) Add(meals []Meal) (err error) {
	if len(meals) == 0 {
		return fmt.Errorf("error: no meals were provided")
	}

	// add basic meal data
	stmt := `INSERT INTO meals (id, name_de, name_en, category_id, price_student, price_guest, date,
		location_id) VALUES `
	args := []any{}
	for i, meal := range meals {
		if i != 0 {
			stmt += ", "
		}
		stmt += "(?, ?, ?, ?, ?, ?, ?, ?)"
		args = append(args, meal.Id, meal.NameDe, meal.NameEn, meal.Category.Id, meal.StudentPrice,
			meal.GuestPrice, meal.Date.Format("2006-01-02"), meal.Location.Id)
	}
	_, err = s.DB.Exec(stmt, args...)
	if err != nil {
		return err
	}

	// add nutrition information
	stmt = `INSERT INTO nutrition (meal_id, kj, kcal, fat, saturated_fat, carbohydrates, sugar,
		protein, salt) VALUES `
	args = []any{}
	for i, meal := range meals {
		if i != 0 {
			stmt += ", "
		}
		stmt += "(?, ?, ?, ?, ?, ?, ?, ?, ?)"
		args = append(args, meal.Id, meal.Nutrition.Kj, meal.Nutrition.Kcal, meal.Nutrition.Fat,
			meal.Nutrition.SaturatedFat, meal.Nutrition.Carbohydrates, meal.Nutrition.Sugar,
			meal.Nutrition.Protein, meal.Nutrition.Salt)
	}
	_, err = s.DB.Exec(stmt, args...)

	// add additives, allergens and features
	stmt = `INSERT INTO meal_attributes (meal_id, attribute_id) VALUES `
	args = []any{}
	for i, meal := range meals {
		mealAttributes := append(meal.Additives, append(meal.Allergens, meal.Features...)...)
		for j, mealAttribute := range mealAttributes {
			if i+j != 0 {
				stmt += ", "
			}
			stmt += "(?, ?)"
			args = append(args, meal.Id, mealAttribute.Id)
		}
	}
	_, err = s.DB.Exec(stmt, args...)

	return err
}

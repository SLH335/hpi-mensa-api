package database

import (
	"database/sql"
	"fmt"
	"time"

	. "github.com/slh335/hpi-mensa-api/types"
)

type MealDBService struct {
	DB *sql.DB
}

func (s *MealDBService) Get(location Location, date time.Time) (meals []Meal, err error) {
	// get basic meal data
	stmt := `SELECT meals.id, meals.name_de, meals.name_en, meals.price_student, meals.price_guest,
			categories.id, categories.name_de, categories.name_en, locations.id, locations.name,
			nutrition.kj, nutrition.kcal, nutrition.fat, nutrition.saturated_fat,
			nutrition.carbohydrates, nutrition.sugar, nutrition.protein, nutrition.salt,
			meals.co2_grams, meals.co2_rating FROM meals
		LEFT JOIN categories ON categories.id=meals.category_id
		LEFT JOIN locations ON locations.id=meals.location_id
		LEFT JOIN nutrition ON nutrition.meal_id=meals.id
		WHERE meals.location_id=? AND meals.date=?`
	dateStr := date.Format("2006-01-02")
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
			&meal.Nutrition.Protein, &meal.Nutrition.Salt, &meal.CO2.Grams, &meal.CO2.Rating,
		)
		if err != nil {
			return []Meal{}, err
		}
		meal.Date, _ = time.Parse("2006-01-02", dateStr)
		meals = append(meals, meal)
	}

	if len(meals) == 0 {
		return meals, nil
	}

	// get meal additives, allergens and features
	stmt = `SELECT meal_attributes.meal_id, attributes.* FROM meal_attributes
		INNER JOIN attributes ON meal_attributes.attribute_id=attributes.id
			AND meal_attributes.attribute_type=attributes.type WHERE `
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
		attribute.Location = &Location{}
		var mealId int
		err = rows.Scan(&mealId, &attribute.Id, &attribute.Type, &attribute.Short,
			&attribute.NameDe, &attribute.NameEn, &attribute.Location.Id)
		if err != nil {
			return []Meal{}, err
		}
		for i := range meals {
			if meals[i].Id == mealId {
				switch attribute.Type {
				case AdditiveAttribute:
					meals[i].Additives = append(meals[i].Additives, attribute)
				case AllergenAttribute:
					meals[i].Allergens = append(meals[i].Allergens, attribute)
				case FeatureAttribute:
					meals[i].Features = append(meals[i].Features, attribute)
				}
				break
			}
		}
	}

	return meals, nil
}

func (s *MealDBService) Add(meals []Meal) (err error) {
	if len(meals) == 0 {
		return fmt.Errorf("error: no meals were provided")
	}

	// add basic meal data
	stmt := `INSERT INTO meals (id, name_de, name_en, category_id, price_student, price_guest,
		co2_grams, co2_rating, date, location_id) VALUES `
	args := []any{}
	for i, meal := range meals {
		if i != 0 {
			stmt += ", "
		}
		stmt += "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
		args = append(args, meal.Id, meal.NameDe, meal.NameEn, meal.Category.Id, meal.StudentPrice,
			meal.GuestPrice, meal.CO2.Grams, meal.CO2.Rating, meal.Date.Format("2006-01-02"),
			meal.Location.Id)
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
	if len(args) > 0 {
		_, err = s.DB.Exec(stmt, args...)
	}
	if err != nil {
		return err
	}

	// add additives, allergens and features
	stmt = `INSERT INTO meal_attributes (meal_id, attribute_id, attribute_type) VALUES `
	args = []any{}
	firstAdd := true
	for _, meal := range meals {
		mealAttributes := append(meal.Additives, append(meal.Allergens, meal.Features...)...)
		for _, mealAttribute := range mealAttributes {
			if firstAdd {
				firstAdd = false
			} else {
				stmt += ", "
			}
			stmt += "(?, ?, ?)"
			args = append(args, meal.Id, mealAttribute.Id, mealAttribute.Type)
		}
	}
	if len(args) > 0 {
		_, err = s.DB.Exec(stmt, args...)
	}
	if err != nil {
		return err
	}

	return err
}

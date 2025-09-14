package repositories

import (
	"context"
	"fmt"
	"hpi-mensa/gen/db"
	"hpi-mensa/internal/domain/meal"
)

type MealRepository interface {
	GetByID(ctx context.Context, mealID string) (*meal.Meal, error)
	ListByMenu(ctx context.Context, menuSlug string) ([]*meal.Meal, error)
	Create(ctx context.Context, meal *meal.Meal) error
	Update(ctx context.Context, meal *meal.Meal) error
	Delete(ctx context.Context, mealID string) error
	SetAttributesForMeal(ctx context.Context, mealID string, attributeSlugs []string) error
}

type mealRepo struct {
	q *db.Queries
}

func NewMealRepository(queries *db.Queries) MealRepository {
	return &mealRepo{
		q: queries,
	}
}

func (r *mealRepo) GetByID(ctx context.Context, mealID string) (*meal.Meal, error) {
	row, err := r.q.GetMealByID(ctx, mealID)
	if err != nil {
		return nil, fmt.Errorf("get meal %s: %w", mealID, err)
	}
	meal := db.Meal{
		ID: row.ID,
		NameDe: row.NameDe,
		NameEn: row.NameEn,
		CategoryDe: row.CategoryDe,
		CategoryEn: row.CategoryEn,
		Date: row.Date,
		MenuSlug: row.MenuSlug,
		LocationSlug: row.LocationSlug,
	}
	location := db.Location{
		Slug: row.LocationSlug,
		NameDe: row.LocationNameDe,
		NameEn: row.LocationNameEn,
	}

	attrs, err := r.q.ListAttributesByMeal(ctx, mealID)
	if err != nil {
		return nil, fmt.Errorf("list attributes for meal %s: %w", mealID, err)
	}
	attributes := []db.Attribute{}
	for _, a := range attrs {
		if !a.Slug.Valid || !a.Type.Valid || !a.NameDe.Valid || !a.NameEn.Valid || !a.ShortDe.Valid || !a.ShortEn.Valid {
			continue
		}
		attributes = append(attributes, db.Attribute{
			Slug: a.Slug.String,
			Type: a.Type.String,
			NameDe: a.NameDe.String,
			NameEn: a.NameEn.String,
			ShortDe: a.ShortDe.String,
			ShortEn: a.ShortEn.String,
		})
	}

	prices, err := r.q.ListPricesByMeal(ctx, mealID)
	if err != nil {
		return nil, fmt.Errorf("list prices for meal %s: %w", mealID, err)
	}
	
	return mapMealToDomain(meal, attributes, prices, location), nil
}


func (r *mealRepo) ListByMenu(ctx context.Context, menuSlug string) ([]*meal.Meal, error) {
	// Get meals for menu
	rows, err := r.q.ListMealsByMenu(ctx, menuSlug)
	if err != nil {
		return nil, fmt.Errorf("list meals for menu %s: %w", menuSlug, err)
	}
	mealIds := []string{}
	for _, row := range rows {
		mealIds = append(mealIds, row.ID)
	}

	// Get all attributes for returned meals
	attributes, err := r.q.ListAttributesForMeals(ctx, mealIds)
	if err != nil {
		return nil, fmt.Errorf("list attributes for menu %s: %w", menuSlug, err)
	}
	attributesByMeal := map[string][]db.Attribute{}
	for _, a := range attributes {
		attributesByMeal[a.MealID] = append(attributesByMeal[a.MealID], db.Attribute{
			Slug: a.Slug,
			Type: a.Type,
			NameDe: a.NameDe,
			NameEn: a.NameEn,
			ShortDe: a.ShortDe,
			ShortEn: a.ShortEn,
		})
	}

	// Map meal data
	meals := []*meal.Meal{}
	for _, row := range rows {
		meal := db.Meal{
			ID: row.ID,
			NameDe: row.NameDe,
			NameEn: row.NameEn,
			CategoryDe: row.CategoryDe,
			CategoryEn: row.CategoryEn,
			Date: row.Date,
			MenuSlug: row.MenuSlug,
			LocationSlug: row.LocationSlug,
		}
		location := db.Location{
			Slug: row.LocationSlug,
			NameDe: row.LocationNameDe,
			NameEn: row.LocationNameEn,
		}
		meals = append(meals, mapMealToDomain(meal, attributesByMeal[meal.ID], nil, location))
	}

	return meals, nil
}

func (r *mealRepo) Create(ctx context.Context, meal *meal.Meal) error {
	err := r.q.CreateMeal(ctx, db.CreateMealParams{
		ID: meal.ID,
		NameDe: meal.Name.De,
		NameEn: meal.Name.En,
		CategoryDe: meal.Category.De,
		CategoryEn: meal.Category.En,
		Date: meal.Date,
		MenuSlug: meal.M,
	})

	return nil
}

func (r *mealRepo) Update(ctx context.Context, meal *meal.Meal) error {}

func (r *mealRepo) Delete(ctx context.Context, mealID string) error {}

func mapMealToDomain(m db.Meal, attrs []db.Attribute, prices []db.Price, location db.Location) *meal.Meal {
	// Map attributes
	allergens := []meal.MealAttribute{}
	additives := []meal.MealAttribute{}
	features := []meal.MealAttribute{}
	for _, a := range attrs {
		aType := meal.MealAttributeType(a.Type)
		if !aType.IsValid() {
			continue
		}
		attribute := meal.MealAttribute{
			Slug: a.Slug,
			Type: aType,
			Name: meal.LangString{
				De: a.NameDe,
				En: a.NameEn,
			},
			Short: meal.LangString{
				De: a.ShortDe,
				En: a.ShortEn,
			},
		}
		switch aType {
		case meal.AllergenType:
			allergens = append(allergens, attribute)
		case meal.AdditiveType:
			additives = append(additives, attribute)
		case meal.FeatureType:
			features = append(features, attribute)
		}
	}

	// Map prices
	priceCategories := []meal.PriceCategory{}
	for _, p := range prices {
		priceCategories = append(priceCategories, meal.PriceCategory{
			Type: meal.LangString{
				De: p.TypeDe,
				En: p.TypeEn,
			},
			Amount: p.Amount,
		})
	}

	// Map meal
	return &meal.Meal{
		ID: m.ID,
		Name: meal.LangString{
			De: m.NameDe,
			En: m.NameEn,
		},
		Category: meal.LangString{
			De: m.CategoryDe,
			En: m.CategoryEn,
		},
		Date: m.Date,
		Prices: priceCategories,
		Allergens: allergens,
		Additives: additives,
		Features: features,
		Location: meal.Location{
			Slug: location.Slug,
			Name: meal.LangString{
				De: location.NameDe,
				En: location.NameEn,
			},
		},
	}
}

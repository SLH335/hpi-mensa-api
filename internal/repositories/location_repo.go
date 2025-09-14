package repositories

import (
	"context"
	"fmt"
	"hpi-mensa/gen/db"
	"hpi-mensa/internal/domain/meal"
)

type LocationRepository interface {
	GetBySlug(ctx context.Context, slug string) (*meal.Location, error)
	ListAll(ctx context.Context) ([]*meal.Location, error)
	ListByProvider(ctx context.Context, provider string) ([]*meal.Location, error)
	Create(ctx context.Context, loc *meal.Location) error
	Delete(ctx context.Context, slug string) error
}

type locationRepo struct {
	q *db.Queries
}

func NewLocationRepository(queries *db.Queries) LocationRepository {
	return &locationRepo{q: queries}
}

func (r *locationRepo) GetBySlug(ctx context.Context, slug string) (*meal.Location, error) {
	row, err := r.q.GetLocationBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("get location by slug: %w", err)
	}

	return mapLocToDomain(row), nil
}

func (r *locationRepo) ListAll(ctx context.Context) ([]*meal.Location, error) {
	rows, err := r.q.ListLocations(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all locations: %w", err)
	}

	locations := []*meal.Location{}
	for _, row := range rows {
		locations = append(locations, mapLocToDomain(row))
	}
	return locations, nil
}

func (r *locationRepo) ListByProvider(ctx context.Context, provider string) ([]*meal.Location, error) {
	rows, err := r.q.ListLocationsByProvider(ctx, provider)
	if err != nil {
		return nil, fmt.Errorf("get locations by provider: %w", err)
	}

	locations := []*meal.Location{}
	for _, row := range rows {
		locations = append(locations, mapLocToDomain(row))
	}
	return locations, nil
}

func (r *locationRepo) Create(ctx context.Context, loc *meal.Location) error {
	err := r.q.CreateLocation(ctx, db.CreateLocationParams{
		Slug: loc.Slug,
		NameDe: loc.Name.De,
		NameEn: loc.Name.En,
		Provider: loc.Provider.Slug(),
	})
	if err != nil {
		return fmt.Errorf("insert location %s: %w", loc.Slug, err)
	}
	return nil
}

func (r *locationRepo) Delete(ctx context.Context, slug string) error {
	err := r.q.DeleteLocationBySlug(ctx, slug)
	if err != nil {
		return fmt.Errorf("delete location %s: %w", slug, err)
	}
	return nil
}

// Helper to map sqlc struct to domain
func mapLocToDomain(row db.Location) *meal.Location {
	return &meal.Location{
		Slug: row.Slug,
		Name: meal.LangString{
			De: row.NameDe,
			En: row.NameEn,
		},
	}
}

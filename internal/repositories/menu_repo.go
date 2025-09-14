package repositories

import (
	"context"
	"fmt"
	"hpi-mensa/gen/db"
	"hpi-mensa/internal/domain/meal"
	"time"
)

type MenuRepository interface {
	GetBySlug(ctx context.Context, slug string) (*meal.Menu, error)
	GetByDate(ctx context.Context, date time.Time) ([]*meal.Menu, error)
	GetByDateAndLocation(ctx context.Context, date time.Time, locationSlug string) (*meal.Menu, error)
	Create(ctx context.Context, menu *meal.Menu) error
	Update(ctx context.Context, menu *meal.Menu) error
	Delete(ctx context.Context, slug string) error
}

type menuRepo struct {
	q *db.Queries
}

func NewMenuRepository(queries *db.Queries) MenuRepository {
	return &menuRepo{
		q: queries,
	}
}

func (r *menuRepo) GetBySlug(ctx context.Context, slug string) (*meal.Menu, error) {
	row, err := r.q.GetMenuBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("get menu by slug %s: %w", slug, err)
	}
	
}

func mapMenuToDomain(row db.Menu) *meal.Menu {
	return &meal.Menu{
		Slug: row.Slug,
		Date: row.Date,
		Location: ,
	}
}

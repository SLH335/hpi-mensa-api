package http

import "github.com/slh335/hpi-mensa-api/services"

type Server struct {
	LocationService *services.LocationService
	MealService     *services.MealService
}

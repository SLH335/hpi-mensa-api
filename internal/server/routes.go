package server

import (
	"fmt"
	"hpi-mensa/internal/domain/meal"
	"hpi-mensa/internal/services"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"https://*", "http://*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Global request logger using Echo middleware + zerolog
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			stop := time.Now()

			if err != nil {
				c.Error(err)
			}

			req := c.Request()
			res := c.Response()

			log.Info().
				Str("method", req.Method).
				Str("path", req.URL.Path).
				Int("status", res.Status).
				Dur("duration", stop.Sub(start)).
				Str("ip", c.RealIP()).
				Msg("Handled request")

			return err
		}
	})

	e.GET("/", s.HelloWorldHandler)
	e.GET("/locations", s.LocationsHandler)
	e.GET("/menu/:location", s.MenuHandler)

	e.GET("/health", s.healthHandler)

	return e
}

func (s *Server) HelloWorldHandler(c echo.Context) error {
	log.Debug().Msg("HelloWorldHandler called")

	resp := map[string]string{
		"message": "Hello World",
	}

	return c.JSON(http.StatusOK, resp)
}

func (s *Server) LocationsHandler(c echo.Context) error {
	log.Info().Msg("Loading locations")

	locations, err := services.GetLocations()
	if err != nil {
		log.Error().Err(err).Msg("Failed to load locations")
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": fmt.Sprintf("Failed to load locations: %v", err),
		})
	}

	log.Info().Int("count", len(locations)).Msg("Loaded locations successfully")
	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"message": "Successfully loaded locations",
		"data": locations,
	})
}

func (s *Server) MenuHandler(c echo.Context) error {
	locationSlug := strings.TrimSpace(c.Param("location"))
	if strings.TrimSpace(locationSlug) == "" {
		log.Warn().Msg("Missing location parameter")
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "path parameter 'location' is required",
		})
	}
	dateStr := strings.TrimSpace(c.QueryParam("date"))
	var date time.Time
	var err error
	if dateStr == "" {
		date = time.Now()
		dateStr = date.Format("2006-01-02")
	} else {
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			log.Warn().Str("date", dateStr).Msg("Invalid date field")
			return c.JSON(http.StatusInternalServerError, map[string]any{
				"success": false,
				"message": fmt.Sprintf("'%s' is not a valid date. Required format is YYYY-MM-DD", dateStr),
			})
		}
	}

	log.Info().Str("location", locationSlug).Str("date", dateStr).Msg("Loading menu")

	locations, err := services.GetLocations()
	if err != nil {
		log.Error().Err(err).Str("location", locationSlug).Msg("Failed to load locations")
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": fmt.Sprintf("Failed to load menu: %v", err),
		})
	}
	location := meal.Location{}
	for _, loc := range locations {
		if loc.Slug == locationSlug {
			location = loc
			break
		}
	}
	if location.Slug == "" {
		log.Warn().Str("location", locationSlug).Msg("Requested location does not exist")
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": fmt.Sprintf("Location '%s' does not exist", locationSlug),
		})
	}

	menu, err := services.GetMenu(location, date)
	if err != nil {
		log.Error().
			Err(err).
			Str("location", locationSlug).
			Str("date", dateStr).
			Msg("Failed to load menu")
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": fmt.Sprintf("Failed to load menu: %v", err),
		})
	}

	log.Info().
		Str("location", locationSlug).
		Str("date", dateStr).
		Int("meals", len(menu.Meals)).
		Msg("Loaded menu successfully")

	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"message": fmt.Sprintf("Successfully loaded menu for %s", dateStr),
		"data": menu,
	})
}

func (s *Server) healthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, s.db.Health())
}

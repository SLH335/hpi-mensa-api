package server

import (
	"fmt"
	"hpi-mensa/internal/mealdata/common"
	"hpi-mensa/internal/mealdata/common/types"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"https://*", "http://*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	e.GET("/", s.HelloWorldHandler)
	e.GET("/locations", s.LocationsHandler)
	e.GET("/menu/:location", s.MenuHandler)

	e.GET("/health", s.healthHandler)

	return e
}

func (s *Server) HelloWorldHandler(c echo.Context) error {
	resp := map[string]string{
		"message": "Hello World",
	}

	return c.JSON(http.StatusOK, resp)
}

func (s *Server) LocationsHandler(c echo.Context) error {
	locations, err := mealdata.GetLocations()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": fmt.Sprintf("Failed to load locations: %v", err),
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"message": "Successfully loaded locations",
		"data": locations,
	})
}

func (s *Server) MenuHandler(c echo.Context) error {
	locationSlug := c.Param("location")
	if strings.TrimSpace(locationSlug) == "" {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": fmt.Sprintf("path parameter 'location' is required"),
		})
	}

	locations, err := mealdata.GetLocations()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": fmt.Sprintf("Failed to load menu: %v", err),
		})
	}
	location := common.Location{}
	for _, loc := range locations {
		if loc.Slug == locationSlug {
			location = loc
			break
		}
	}
	if location.Slug == "" {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": fmt.Sprintf("Location '%s' does not exist", locationSlug),
		})
	}

	meals, err := mealdata.GetMeals(location)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": fmt.Sprintf("Failed to load menu: %v", err),
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"message": "Successfully loaded menu",
		"data": meals,
	})
}

func (s *Server) healthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, s.db.Health())
}

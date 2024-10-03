package main

import (
	"log"

	"github.com/labstack/echo/v4"

	"github.com/slh335/hpi-mensa-api/database"
	"github.com/slh335/hpi-mensa-api/http"
	"github.com/slh335/hpi-mensa-api/services"
)

func main() {
	db, err := database.Open("file:app.db")
	if err != nil {
		log.Fatal(err)
		return
	}
	err = database.CreateTables(db)
	if err != nil {
		log.Fatal(err)
		return
	}

	server := http.Server{
		LocationService: &services.LocationService{
			DbService: &database.LocationDBService{
				DB: db,
			},
		},
		MealService: &services.MealService{
			DbService: &database.MealDBService{
				DB: db,
			},
			AttributeService: &services.MealAttributeService{
				DbService: &database.MealAttributeDBService{
					DB: db,
				},
			},
			CategoryService: &services.MealCategoryService{
				DbService: &database.MealCategoryDBService{
					DB: db,
				},
			},
		},
	}

	e := echo.New()

	e.Static("/static", "static")

	e.GET("/", server.Index)
	e.GET("/de", server.Index)
	e.GET("/en", server.Index)

	api := e.Group("/api/v1")
	api.GET("/locations", server.GetLocations)
	api.GET("/additives/:location", server.GetAdditives)
	api.GET("/allergens/:location", server.GetAllergens)
	api.GET("/features/:location", server.GetFeatures)
	api.GET("/meals/:location", server.GetMeals)

	e.Logger.Fatal(e.Start(":3000"))
}

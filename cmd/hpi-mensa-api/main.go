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

	e.GET("/locations", server.GetLocations)
	e.GET("/additives/:location", server.GetAdditives)
	e.GET("/allergens/:location", server.GetAllergens)
	e.GET("/features/:location", server.GetFeatures)
	e.GET("/meals/:location", server.GetMeals)

	e.Logger.Fatal(e.Start(":3000"))
}

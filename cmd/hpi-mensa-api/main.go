package main

import (
	"fmt"
	"log"
	"time"

	"github.com/slh335/hpi-mensa-api/database"
	"github.com/slh335/hpi-mensa-api/services"
	. "github.com/slh335/hpi-mensa-api/types"
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

	locations, err := mensadata.GetLocations()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(locations)

	fmt.Println(mensadata.GetMeals(locations[1], English, time.Now().Add(time.Hour*24)))
}

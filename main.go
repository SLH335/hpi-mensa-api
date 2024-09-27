package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/valyala/fastjson"
)

type Location int

const (
	Griebnitzsee Location = 9601
)

type Language int

const (
	German  Language = 1
	English Language = 2
)

type Meal struct {
	Name         string
	StudentPrice float64
	GuestPrice   float64
	Date         time.Time
	Id           int
	Nutrition    Nutrition
}

type Nutrition struct {
	Kj            int
	Kcal          int
	Fat           float64
	SaturatedFat  float64
	Carbohydrates float64
	Sugar         float64
	Protein       float64
	Salt          float64
}

const baseUrl string = "https://swp.webspeiseplan.de/index.php"

func getMeals(location Location, language Language, day time.Time) (meals []Meal, err error) {
	jsonData, err := getMenuData(location, language)
	if err != nil {
		return []Meal{}, err
	}

	for _, mealplan := range jsonData.GetArray() {
		if string(mealplan.GetStringBytes("speiseplanAdvanced", "titel")) != "Mittagessen" {
			continue
		}
		for _, mealList := range mealplan.GetArray("speiseplanGerichtData") {
			meal := Meal{}

			meal.Date, err = time.Parse(time.RFC3339, string(mealList.GetStringBytes("speiseplanAdvancedGericht", "datum")))
			if err != nil {
				return meals, err
			}
			if meal.Date.Year() != day.Year() || meal.Date.YearDay() != day.YearDay() {
				continue
			}

			meal.Name = string(mealList.GetStringBytes("speiseplanAdvancedGericht", "gerichtname"))
			meal.StudentPrice = mealList.GetFloat64("zusatzinformationen", "mitarbeiterpreisDecimal2")
			meal.GuestPrice = mealList.GetFloat64("zusatzinformationen", "gaestepreisDecimal2")
			meal.Id = mealList.GetInt("speiseplanAdvancedGericht", "id")
			meal.Nutrition.Kj = mealList.GetInt("zusatzinformationen", "nwkjInteger")
			meal.Nutrition.Kcal = mealList.GetInt("zusatzinformationen", "nwkcalInteger")
			meal.Nutrition.Fat = mealList.GetFloat64("zusatzinformationen", "nwfettDecimal1")
			meal.Nutrition.SaturatedFat = mealList.GetFloat64("zusatzinformationen", "nwfettsaeurenDecimal1")
			meal.Nutrition.Carbohydrates = mealList.GetFloat64("zusatzinformationen", "nwkohlehydrateDecimal1")
			meal.Nutrition.Sugar = mealList.GetFloat64("zusatzinformationen", "nwzuckerDecimal1")
			meal.Nutrition.Protein = mealList.GetFloat64("zusatzinformationen", "nweiweissDecimal1")
			meal.Nutrition.Salt = mealList.GetFloat64("zusatzinformationen", "nwsalzDecimal1")
			meals = append(meals, meal)
		}
	}

	return meals, nil
}

func getMenuData(location Location, language Language) (jsonData *fastjson.Value, err error) {
	params := url.Values{}
	params.Add("token", "55ed21609e26bbf68ba2b19390bf7961")
	params.Add("model", "menu")
	params.Add("location", strconv.Itoa(int(location)))
	params.Add("languagetype", strconv.Itoa(int((language))))
	params.Add("_", fmt.Sprintf("%d", time.Now().UnixMilli()))

	req, err := http.NewRequest(http.MethodGet, baseUrl+"?"+params.Encode(), nil)
	if err != nil {
		return jsonData, err
	}

	req.Header.Add("Referer", "https://sqp.webspeiseplan.de/Menu")

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return jsonData, err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return jsonData, err
	}

	var parser fastjson.Parser
	jsonData, err = parser.Parse(string(resBody))
	if err != nil {
		return jsonData, err
	}

	return jsonData.Get("content"), nil
}

func main() {
	fmt.Println(getMeals(Griebnitzsee, German, time.Now()))
}

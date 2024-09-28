package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/valyala/fastjson"
)

const baseUrl string = "https://swp.webspeiseplan.de/index.php"

type Location int

const (
	None         Location = 0
	Griebnitzsee Location = 9601
)

type Language int

const (
	German  Language = 1
	English Language = 2
)

type Model string
type FeatureModel Model

const (
	Additives FeatureModel = "additives"
	Allergens FeatureModel = "allergens"
	Features  FeatureModel = "features"
	Locations Model        = "location"
	Menu      Model        = "menu"
	Outlets   Model        = "outlet"
)

type MealType struct {
	Id     int
	NameDe string
	NameEn string
}

var UnknownOffer MealType = MealType{
	Id:     0,
	NameDe: "Unbekanntes Angebot",
	NameEn: "Unknown Offer",
}
var Offer1 MealType = MealType{
	Id:     149,
	NameDe: "Angebot 1",
	NameEn: "Offer 1",
}
var Offer2 MealType = MealType{
	Id:     150,
	NameDe: "Angebot 2",
	NameEn: "Offer 2",
}
var Offer3 MealType = MealType{
	Id:     151,
	NameDe: "Angebot 3",
	NameEn: "Offer 3",
}
var Offer4 MealType = MealType{
	Id:     152,
	NameDe: "Angebot 4",
	NameEn: "Offer 4",
}
var DailySpecial MealType = MealType{
	Id:     118,
	NameDe: "Tagesangebot",
	NameEn: "Daily Special",
}
var SaladBar MealType = MealType{
	Id:     112,
	NameDe: "Salattheke",
	NameEn: "Salad bar",
}
var Dessert1 MealType = MealType{
	Id:     119,
	NameDe: "Dessert 1",
	NameEn: "Dessert 1",
}
var Dessert2 MealType = MealType{
	Id:     120,
	NameDe: "Dessert 1",
	NameEn: "Dessert 2",
}

type Meal struct {
	Name         string
	Type         MealType
	StudentPrice float64
	GuestPrice   float64
	Date         time.Time
	Id           int
	Nutrition    Nutrition
	Additives    []Feature
	Allergens    []Feature
	Features     []Feature
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

type Feature struct {
	Id    int
	Name  string
	Short string
}

func getMeals(location Location, language Language, day time.Time) (meals []Meal, err error) {
	jsonData, err := getData(Menu, location, language)
	if err != nil {
		return []Meal{}, err
	}
	additives, err := getFeatures(Additives, location, language)
	if err != nil {
		return []Meal{}, err
	}
	allergens, err := getFeatures(Allergens, location, language)
	if err != nil {
		return []Meal{}, err
	}
	features, err := getFeatures(Features, location, language)
	if err != nil {
		return []Meal{}, err
	}

	for _, mealplan := range jsonData.GetArray() {
		if string(mealplan.GetStringBytes("speiseplanAdvanced", "titel")) != "Mittagessen" {
			continue
		}
		for _, mealData := range mealplan.GetArray("speiseplanGerichtData") {
			meal := Meal{}

			meal.Date, err = time.Parse(time.RFC3339, string(mealData.GetStringBytes("speiseplanAdvancedGericht", "datum")))
			if err != nil {
				return meals, err
			}
			if meal.Date.UTC().Year() != day.UTC().Year() || meal.Date.UTC().YearDay() != day.UTC().YearDay() {
				continue
			}

			meal.Name = string(mealData.GetStringBytes("speiseplanAdvancedGericht", "gerichtname"))
			meal.StudentPrice = mealData.GetFloat64("zusatzinformationen", "mitarbeiterpreisDecimal2")
			meal.GuestPrice = mealData.GetFloat64("zusatzinformationen", "gaestepreisDecimal2")
			meal.Id = mealData.GetInt("speiseplanAdvancedGericht", "id")
			meal.Nutrition.Kj = mealData.GetInt("zusatzinformationen", "nwkjInteger")
			meal.Nutrition.Kcal = mealData.GetInt("zusatzinformationen", "nwkcalInteger")
			meal.Nutrition.Fat = mealData.GetFloat64("zusatzinformationen", "nwfettDecimal1")
			meal.Nutrition.SaturatedFat = mealData.GetFloat64("zusatzinformationen", "nwfettsaeurenDecimal1")
			meal.Nutrition.Carbohydrates = mealData.GetFloat64("zusatzinformationen", "nwkohlehydrateDecimal1")
			meal.Nutrition.Sugar = mealData.GetFloat64("zusatzinformationen", "nwzuckerDecimal1")
			meal.Nutrition.Protein = mealData.GetFloat64("zusatzinformationen", "nweiweissDecimal1")
			meal.Nutrition.Salt = mealData.GetFloat64("zusatzinformationen", "nwsalzDecimal1")

			additiveIds := strings.Split(string(mealData.GetStringBytes("zusatzstoffeIds")), ",")
			allergenIds := strings.Split(string(mealData.GetStringBytes("allergeneIds")), ",")
			featureIds := strings.Split(string(mealData.GetStringBytes("gerichtmerkmaleIds")), ",")

			for _, additiveId := range additiveIds {
				for _, additive := range additives {
					if strconv.Itoa(additive.Id) == additiveId {
						meal.Additives = append(meal.Additives, additive)
						break
					}
				}
			}
			for _, allergenId := range allergenIds {
				for _, allergen := range allergens {
					if strconv.Itoa(allergen.Id) == allergenId {
						meal.Allergens = append(meal.Allergens, allergen)
						break
					}
				}
			}
			for _, featureId := range featureIds {
				for _, feature := range features {
					if strconv.Itoa(feature.Id) == featureId {
						meal.Features = append(meal.Features, feature)
						break
					}
				}
			}

			meal.Type = getMealTypeFromId(mealData.GetInt("speiseplanAdvancedGericht", "gerichtkategorieID"))

			meals = append(meals, meal)
		}
	}

	return meals, nil
}

func getFeatures(model FeatureModel, location Location, language Language) (features []Feature, err error) {
	jsonData, err := getData(Model(model), location, language)
	if err != nil {
		return features, err
	}

	for _, featureData := range jsonData.GetArray() {
		feature := Feature{}
		feature.Id = featureData.GetInt("id")
		feature.Name = string(featureData.GetStringBytes("name"))
		feature.Short = string(featureData.GetStringBytes("kuerzel"))
		features = append(features, feature)
	}
	return features, nil
}

func getMealTypeFromId(id int) MealType {
	switch id {
	case Offer1.Id:
		return Offer1
	case Offer2.Id:
		return Offer2
	case Offer3.Id:
		return Offer3
	case Offer4.Id:
		return Offer4
	case DailySpecial.Id:
		return DailySpecial
	case SaladBar.Id:
		return SaladBar
	case Dessert1.Id:
		return Dessert1
	case Dessert2.Id:
		return Dessert2
	default:
		return UnknownOffer
	}
}

func getData(model Model, location Location, language Language) (jsonData *fastjson.Value, err error) {
	params := url.Values{}
	params.Add("token", "55ed21609e26bbf68ba2b19390bf7961")
	params.Add("model", string(model))
	params.Add("location", strconv.Itoa(int(location)))
	params.Add("languagetype", strconv.Itoa(int((language))))

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
	fmt.Println(getMeals(Griebnitzsee, German, time.Now().Add(time.Hour*24)))
}

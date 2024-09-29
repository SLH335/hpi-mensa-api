package mensadata

import (
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/valyala/fastjson"

	. "github.com/slh335/hpi-mensa-api/types"
)

const baseUrl string = "https://swp.webspeiseplan.de/index.php"

func GetMeals(location Location, language Language, date time.Time) (meals []Meal, err error) {
	jsonData, err := getData(MenuModel, location, language)
	if err != nil {
		return []Meal{}, err
	}

	meals = []Meal{}
	for _, mealplan := range jsonData.GetArray() {
		if string(mealplan.GetStringBytes("speiseplanAdvanced", "titel")) != "Mittagessen" {
			continue
		}
		for _, mealData := range mealplan.GetArray("speiseplanGerichtData") {
			var meal Meal

			meal.Date, err = time.Parse(time.RFC3339, string(mealData.GetStringBytes("speiseplanAdvancedGericht", "datum")))
			if err != nil {
				return meals, err
			}
			if meal.Date.Format("2006-01-02") != date.Format("2006-01-02") {
				continue
			}

			meal.Id = mealData.GetInt("speiseplanAdvancedGericht", "id")
			meal.Location = location
			meal.NameDe = string(mealData.GetStringBytes("speiseplanAdvancedGericht", "gerichtname"))
			meal.NameEn = string(mealData.GetStringBytes("zusatzinformationen", "gerichtnameAlternative"))
			meal.StudentPrice = mealData.GetFloat64("zusatzinformationen", "mitarbeiterpreisDecimal2")
			meal.GuestPrice = mealData.GetFloat64("zusatzinformationen", "gaestepreisDecimal2")
			meal.Nutrition.Kj = mealData.GetInt("zusatzinformationen", "nwkjInteger")
			meal.Nutrition.Kcal = mealData.GetInt("zusatzinformationen", "nwkcalInteger")
			meal.Nutrition.Fat = mealData.GetFloat64("zusatzinformationen", "nwfettDecimal1")
			meal.Nutrition.SaturatedFat = mealData.GetFloat64("zusatzinformationen", "nwfettsaeurenDecimal1")
			meal.Nutrition.Carbohydrates = mealData.GetFloat64("zusatzinformationen", "nwkohlehydrateDecimal1")
			meal.Nutrition.Sugar = mealData.GetFloat64("zusatzinformationen", "nwzuckerDecimal1")
			meal.Nutrition.Protein = mealData.GetFloat64("zusatzinformationen", "nweiweissDecimal1")
			meal.Nutrition.Salt = mealData.GetFloat64("zusatzinformationen", "nwsalzDecimal1")

			additiveIds := strings.Split(string(mealData.GetStringBytes("zusatzstoffeIds")), ",")
			for _, additiveIdStr := range additiveIds {
				additiveId, err := strconv.Atoi(additiveIdStr)
				if err != nil {
					continue
				}
				meal.Additives = append(meal.Additives, MealAttribute{
					Id:   additiveId,
					Type: AdditiveAttribute,
				})
			}

			allergenIds := strings.Split(string(mealData.GetStringBytes("allergeneIds")), ",")
			for _, allergenIdStr := range allergenIds {
				allergenId, err := strconv.Atoi(allergenIdStr)
				if err != nil {
					continue
				}
				meal.Allergens = append(meal.Allergens, MealAttribute{
					Id:   allergenId,
					Type: AllergenAttribute,
				})
			}

			featureIds := strings.Split(string(mealData.GetStringBytes("gerichtmerkmaleIds")), ",")
			for _, featureIdStr := range featureIds {
				featureId, err := strconv.Atoi(featureIdStr)
				if err != nil {
					continue
				}
				meal.Features = append(meal.Features, MealAttribute{
					Id:   featureId,
					Type: FeatureAttribute,
				})
			}

			meal.Category = MealCategory{
				Id: mealData.GetInt("speiseplanAdvancedGericht", "gerichtkategorieID"),
			}

			meals = append(meals, meal)
		}
	}

	return meals, nil
}

func GetMealAttributes(attributeType MealAttributeType, location Location, language Language) (attributes []MealAttribute, err error) {
	jsonData, err := getData(Model(attributeType+"s"), location, language)
	if err != nil {
		return attributes, err
	}

	for _, attributeData := range jsonData.GetArray() {
		attribute := MealAttribute{}
		switch attributeType {
		case AdditiveAttribute:
			attribute.Id = attributeData.GetInt("zusatzstoffeID")
		case AllergenAttribute:
			attribute.Id = attributeData.GetInt("allergeneID")
		case FeatureAttribute:
			attribute.Id = attributeData.GetInt("gerichtmerkmalID")
		}
		if language == German {
			attribute.NameDe = string(attributeData.GetStringBytes("name"))
		} else {
			attribute.NameEn = string(attributeData.GetStringBytes("name"))
		}
		attribute.Short = string(attributeData.GetStringBytes("kuerzel"))
		attribute.Type = attributeType
		attribute.Location = &location
		attributes = append(attributes, attribute)
	}
	return attributes, nil
}

func GetLocations() (locations []Location, err error) {
	jsonData, err := getData(LocationsModel, Location{}, Language(0))
	if err != nil {
		return []Location{}, err
	}

	for _, locationData := range jsonData.GetArray() {
		var location Location
		location.Id = locationData.GetInt("id")
		location.Name = string(locationData.GetStringBytes("name"))
		locations = append(locations, location)
	}
	return locations, nil
}

func GetMealCategories(location Location, lang Language) (categories []MealCategory, err error) {
	jsonData, err := getData(MealCategoryModel, location, lang)
	if err != nil {
		return []MealCategory{}, err
	}

	for _, categoryData := range jsonData.GetArray() {
		var category MealCategory
		category.Id = categoryData.GetInt("gerichtkategorieID")
		category.Name = string(categoryData.GetStringBytes("name"))
		category.Location = &location
		categories = append(categories, category)
	}
	return categories, nil
}

func getData(model Model, location Location, language Language) (jsonData *fastjson.Value, err error) {
	params := url.Values{}
	params.Add("token", "55ed21609e26bbf68ba2b19390bf7961")
	params.Add("model", string(model))
	params.Add("location", strconv.Itoa(location.Id))
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

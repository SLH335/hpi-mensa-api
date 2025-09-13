package stwwb

import (
	"errors"
	"fmt"
	"hpi-mensa/internal/mealdata/util"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/valyala/fastjson"
)

// Get all meals of current menu (usually current and following week) for the given location
func getMenu(location Location) (menu []Meal, err error) {
	jsonData, err := stwwbRequest(MenuModel, location, LangGerman)
	if err != nil {
		return []Meal{}, fmt.Errorf("menu request: %w", err)
	}

	for _, plan := range jsonData.GetArray() {
		outletId := plan.GetInt("speiseplanAdvanced", "outletID")
		outlet := outlets[outletId]

		for _, meal := range plan.GetArray("speiseplanGerichtData") {
			dishData := meal.Get("speiseplanAdvancedGericht")
			extraData := meal.Get("zusatzinformationen")

			// Fill in meal category data from global state
			categoryID := dishData.GetInt("gerichtkategorieID")
			mealCategory := categories[location.ID][categoryID]
			// Abort if category was not found
			if mealCategory.ID == 0 {
				return []Meal{}, fmt.Errorf("meal category not found: '%d'", categoryID)
			}

			dateStr := string(dishData.GetStringBytes("datum"))
			if len(dateStr) >= 10 {
				dateStr = dateStr[:10]
			}
			date, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				return []Meal{}, fmt.Errorf("parse meal date: %w", err)
			}

			// Fill in meal attribute data from global state
			allergens, err := AllergenAttribute.getAttributeData(string(meal.GetStringBytes("allergeneIds")), location)
			if err != nil {
				return []Meal{}, fmt.Errorf("get allergen data: %w", err)
			}
			additives, err := AdditiveAttribute.getAttributeData(string(meal.GetStringBytes("zusatzstoffeIds")), location)
			if err != nil {
				return []Meal{}, fmt.Errorf("get additive data: %w", err)
			}
			features, err := FeatureAttribute.getAttributeData(string(meal.GetStringBytes("gerichtmerkmaleIds")), location)
			if err != nil {
				return []Meal{}, fmt.Errorf("get feature data: %w", err)
			}

			menu = append(menu, Meal{
				ID: dishData.GetInt("id"),
				Name: util.LangString{
					De: string(dishData.GetStringBytes("gerichtname")),
					En: string(extraData.GetStringBytes("gerichtnameAlternative")),
				},
				Category: mealCategory,
				Date: date,
				IsActive: dishData.GetBool("aktiv"),
				PriceStudent: extraData.GetFloat64("mitarbeiterpreisDecimal2"),
				PriceEmployee: extraData.GetFloat64("gaestepreisDecimal2"),
				PriceGuest: extraData.GetFloat64("price3Decimal2"),
				Contingent: extraData.GetInt("contingent"),
				Nutrition: MealNutrition{
					Kj: extraData.GetInt("nwkjInteger"),
					Kcal: extraData.GetInt("nwkcalInteger"),
					Fat: extraData.GetFloat64("nwfettDecimal1"),
					SaturatedFat: extraData.GetFloat64("nwfettsaeurenDecimal1"),
					Carbs: extraData.GetFloat64("nwkohlehydrateDecimal1"),
					Sugar: extraData.GetFloat64("nwzuckerDecimal1"),
					Protein: extraData.GetFloat64("nweiweissDecimal1"),
					Salt: extraData.GetFloat64("nwsalzDecimal1"),
				},
				Allergens: allergens,
				Additives: additives,
				Features: features,
				Sustainability: MealSustainability{
					CO2: MealCO2{
						Grams: extraData.GetInt("sustainability", "co2", "co2Value"),
						Rating: string(extraData.GetStringBytes("sustainability", "co2", "co2RatingIdentifier")),
					},
				},
				Location: location,
				Outlet: outlet,
			})
		}
	}

	return menu, nil
}

func getLocations() (locations []Location, err error) {
	jsonData, err := stwwbRequest(LocationsModel, Location{}, LangGerman)
	if err != nil {
		return []Location{}, fmt.Errorf("locations request: %w", err)
	}

	for _, location := range jsonData.GetArray() {
		isActive := location.GetBool("active")
		isPublic := location.GetBool("isPublic")
		if !isActive || !isPublic {
			continue
		}

		locations = append(locations, Location{
			ID: location.GetInt("id"),
			Name: string(location.GetStringBytes("name")),
		})
	}

	return locations, nil
}

func getMealCategories(location Location) (categories map[int]MealCategory, err error) {
	categories = map[int]MealCategory{}

	jsonDataGer, err := stwwbRequest(MealCategoryModel, location, LangGerman)
	if err != nil {
		return map[int]MealCategory{}, fmt.Errorf("meal categories request: %w", err)
	}

	// Get categories for English separately, because the name is only returned in the specified language
	jsonDataEng, err := stwwbRequest(MealCategoryModel, location, LangEnglish)
	if err != nil {
		return map[int]MealCategory{}, fmt.Errorf("meal categories request: %w", err)
	}

	// Store English names in map for faster access
	englishNames := map[int]string{}
	for _, category := range jsonDataEng.GetArray() {
		id := category.GetInt("gerichtkategorieID")
		name := string(category.GetStringBytes("name"))
		englishNames[id] = name
	}

	// Iterate through German categories and add English names from map
	for _, category := range jsonDataGer.GetArray() {
		id := category.GetInt("gerichtkategorieID")
		categories[id] = MealCategory{
			ID: id,
			Name: util.LangString{
				De: string(category.GetStringBytes("name")),
				En: englishNames[id],
			},
		}
	}

	return categories, nil
}

func getMealAttributes(location Location, attributeType MealAttributeType) (attributes map[int]MealAttribute, err error) {
	attributes = map[int]MealAttribute{}

	model, err := attributeType.getAttributeModel()
	if err != nil {
		return map[int]MealAttribute{}, fmt.Errorf("get attribute model: %w", err)
	}

	jsonDataGer, err := stwwbRequest(model, location, LangGerman)
	if err != nil {
		return map[int]MealAttribute{}, fmt.Errorf("meal attributes request: %w", err)
	}

	// Get attributes for English separately, because the name is only returned in the specified language
	jsonDataEng, err := stwwbRequest(model, location, LangEnglish)
	if err != nil {
		return map[int]MealAttribute{}, fmt.Errorf("meal attributes request: %w", err)
	}

	// Get correct id key, since they are different for each attribute type
	idKey, err := attributeType.getAttributeIDKey()
	if err != nil {
		return map[int]MealAttribute{}, fmt.Errorf("get attribute id key: %w", err)
	}

	// Store English names in map for faster access
	englishNames := map[int]string{}
	for _, attribute := range jsonDataEng.GetArray() {
		id := attribute.GetInt(idKey)
		name := string(attribute.GetStringBytes("name"))
		englishNames[id] = name
	}

	// Iterate through German attributes and add English names from map
	for _, attribute := range jsonDataGer.GetArray() {
		id := attribute.GetInt(idKey)
		attributes[id] = MealAttribute{
			ID: id,
			Name: util.LangString{
				De: string(attribute.GetStringBytes("name")),
				En: englishNames[id],
			},
			Short: string(attribute.GetStringBytes("kuerzel")),
		}
	}

	return attributes, nil
}

func getOutlets() (outlets map[int]Outlet, err error) {
	outlets = map[int]Outlet{}

	jsonData, err := stwwbRequest(LocationsModel, Location{}, LangGerman)
	if err != nil {
		return map[int]Outlet{}, fmt.Errorf("outlets request: %w", err)
	}

	for _, outlet := range jsonData.GetArray() {
		id := outlet.GetInt("id")
		outlets[id] = Outlet{
			ID: id,
			Name: string(outlet.GetStringBytes("name")),
			Location: Location{
				ID: outlet.GetInt("locationInfo", "id"),
				Name: string(outlet.GetStringBytes("locationInfo", "name")),
			},
			OpeningHours: parseOutletOpeningHours(outlet),
			Address: OutletAddress{
				ID: outlet.GetInt("addressInfo", "id"),
				Street: string(outlet.GetStringBytes("addressInfo", "street")),
				City: string(outlet.GetStringBytes("addressInfo", "city")),
				PostalCode: string(outlet.GetStringBytes("addressInfo", "postalCode")),
				Country: string(outlet.GetStringBytes("addressInfo", "countryName")),
				CountryCode: string(outlet.GetStringBytes("addressInfo", "countryCode")),
				Coords: Coords{
					Lat: outlet.GetFloat64("positionInfo", "latitude"),
					Lon: outlet.GetFloat64("positionInfo", "longitude"),
				},
			},
			ContactEmail: string(outlet.GetStringBytes("contactInfo", "email")),
			ImageLink: string(outlet.GetStringBytes("logoImage")),
		}
	}

	return outlets, nil
}

func stwwbRequest(model Model, location Location, language Language) (jsonData *fastjson.Value, err error) {
	baseUrl := os.Getenv("STWWB_BASE_URL")
	token := os.Getenv("STWWB_TOKEN")

	params := url.Values{}
	params.Add("token", token)
	params.Add("model", string(model))
	params.Add("location", strconv.Itoa(location.ID))
	params.Add("languagetype", strconv.Itoa(int(language)))

	url := fmt.Sprintf("%s/index.php?%s", baseUrl, params.Encode())
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create STWWB request: %w", err)
	}

	req.Header.Add("Referer", baseUrl+"/menu")

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send STWWB request: %w", err)
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read STWWB response body: %w", err)
	}

	var parser fastjson.Parser
	jsonData, err = parser.Parse(string(resBody))
	if err != nil {
		return nil, fmt.Errorf("parse STWWB response data: %w", err)
	}

	success := jsonData.GetBool("success")
	if !success {
		errorMessage := "Unknown error"
		errorMessageBytes := jsonData.GetStringBytes("content")
		if errorMessageBytes != nil {
			errorMessage = string(errorMessageBytes)
		}
		return nil, fmt.Errorf("STWWB request did not succeed: %w", errors.New(errorMessage))
	}

	return jsonData.Get("content"), nil
}


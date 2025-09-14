package stwwb

import (
	"errors"
	"fmt"
	"hpi-mensa/internal/mealdata/common/types"
	"hpi-mensa/internal/mealdata/util"
	"strconv"
	"strings"

	"github.com/valyala/fastjson"
)

func (location Location) Slug() (slug string) {
	return Provider.Slug() + "-" + strings.ToLower(strings.ReplaceAll(location.Name, " ", "-"))
}

// Helper function to parse STWWB JSON response data for outlet opening hours
func parseOutletOpeningHours(jsonData *fastjson.Value) (openingHours OutletOpeningHours) {
	return OutletOpeningHours{
		Mo: parseOutletOpeningTime(jsonData, "mo"),
		Tu: parseOutletOpeningTime(jsonData, "di"),
		We: parseOutletOpeningTime(jsonData, "mi"),
		Th: parseOutletOpeningTime(jsonData, "do"),
		Fr: parseOutletOpeningTime(jsonData, "fr"),
		Sa: parseOutletOpeningTime(jsonData, "sa"),
		Su: parseOutletOpeningTime(jsonData, "so"),
	}
}

func parseOutletOpeningTime(jsonData *fastjson.Value, day string) (openingTime OutletOpeningTime) {
	text := string(jsonData.GetStringBytes(day+"Zeit2"))
	hoursStr := strings.ReplaceAll(string(jsonData.GetStringBytes(day+"Zeit1")), " Uhr", "")
	hours := strings.Split(hoursStr, " - ")
	opening, closing := "", ""
	if len(hours) >= 2 {
		opening = hours[0]
		closing = hours[1]
	}

	return OutletOpeningTime{
		Text: text,
		Opening: opening,
		Closing: closing,
	}
}

func convertPriceCategories(meal Meal) (prices []common.PriceCategory) {
	if meal.PriceStudent > 0 {
		prices = append(prices, common.PriceCategory{
			Type: util.LangString{
				De: "Studi",
				En: "Student",
			},
			Amount: meal.PriceStudent,
		})
	}
	if meal.PriceEmployee > 0 {
		prices = append(prices, common.PriceCategory{
			Type: util.LangString{
				De: "Mitarbeiter",
				En: "Employee",
			},
			Amount: meal.PriceEmployee,
		})
	}
	if meal.PriceGuest > 0 {
		prices = append(prices, common.PriceCategory{
			Type: util.LangString{
				De: "Gast",
				En: "Guest",
			},
			Amount: meal.PriceGuest,
		})
	}

	if len(prices) == 1 {
		prices[0].Type = util.LangString{
			De: "Alle",
			En: "All",
		}
	}

	return prices
}

func convertMealAttributes(mealAttributes []MealAttribute) (attributes []common.MealAttribute) {
	for _, mealAttribute := range mealAttributes {
		// Skip attributes like 'no allergen' or 'no additive'
		if strings.TrimSpace(mealAttribute.Short) == "" {
			continue
		}
		attributes = append(attributes, common.MealAttribute{
			Slug: fmt.Sprintf("%s-%d", Provider.Slug(), mealAttribute.ID),
			Name: mealAttribute.Name,
			Short: util.LangString{
				De: mealAttribute.Short,
				En: mealAttribute.Short,
			},
		})
	}
	return attributes
}

func (attributeType MealAttributeType) getAttributeIDKey() (key string, err error) {
	switch attributeType {
	case AllergenAttribute:
		return "allergeneID", nil
	case AdditiveAttribute:
		return "zusatzstoffeID", nil
	case FeatureAttribute:
		return "gerichtmerkmalID", nil
	default:
		return "", errors.New(fmt.Sprintf("invalid attribute type '%s'", attributeType))
	}
}

func (attributeType MealAttributeType) getAttributeModel() (key Model, err error) {
	switch attributeType {
	case AllergenAttribute:
		return AllergensModel, nil
	case AdditiveAttribute:
		return AdditivesModel, nil
	case FeatureAttribute:
		return FeaturesModel, nil
	default:
		return "", errors.New(fmt.Sprintf("invalid attribute type '%s'", attributeType))
	}
}

func (attributeType MealAttributeType) getAttributeData(attributeIDStr string, location Location) (attributes []MealAttribute, err error) {
	if attributeIDStr == "" {
		return []MealAttribute{}, nil
	}
	attributeIDs := strings.SplitSeq(attributeIDStr, ",")
	for attributeIDStr := range attributeIDs {
		attributeID, err := strconv.Atoi(attributeIDStr)
		if err != nil {
			return []MealAttribute{}, fmt.Errorf("convert attribute id '%s' to int: %w", attributeIDStr, err)
		}

		attribute := MealAttribute{}
		switch attributeType {
		case AllergenAttribute:
			attribute = Provider.allergens[location.ID][attributeID]
		case AdditiveAttribute:
			attribute = Provider.additives[location.ID][attributeID]
		case FeatureAttribute:
			attribute = Provider.features[location.ID][attributeID]
		default:
			return []MealAttribute{}, errors.New(fmt.Sprintf("invalid attribute type: '%d'", attributeID))
		}

		// Abort if category was not found
		if attribute.ID == 0 {
			return []MealAttribute{}, fmt.Errorf("meal attribute not found: '%d'", attributeID)
		}

		attributes = append(attributes, attribute)
	}

	return attributes, nil
}

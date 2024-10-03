package services

import (
	"slices"

	"github.com/slh335/hpi-mensa-api/database"
	"github.com/slh335/hpi-mensa-api/services/mensadata"
	. "github.com/slh335/hpi-mensa-api/types"
)

type MealAttributeService struct {
	DbService *database.MealAttributeDBService
}

func (s *MealAttributeService) Get(
	attributeType MealAttributeType,
	location Location,
	lang Language,
) (attributes []MealAttribute, err error) {
	attributes, err = s.DbService.Get(attributeType, location)
	if err != nil {
		return []MealAttribute{}, err
	}

	if len(attributes) == 0 {
		attributesDe, err := mensadata.GetMealAttributes(attributeType, location, German)
		if err != nil {
			return []MealAttribute{}, err
		}
		attributesEn, err := mensadata.GetMealAttributes(attributeType, location, English)
		if err != nil {
			return []MealAttribute{}, err
		}
		for _, attributeDe := range attributesDe {
			// filter out empty attributes like 'no allergen'
			noAttribute := []string{"", " ", "NON", "X99", "Additives"}
			if slices.Contains(noAttribute, attributeDe.Short) {
				continue
			}
			for _, attributeEn := range attributesEn {
				if attributeDe.Id == attributeEn.Id {
					attribute := attributeEn
					attribute.NameDe = attributeDe.NameDe
					attributes = append(attributes, attribute)
					break
				}
			}
		}
		if location.Id == 999 {
			for _, attribute := range mensadata.GetCKOAttributes() {
				if attribute.Type == attributeType {
					attributes = append(attributes, attribute)
				}
			}
		}
		err = s.DbService.Add(attributes)
		if err != nil {
			return []MealAttribute{}, err
		}
	}

	for i := range attributes {
		attributes[i] = attributes[i].Translated(lang)
	}

	return attributes, nil
}

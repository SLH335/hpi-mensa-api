package mensadata

import (
	"net/http"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	. "github.com/slh335/hpi-mensa-api/types"
)

const ckoUrl string = "https://www.campus-kitchen-one.de"

var ckoLocation Location = Location{
	Id:   999,
	Name: "Campus Kitchen One",
}

func getCKOData(lang Language) (doc *goquery.Document, err error) {
	url := ckoUrl
	if lang == English {
		url += "/en"
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return doc, err
	}
	req.Header.Add("User-Agent", "HPI-Mensa-API v0.1.0")

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return doc, err
	}
	defer res.Body.Close()

	doc, err = goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return doc, err
	}

	return doc, nil
}

func GetCKOMeals(date time.Time, lang Language) (meals []Meal, err error) {
	doc, err := getCKOData(lang)
	if err != nil {
		return []Meal{}, err
	}

	// the HTML structure of the German and English pages are completely different
	switch lang {
	case German:
		meals, err = getCKOMealsGerman(doc)
	case English:
		meals, err = getCKOMealsEnglish(doc)
	}
	if err != nil {
		return []Meal{}, err
	}

	meals = slices.DeleteFunc(meals, func(meal Meal) bool {
		return meal.Date.Format("2006-01-02") != date.Format("2006-01-02")
	})

	return meals, nil
}

func getCKOMealsGerman(doc *goquery.Document) (meals []Meal, err error) {
	data := doc.Find("[id^=cc-m-] [id^=cc-matrix-]").FilterFunction(func(i int, s *goquery.Selection) bool {
		return s.Find(".j-module").Length() > 5
	})

	// parse prices first
	prices := [][][]float64{}
	data.Last().Find(".j-text").Each(func(i int, s *goquery.Selection) {
		if strings.HasPrefix(strings.TrimSpace(s.Text()), "Stud") {
			return
		}
		dayPrices := [][]float64{}

		s.Find("p span").Each(func(j int, s2 *goquery.Selection) {
			mealPrices := []float64{}
			pricesStr := strings.Split(strings.TrimSpace(s2.Text()), "|")
			for _, priceStr := range pricesStr {
				price, err := parsePrice(priceStr)
				if err != nil {
					return
				}
				mealPrices = append(mealPrices, price)
			}
			dayPrices = append(dayPrices, mealPrices)
		})
		prices = append(prices, dayPrices)
	})

	// parse dates and meals
	var currentDate time.Time // keep track of the current date
	data.First().Find(".j-text").Each(func(i int, s *goquery.Selection) {
		if i%2 == 0 {
			pattern := regexp.MustCompile("[0-9]{2}\\.[0-9]{2}\\.[0-9]{4}")
			dateStr := pattern.FindString(strings.TrimSpace(s.Text()))
			if dateStr != "" {
				currentDate, err = time.Parse("02.01.2006", dateStr)
			}
		} else {
			j := 0
			s.Find("p > span").Each(func(_ int, mealS *goquery.Selection) {
				name := strings.Join(strings.Fields(mealS.Text()), " ")
				if slices.Contains([]string{"", "-", "Geschlossen"}, name) {
					return
				}
				attributes, name := parseCKOAttributes(name)
				if name == "" {
					return
				}

				var meal Meal
				meal.NameDe = name
				dateInt, err := strconv.Atoi(currentDate.Format("20060102"))
				if err != nil {
					return
				}
				meal.Id = dateInt*100 + j
				meal.Category.Id = 0
				meal.Date = currentDate
				meal.StudentPrice = prices[i/2][j][0]
				meal.GuestPrice = prices[i/2][j][1]
				meal.Location = ckoLocation
				for _, short := range attributes {
					attribute := GetAttribute(short)
					switch attribute.Type {
					case AllergenAttribute:
						meal.Allergens = append(meal.Allergens, attribute)
					case AdditiveAttribute:
						meal.Additives = append(meal.Additives, attribute)
					case FeatureAttribute:
						meal.Features = append(meal.Features, attribute)
					}
				}
				meals = append(meals, meal)
				j++
			})
		}
	})

	return meals, nil
}

func getCKOMealsEnglish(doc *goquery.Document) (meals []Meal, err error) {
	// parse prices first
	prices := [][]float64{}
	doc.Find("table tbody td span").Each(func(i int, s *goquery.Selection) {
		if !strings.Contains(s.Text(), "|") || s.Parent().Is("span") {
			return
		}
		mealPrices := []float64{}
		for _, priceStr := range strings.Split(s.Text(), "|") {
			price, err := parsePrice(strings.TrimSpace(priceStr))
			if err != nil {
				return
			}
			mealPrices = append(mealPrices, price)
		}
		prices = append(prices, mealPrices)

	})

	// parse dates and meals
	datePattern := regexp.MustCompile("[0-9]{2}\\.[0-9]{2}\\.[0-9]{4}")
	var currentDate time.Time // keep track of the current date
	currentPriceIndex := 0    // keep track of which price should be added to which meal
	doc.Find("table tbody td").Each(func(i int, s *goquery.Selection) {
		dateStr := datePattern.FindString(strings.TrimSpace(s.Text()))
		if dateStr != "" {
			currentDate, err = time.Parse("02.01.2006", dateStr)
			return
		}

		s.Find("span").Each(func(j int, mealS *goquery.Selection) {
			// filter out nested <span> tags
			if mealS.Parent().Is("span") {
				return
			}
			// filter out prices
			if strings.Contains(mealS.Text(), "|") {
				return
			}
			// filter out empty slots
			name := strings.TrimSpace(mealS.Text())
			if slices.Contains([]string{"", "-", "Closed"}, name) {
				return
			}

			attributes, name := parseCKOAttributes(name)

			var meal Meal
			meal.NameEn = strings.Join(strings.Fields(name), " ")
			dateInt, err := strconv.Atoi(currentDate.Format("20060102"))
			if err != nil {
				return
			}
			meal.Id = dateInt*1000 + i*10 + j
			meal.Category.Id = 0
			meal.Date = currentDate
			meal.StudentPrice = prices[currentPriceIndex][0]
			meal.GuestPrice = prices[currentPriceIndex][1]
			meal.Location = ckoLocation
			for _, short := range attributes {
				attribute := GetAttribute(short)
				switch attribute.Type {
				case AllergenAttribute:
					meal.Allergens = append(meal.Allergens, attribute)
				case AdditiveAttribute:
					meal.Additives = append(meal.Additives, attribute)
				case FeatureAttribute:
					meal.Features = append(meal.Features, attribute)
				}
			}
			meal.Date = currentDate
			meals = append(meals, meal)

			currentPriceIndex++
		})
	})

	return meals, nil
}

func parseCKOAttributes(meal string) (attributes []string, name string) {
	name = strings.Join(strings.Fields(meal), " ")
	// parse vegan and veggie attributes separately, as they are longer than the others
	if strings.Contains(strings.ToLower(name), "vegan") {
		attributes = append(attributes, "Vegan")
		if strings.HasSuffix(strings.ToLower(name), "vegan") {
			name = strings.TrimSpace(strings.ReplaceAll(name, "vegan", ""))
		}
	}
	if strings.Contains(strings.ToLower(name), "veggie") {
		attributes = append(attributes, "Veggie")
		if strings.HasSuffix(strings.ToLower(name), "veggie") {
			name = strings.TrimSpace(strings.ReplaceAll(name, "veggie", ""))
		}
	}
	// regex to handle the inconsistent attribute formatting
	attributePattern := regexp.MustCompile(" (\\(?([A-Z]|[0-9]{1,2})\\)?(,?|, ))*$")
	matches := attributePattern.FindAllString(name, -1)
	var attributeStr string
	if len(matches) > 0 {
		attributeStr = matches[len(matches)-1]
		attributeStr = strings.ReplaceAll(attributeStr, " ", "")
		attributeStr = strings.ReplaceAll(attributeStr, "(", "")
		attributeStr = strings.ReplaceAll(attributeStr, ")", "")
	}
	// remove attributes from meal name
	name = strings.Replace(name, attributeStr, "", 1)
	attributes = append(attributes, strings.Split(strings.TrimSpace(attributeStr), ",")...)
	return attributes, strings.TrimSpace(name)
}

func parsePrice(priceStr string) (price float64, err error) {
	priceStrF := strings.ReplaceAll(priceStr, "€", "")
	priceStrF = strings.ReplaceAll(priceStrF, ",", ".")
	priceStrF = strings.TrimSpace(priceStrF)
	price, err = strconv.ParseFloat(priceStrF, 64)
	return price, err
}

func GetAttribute(short string) (attribute MealAttribute) {
	return GetCKOAttributes()[short]
}

// hardcoded list of all attributes from Campus Kitchen One, as they are not available in
// machine-readable format - this should be improved in the future
func GetCKOAttributes() (attributes map[string]MealAttribute) {
	attributes = map[string]MealAttribute{}

	attributes["1"] = MealAttribute{Id: 10001, Short: "1", NameDe: "geschwärzt", NameEn: "blackened", Type: AdditiveAttribute, Location: &ckoLocation}
	attributes["2"] = MealAttribute{Id: 10002, Short: "2", NameDe: "geschwefelt", NameEn: "sulfites", Type: AdditiveAttribute, Location: &ckoLocation}
	attributes["3"] = MealAttribute{Id: 10003, Short: "3", NameDe: "gewachst", NameEn: "waxed", Type: AdditiveAttribute, Location: &ckoLocation}
	attributes["4"] = MealAttribute{Id: 10004, Short: "4", NameDe: "Ei", NameEn: "chicken egg", Type: AllergenAttribute, Location: &ckoLocation}
	attributes["5"] = MealAttribute{Id: 10005, Short: "5", NameDe: "Soja", NameEn: "soy", Type: AllergenAttribute, Location: &ckoLocation}
	attributes["6"] = MealAttribute{Id: 10006, Short: "6", NameDe: "Milcherzeugnis", NameEn: "dairy", Type: AllergenAttribute, Location: &ckoLocation}
	attributes["7"] = MealAttribute{Id: 10007, Short: "7", NameDe: "Gluten", NameEn: "gluten", Type: AllergenAttribute, Location: &ckoLocation}
	attributes["8"] = MealAttribute{Id: 10008, Short: "8", NameDe: "Farbstoff", NameEn: "colouring", Type: AdditiveAttribute, Location: &ckoLocation}
	attributes["10"] = MealAttribute{Id: 10009, Short: "10", NameDe: "Geschmacksverstärker", NameEn: "flavour enhancers", Type: AdditiveAttribute, Location: &ckoLocation}
	attributes["14"] = MealAttribute{Id: 10010, Short: "14", NameDe: "chininhaltig", NameEn: "quinine", Type: AdditiveAttribute, Location: &ckoLocation}
	attributes["15"] = MealAttribute{Id: 10011, Short: "15", NameDe: "koffeinhaltig", NameEn: "caffeinated", Type: FeatureAttribute, Location: &ckoLocation}
	attributes["16"] = MealAttribute{Id: 10012, Short: "16", NameDe: "Alkohol", NameEn: "alcohol", Type: FeatureAttribute, Location: &ckoLocation}
	attributes["A"] = MealAttribute{Id: 10013, Short: "A", NameDe: "Schalenfrüchte", NameEn: "nuts", Type: AllergenAttribute, Location: &ckoLocation}
	attributes["B"] = MealAttribute{Id: 10014, Short: "B", NameDe: "Sellerie", NameEn: "celery", Type: AllergenAttribute, Location: &ckoLocation}
	attributes["C"] = MealAttribute{Id: 10015, Short: "C", NameDe: "Senf", NameEn: "mustard", Type: AllergenAttribute, Location: &ckoLocation}
	attributes["D"] = MealAttribute{Id: 10016, Short: "D", NameDe: "Sesam", NameEn: "sesame", Type: AllergenAttribute, Location: &ckoLocation}
	attributes["N"] = MealAttribute{Id: 10017, Short: "N", NameDe: "Nüsse", NameEn: "nuts", Type: AllergenAttribute, Location: &ckoLocation}
	attributes["F"] = MealAttribute{Id: 10018, Short: "F", NameDe: "Fisch", NameEn: "fish", Type: FeatureAttribute, Location: &ckoLocation}
	attributes["G"] = MealAttribute{Id: 10019, Short: "G", NameDe: "Geflügelfleisch", NameEn: "poultry", Type: FeatureAttribute, Location: &ckoLocation}
	attributes["K"] = MealAttribute{Id: 10020, Short: "K", NameDe: "Krebstiere", NameEn: "shellfish", Type: FeatureAttribute, Location: &ckoLocation}
	attributes["S"] = MealAttribute{Id: 10021, Short: "S", NameDe: "Schweinefleisch", NameEn: "pork", Type: FeatureAttribute, Location: &ckoLocation}
	attributes["P"] = MealAttribute{Id: 10022, Short: "P", NameDe: "Nitrit", NameEn: "nitrite", Type: AdditiveAttribute, Location: &ckoLocation}
	attributes["Vegan"] = MealAttribute{Id: 10023, Short: "Vegan", NameDe: "vegan", NameEn: "vegan", Type: FeatureAttribute, Location: &ckoLocation}
	attributes["Veggie"] = MealAttribute{Id: 10024, Short: "Veggie", NameDe: "vegetarisch", NameEn: "vegetarian", Type: FeatureAttribute, Location: &ckoLocation}

	return attributes
}

package http

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	. "github.com/slh335/hpi-mensa-api/types"
)

func (server *Server) parseLocation(c echo.Context, locationStr string) (ok bool, err error, location Location) {
	if strings.TrimSpace(locationStr) == "" {
		return false, c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "path param 'location' is required",
		}), Location{}
	}
	locations, err := server.LocationService.Get()
	if err != nil {
		return false, c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "failed to load locations",
		}), Location{}
	}

	for _, currentLocation := range locations {
		if strings.ReplaceAll(strings.ReplaceAll(strings.ToLower(locationStr), " ", ""), "-", "") ==
			strings.ReplaceAll(strings.ToLower(currentLocation.Name), " ", "") {

			location = currentLocation
			break
		}
	}
	if location.Id == 0 {
		return false, c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: fmt.Sprintf("invalid location: %s", locationStr),
		}), Location{}
	}
	return true, nil, location
}

func (server *Server) parseDate(c echo.Context, dateStr string) (ok bool, err error, date time.Time) {
	if strings.TrimSpace(dateStr) == "" {
		date, _ = time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
		return true, nil, date
	}
	date, err = time.Parse("2006-01-02", dateStr)
	if err != nil {
		return false, c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: fmt.Sprintf("invalid date: %s", dateStr),
		}), time.Time{}
	}
	return true, nil, date
}

func (server *Server) parseLanguage(c echo.Context, langStr string) (ok bool, err error, lang Language) {
	if strings.TrimSpace(langStr) == "" {
		return true, nil, English
	}
	switch strings.TrimSpace(langStr) {
	case "de":
		return true, nil, German
	case "en":
		return true, nil, English
	}

	return false, c.JSON(http.StatusInternalServerError, Response{
		Success: false,
		Message: fmt.Sprintf("invalid language: %s", langStr),
	}), English
}

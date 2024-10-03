package http

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/slh335/hpi-mensa-api/components"
	. "github.com/slh335/hpi-mensa-api/types"
	"github.com/slh335/hpi-mensa-api/util"
)

func (server *Server) GetMeals(c echo.Context) error {
	type Params struct {
		Location string `param:"location"`
		Date     string `query:"date"`
		Lang     string `query:"lang"`
		Format   string `query:"format"`
	}
	var params Params
	if err := c.Bind(&params); err != nil {
		return err
	}
	ok, err, location := server.parseLocation(c, params.Location)
	if !ok {
		return err
	}
	ok, err, date := server.parseDate(c, params.Date)
	if !ok {
		return err
	}
	ok, err, lang := server.parseLanguage(c, params.Lang)
	if !ok {
		return err
	}
	ok, err, format := server.parseFormat(c, params.Format)
	if !ok {
		return err
	}

	meals, err := server.MealService.Get(location, lang, date)
	if err != nil {
		fmt.Printf("Error: %+v\n", err)
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "failed to load meals",
		})
	}
	switch format {
	case FormatJSON:
		return c.JSON(http.StatusOK, Response{Success: true, Data: meals})
	case FormatHTML:
		return c.HTML(http.StatusOK, util.RenderComponent(components.Meals(meals)))
	}
	return nil
}

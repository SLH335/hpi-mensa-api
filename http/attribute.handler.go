package http

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	. "github.com/slh335/hpi-mensa-api/types"
)

func (server *Server) GetAdditives(c echo.Context) error {
	return server.getAttributes(c, AdditiveAttribute)
}

func (server *Server) GetAllergens(c echo.Context) error {
	return server.getAttributes(c, AllergenAttribute)
}

func (server *Server) GetFeatures(c echo.Context) error {
	return server.getAttributes(c, FeatureAttribute)
}

func (server *Server) getAttributes(c echo.Context, attributeType MealAttributeType) error {
	type Params struct {
		Location string `param:"location"`
		Lang     string `query:"lang"`
	}
	var params Params
	if err := c.Bind(&params); err != nil {
		return err
	}

	ok, err, location := server.parseLocation(c, params.Location)
	if !ok {
		return err
	}

	ok, err, lang := server.parseLanguage(c, params.Lang)
	if !ok {
		return err
	}

	attributes, err := server.MealService.AttributeService.Get(attributeType, location, lang)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: fmt.Sprintf("failed to load %ss", attributeType),
		})
	}

	return c.JSON(http.StatusOK, Response{Success: true, Data: attributes})
}

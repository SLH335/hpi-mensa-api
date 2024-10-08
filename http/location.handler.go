package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/slh335/hpi-mensa-api/components"
	. "github.com/slh335/hpi-mensa-api/types"
	"github.com/slh335/hpi-mensa-api/util"
)

func (server *Server) GetLocations(c echo.Context) error {
	type Params struct {
		Format string `query:"format"`
	}
	var params Params
	if err := c.Bind(&params); err != nil {
		return err
	}
	ok, err, format := server.parseFormat(c, params.Format)
	if !ok {
		return err
	}

	locations, err := server.LocationService.Get()
	if err != nil {
		fmt.Printf("Error: %+v\n", err)
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "failed to load locations",
		})
	}

	location := Location{Id: 9601, Name: "Griebnitzsee"}

	switch format {
	case FormatJSON:
		return c.JSON(http.StatusOK, Response{Success: true, Data: locations})
	case FormatHTML:
		return c.HTML(http.StatusOK, util.RenderComponent(components.LocationSelector(location, locations, Language{Id: 1, Short: "de"}, time.Now())))
	}
	return nil
}

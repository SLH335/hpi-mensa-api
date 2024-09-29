package http

import (
	"net/http"

	"github.com/labstack/echo/v4"

	. "github.com/slh335/hpi-mensa-api/types"
)

func (server *Server) GetLocations(c echo.Context) error {
	locations, err := server.LocationService.Get()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "error: failed to load locations",
		})
	}

	return c.JSON(http.StatusOK, Response{Success: true, Data: locations})
}

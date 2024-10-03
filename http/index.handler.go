package http

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/slh335/hpi-mensa-api/components"
	"github.com/slh335/hpi-mensa-api/util"
)

func (server *Server) Index(c echo.Context) error {
	return c.HTML(http.StatusOK, util.RenderComponent(components.Index()))
}

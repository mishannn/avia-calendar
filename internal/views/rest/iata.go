package rest

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type FindIATAParams struct {
	Filter string `query:"filter" validate:"required"`
}

func (s *Server) findIATA(c echo.Context) error {
	var params FindIATAParams
	if err := c.Bind(&params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("can't bind params: %w", err))
	}

	if err := c.Validate(&params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid params: %w", err))
	}

	items, err := s.ticketsService.SearchIATA(params.Filter)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("can't search IATA: %w", err))
	}

	return c.JSON(http.StatusOK, items)
}

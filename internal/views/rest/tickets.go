package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mishannn/avia-calendar/internal/views/rest/resttypes"
)

type FindCheapestTicketParams struct {
	From                string             `query:"from" validate:"required"`
	To                  string             `query:"to" validate:"required"`
	StartDate           resttypes.DateOnly `query:"start_date" validate:"required"`
	EndDate             resttypes.DateOnly `query:"end_date" validate:"required"`
	MaxTransfersCount   int                `query:"max_transfers_amount"`
	MaxTransferDuration int                `query:"max_transfer_duration"`
}

type FindCheapestTicketResponseEvent struct {
	Date            resttypes.DateOnly `json:"date"`
	TransfersAmount int                `json:"transfers_amount"`
	Price           float64            `json:"price"`
	SearchLink      string             `json:"search_link"`
	Error           string             `json:"error"`
}

// Cheapest Prices Page
func (s *Server) findCheapestTickets(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
	c.Response().Header().Set(echo.HeaderConnection, "keep-alive")
	c.Response().WriteHeader(http.StatusOK)

	var params FindCheapestTicketParams
	if err := c.Bind(&params); err != nil {
		writeRequestErrorEvent(c, fmt.Errorf("can't bind params: %w", err))
		return nil
	}

	if err := c.Validate(&params); err != nil {
		writeRequestErrorEvent(c, fmt.Errorf("invalid params: %w", err))
		return nil
	}

	startDate := time.Time(params.StartDate)
	endDate := time.Time(params.EndDate)

	if endDate.Before(startDate) {
		writeRequestErrorEvent(c, errors.New("end date shouldn't be before start date"))
		return nil
	}

	if endDate.Sub(startDate).Hours() > (30 * 24) {
		writeRequestErrorEvent(c, errors.New("max interval is 31 days"))
		return nil
	}

	if params.MaxTransfersCount < 0 {
		writeRequestErrorEvent(c, errors.New("transfers count can't be negative"))
		return nil
	}

	if params.MaxTransferDuration < 0 {
		writeRequestErrorEvent(c, errors.New("transfer duration can't be negative"))
		return nil
	}

	pricesCh, err := s.ticketsService.FindCheapestPrices(params.From, params.To, startDate, endDate, params.MaxTransfersCount, params.MaxTransferDuration)
	if err != nil {
		writeRequestErrorEvent(c, fmt.Errorf("can't start multisearch: %w", err))
		return nil
	}

	for price := range pricesCh {
		event := FindCheapestTicketResponseEvent{
			Date: resttypes.DateOnly(price.Date),
		}

		if price.Error != nil {
			event.Error = price.Error.Error()
		} else {
			event.TransfersAmount = price.TransfersAmount
			event.Price = price.Price
			event.SearchLink = price.SearchLink
		}

		eventJSON, err := json.Marshal(event)
		if err != nil {
			writeMarshalErrorEvent(c, fmt.Errorf("can't marshal event for date %s: %w", price.Date.Format("Mon, 02 Jan 2006"), err))
			continue
		}

		writeEvent(c, "price", eventJSON)
	}

	writeCloseEvent(c)
	return nil
}

func writeEvent(c echo.Context, name string, data []byte) {
	c.Response().Writer.Write([]byte(fmt.Sprintf("event: %s\n", name)))
	c.Response().Writer.Write([]byte(fmt.Sprintf("data: %s\n\n", data)))
	c.Response().Flush()
}

func writeRequestErrorEvent(c echo.Context, err error) {
	c.Response().Writer.Write([]byte("event: requesterror\n"))
	c.Response().Writer.Write([]byte(fmt.Sprintf("data: %s\n\n", err)))
	c.Response().Flush()
}

func writeMarshalErrorEvent(c echo.Context, err error) {
	c.Response().Writer.Write([]byte("event: marshalerror\n"))
	c.Response().Writer.Write([]byte(fmt.Sprintf("data: %s\n\n", err)))
	c.Response().Flush()
}
func writeCloseEvent(c echo.Context) {
	c.Response().Writer.Write([]byte("event: close\n"))
	c.Response().Writer.Write([]byte("data: [DONE]\n\n"))
	c.Response().Flush()
}

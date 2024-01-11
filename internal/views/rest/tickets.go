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

type FindCheapestTicketsParams struct {
	From                string             `query:"from" validate:"required"`
	To                  string             `query:"to" validate:"required"`
	StartDate           resttypes.DateOnly `query:"start_date" validate:"required"`
	EndDate             resttypes.DateOnly `query:"end_date" validate:"required"`
	MaxTransfersCount   int                `query:"max_transfers_amount"`
	MaxTransferDuration int                `query:"max_transfer_duration"`
}

type FindCheapestTicketsResponseFlight struct {
	From          string `json:"from"`
	DerartureTime string `json:"departure_time"`
	To            string `json:"to"`
	ArrivalTime   string `json:"arrival_time"`
	Airline       string `json:"airline"`
}

type FindCheapestTicketsResponseEvent struct {
	Date    resttypes.DateOnly `json:"date"`
	Tickets []resttypes.Ticket `json:"tickets"`
	Error   string             `json:"error"`
}

// Cheapest Prices Page
func (s *Server) findCheapestTickets(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
	c.Response().Header().Set(echo.HeaderConnection, "keep-alive")
	c.Response().WriteHeader(http.StatusOK)

	var params FindCheapestTicketsParams
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

	airlineNames := make(map[string]string)
	airlines, err := s.ticketsService.GetAirlines()
	if err != nil {
		writeRequestErrorEvent(c, fmt.Errorf("can't get airlines: %w", err))
		return nil
	}
	for _, airline := range airlines {
		if airline.Name != "" {
			airlineNames[airline.Code] = airline.Name
		} else if airline.NameTranslations["en"] != "" {
			airlineNames[airline.Code] = airline.NameTranslations["en"]
		}
	}

	ticketsEventsCh, err := s.ticketsService.FindTickets(params.From, params.To, startDate, endDate, params.MaxTransfersCount, params.MaxTransferDuration)
	if err != nil {
		writeRequestErrorEvent(c, fmt.Errorf("can't start multisearch: %w", err))
		return nil
	}

	for ticketsEvent := range ticketsEventsCh {
		event := FindCheapestTicketsResponseEvent{
			Date: resttypes.DateOnly(ticketsEvent.Date),
		}

		if ticketsEvent.Error != nil {
			event.Error = ticketsEvent.Error.Error()
		} else {
			event.Tickets = make([]resttypes.Ticket, 0, len(ticketsEvent.Tickets))
			for _, ticket := range ticketsEvent.Tickets {
				flights := make([]resttypes.Flight, 0, len(ticket.Flights))

				for _, flight := range ticket.Flights {
					var airlineName string
					if name, ok := airlineNames[flight.AirlineCode]; ok {
						airlineName = name
					} else {
						airlineName = flight.AirlineCode
					}

					flights = append(flights, resttypes.Flight{
						From:          flight.From,
						DerartureTime: flight.DerartureTime,
						To:            flight.To,
						ArrivalTime:   flight.ArrivalTime,
						AirlineName:   airlineName,
					})
				}

				event.Tickets = append(event.Tickets, resttypes.Ticket{
					SearchLink:      ticket.SearchLink,
					Date:            resttypes.DateOnly(ticket.Date),
					TransfersAmount: ticket.TransfersAmount,
					Price:           ticket.Price,
					Flights:         flights,
				})
			}
		}

		eventJSON, err := json.Marshal(event)
		if err != nil {
			writeMarshalErrorEvent(c, fmt.Errorf("can't marshal event for date %s: %w", ticketsEvent.Date.Format("Mon, 02 Jan 2006"), err))
			continue
		}

		writeEvent(c, "tickets", eventJSON)
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

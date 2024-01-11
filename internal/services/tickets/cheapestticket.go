package tickets

import (
	"errors"
	"fmt"
	"math"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mishannn/avia-calendar/internal/dal/aviasales"
)

type Flight struct {
	From          string
	DerartureTime string
	To            string
	ArrivalTime   string
	AirlineCode   string
}

type Ticket struct {
	SearchLink      string
	Date            time.Time
	TransfersAmount int
	Price           float64
	Flights         []Flight
}

func generateTicketSignature(ticket aviasales.Ticket, chunk aviasales.SearchResultsResponseChunk) string {
	airports := make([]string, 0)
	timestamps := make([]string, 0)

	var airlineCode string

	for legIndex, flightIndex := range ticket.Segments[0].Flights {
		flight := chunk.FlightLegs[flightIndex]

		if legIndex == 0 {
			airports = append(airports, flight.Origin)
			timestamps = append(timestamps, fmt.Sprint(flight.DepartureUnixTimestamp))
			airlineCode = flight.OperatingCarrierDesignator.AirlineID
		} else if airports[len(airports)-1] != flight.Origin {
			airports = append(airports, flight.Origin)
		}

		airports = append(airports, flight.Destination)
		if legIndex == len(ticket.Segments[0].Flights)-1 {
			timestamps = append(timestamps, fmt.Sprint(flight.ArrivalUnixTimestamp))
		}
	}

	return fmt.Sprintf("%s%s%s", airlineCode, strings.Join(timestamps, ""), strings.Join(airports, ""))
}

func getAgentLabel(proposal aviasales.TicketProposal, chunk aviasales.SearchResultsResponseChunk) string {
	agentIDStr := fmt.Sprintf("%d", proposal.AgentID)
	labels := chunk.Agents[agentIDStr].Label

	if label, ok := labels["ru"]; ok {
		return label.Default
	}

	if label, ok := labels["en"]; ok {
		return label.Default
	}

	return agentIDStr
}

func generateSearchLink(from string, to string, date time.Time, proposal aviasales.TicketProposal, ticket aviasales.Ticket, chunk aviasales.SearchResultsResponseChunk) string {
	urlValues := url.Values{}
	urlValues.Set("expected_price", fmt.Sprint(proposal.Price.Value))
	urlValues.Set("expected_price_currency", proposal.Price.CurrencyCode)
	urlValues.Set("expected_price_source", "share")
	urlValues.Set("expected_price_uuid", uuid.NewString())
	urlValues.Set("search_date", date.Format("02012006"))
	urlValues.Set("search_label", getAgentLabel(proposal, chunk))
	urlValues.Set("t", fmt.Sprintf("%s_%s_%.0f", generateTicketSignature(ticket, chunk), ticket.ID, proposal.Price.Value))

	return fmt.Sprintf("https://www.aviasales.ru/search/%s%s%s1?%s", from, date.Format("0201"), to, urlValues.Encode())
}

func getTransfersAmount(ticket aviasales.Ticket) int {
	return len(ticket.Segments[0].Transfers)
}

func getFlights(ticket aviasales.Ticket, chunk aviasales.SearchResultsResponseChunk) []Flight {
	flights := make([]Flight, 0, len(ticket.Segments[0].Flights))

	for _, flightIndex := range ticket.Segments[0].Flights {
		aviasalesFlight := chunk.FlightLegs[flightIndex]

		flights = append(flights, Flight{
			From:          aviasalesFlight.Origin,
			DerartureTime: aviasalesFlight.LocalDepartureDateTime,
			To:            aviasalesFlight.Destination,
			ArrivalTime:   aviasalesFlight.LocalArrivalDateTime,
			AirlineCode:   aviasalesFlight.OperatingCarrierDesignator.AirlineID,
		})
	}

	return flights
}

func createTicketFromAviasalesTicket(from string, to string, date time.Time, ticket aviasales.Ticket, chunk aviasales.SearchResultsResponseChunk) (*Ticket, error) {
	var cheapestProposal aviasales.TicketProposal

	cheapestPrice := math.MaxFloat64
	for _, proposal := range ticket.Proposals {
		if proposal.Price.Value < cheapestPrice {
			cheapestProposal = proposal
			cheapestPrice = proposal.Price.Value
		}
	}

	if cheapestPrice == math.MaxFloat64 {
		return nil, errors.New("can't find cheapest price")
	}

	return &Ticket{
		SearchLink:      generateSearchLink(from, to, date, cheapestProposal, ticket, chunk),
		Date:            date,
		TransfersAmount: getTransfersAmount(ticket),
		Price:           cheapestProposal.Price.Value,
		Flights:         getFlights(ticket, chunk),
	}, nil
}

func getTicketsFromChunk(from string, to string, date time.Time, chunk aviasales.SearchResultsResponseChunk) ([]Ticket, error) {
	tickets := make([]Ticket, 0, len(chunk.Tickets))

	for _, aviasalesTicket := range chunk.Tickets {
		ticket, err := createTicketFromAviasalesTicket(from, to, date, aviasalesTicket, chunk)
		if err != nil {
			return nil, fmt.Errorf("can't create cheapest ticket from aviasales ticket: %w", err)
		}

		tickets = append(tickets, *ticket)
	}

	return tickets, nil
}

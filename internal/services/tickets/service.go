package tickets

import (
	"fmt"
	"sync"
	"time"

	"github.com/mishannn/avia-calendar/internal/dal/aviasales"
	"github.com/patrickmn/go-cache"
)

type Service struct {
	cacheStorage    *cache.Cache
	aviasalesClient *aviasales.Client
}

func NewService(aviasalesClient *aviasales.Client) *Service {
	return &Service{
		cacheStorage:    cache.New(1*time.Hour, 1*time.Hour),
		aviasalesClient: aviasalesClient,
	}
}

func (s *Service) findTickets(from string, to string, date time.Time, maxTransfersCount int, maxTransferDuration int) ([]Ticket, error) {
	startSearchReqBody := aviasales.SearchStartRequestBody{
		SearchParams: aviasales.SearchParams{
			Directions: []aviasales.Direction{
				{
					Origin:      from,
					Destination: to,
					Date:        date.UTC().Format(time.DateOnly),
				},
			},
			Passengers: aviasales.Passengers{
				Adults:   1,
				Children: 0,
				Infants:  0,
			},
			TripClass: "Y",
		},
		MarketCode:   "ru",
		Marker:       "direct",
		Citizenship:  "RU",
		CurrencyCode: "rub",
	}

	startSearchRespBody, err := s.aviasalesClient.StartSearch(startSearchReqBody)
	if err != nil {
		return nil, fmt.Errorf("can't start search: %w", err)
	}

	transfersCountValues := make([]int, maxTransfersCount+1)
	for i := 0; i <= maxTransfersCount; i++ {
		transfersCountValues[i] = i
	}

	searchResultsReqBody := aviasales.SearchResultsRequestBody{
		SearchID: startSearchRespBody.SearchID,
		Order:    "cheapest",
		Limit:    3,
		Filters: map[string]any{
			"baggage":         []string{"full_baggage"},
			"transfers_count": transfersCountValues,
			"transfers_duration": []map[string]int{{
				"min": 0,
				"max": maxTransferDuration * 60,
			}},
			"transfers_without_airport_change": true,
			"transfers_without_visa":           true,
			"without_short_layover":            false,
		},
	}

	time.Sleep(3 * time.Second)

	var searchResultsRespBody aviasales.SearchResultsResponseBody
	for {
		var isLastResults bool

		searchResultsRespBody, isLastResults, err = s.aviasalesClient.GetSearchResults(startSearchRespBody.ResultsURL, searchResultsReqBody)
		if err != nil {
			return nil, fmt.Errorf("can't get search results: %w", err)
		}

		if isLastResults {
			break
		}

		time.Sleep(1 * time.Second)
	}

	return getTicketsFromChunk(from, to, date, searchResultsRespBody[0])
}

type TicketsEvent struct {
	Date    time.Time
	Tickets []Ticket
	Error   error
}

func (s *Service) FindTickets(from string, to string, startDate time.Time, endDate time.Time, maxTransfersCount int, maxTransferDuration int) (<-chan TicketsEvent, error) {
	cheapestPricesCh := make(chan TicketsEvent)

	go func() {
		fetcher := func(wg *sync.WaitGroup, date time.Time, output chan<- TicketsEvent) {
			defer wg.Done()

			cacheKey := fmt.Sprintf("%s:%s:%s:%d:%d", from, to, date.Format(time.DateOnly), maxTransfersCount, maxTransferDuration)

			if cachedTickets, ok := s.cacheStorage.Get(cacheKey); ok {
				output <- TicketsEvent{Date: date, Tickets: cachedTickets.([]Ticket)}
				return
			}

			tickets, err := s.findTickets(from, to, date, maxTransfersCount, maxTransferDuration)
			if err != nil {
				output <- TicketsEvent{Date: date, Error: err}
				return
			}

			s.cacheStorage.Set(cacheKey, tickets, 0)
			output <- TicketsEvent{
				Date:    date,
				Tickets: tickets,
			}
		}

		var wg sync.WaitGroup

		iterDate := startDate
		for iterDate.Before(endDate) || iterDate.Equal(endDate) {
			wg.Add(1)
			go fetcher(&wg, iterDate, cheapestPricesCh)
			iterDate = iterDate.AddDate(0, 0, 1)
		}

		wg.Wait()
		close(cheapestPricesCh)
	}()

	return cheapestPricesCh, nil
}

func (s *Service) SearchIATA(filter string) (aviasales.SearchIATAResponseBody, error) {
	return s.aviasalesClient.SearchIATA(filter)
}

func (s *Service) GetAirlines() (aviasales.GetAirlinesResponseBody, error) {
	return s.aviasalesClient.GetAirlines()
}

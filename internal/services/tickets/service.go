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

func (s *Service) findCheapestTicket(from string, to string, date time.Time, maxTransfersCount int, maxTransferDuration int) (*Ticket, error) {
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
		Limit:    1,
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

	return findCheapestTicketInChunk(from, to, date, searchResultsRespBody[0])
}

func (s *Service) FindCheapestPrices(from string, to string, startDate time.Time, endDate time.Time, maxTransfersCount int, maxTransferDuration int) (<-chan CheapestPriceForDate, error) {
	cheapestPricesCh := make(chan CheapestPriceForDate)

	go func() {
		fetcher := func(wg *sync.WaitGroup, date time.Time, output chan<- CheapestPriceForDate) {
			defer wg.Done()

			cacheKey := fmt.Sprintf("%s:%s:%s:%d:%d", from, to, date.Format(time.DateOnly), maxTransfersCount, maxTransferDuration)

			cachedItem, ok := s.cacheStorage.Get(cacheKey)
			if ok {
				output <- cachedItem.(CheapestPriceForDate)
				return
			}

			cheapestTicket, err := s.findCheapestTicket(from, to, date, maxTransfersCount, maxTransferDuration)
			if err != nil {
				output <- CheapestPriceForDate{Date: date, Error: err}
				return
			}

			result := CheapestPriceForDate{
				Date:            date,
				Price:           cheapestTicket.Price,
				TransfersAmount: cheapestTicket.TransfersAmount,
				SearchLink:      cheapestTicket.SearchLink,
			}

			s.cacheStorage.Set(cacheKey, result, 0)
			output <- result
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

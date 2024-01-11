package tickets

import (
	_ "embed"
	"encoding/json"
	"testing"
	"time"

	"github.com/mishannn/avia-calendar/internal/dal/aviasales"
	"github.com/stretchr/testify/assert"
)

//go:embed cheapestticket_test_response.json
var testRespBody []byte

func TestFindCheapestTicketInChunk(t *testing.T) {
	var response aviasales.SearchResultsResponseBody
	err := json.Unmarshal(testRespBody, &response)
	assert.Nil(t, err, "can't unmarshal test response body")

	date := time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC)

	cheapestTicket, err := getTicketsFromChunk("MOW", "HKT", date, response[0])
	assert.Nil(t, err, "can't find cheapest ticket in chunk")

	t.Logf("%+v", cheapestTicket)

	// expected := &Ticket{
	// 	SearchLink:      "https://www.aviasales.ru/search/MOW0801HKT1?expected_price=36721&expected_price_currency=RUB&expected_price_source=share&expected_price_uuid=4539cd1b-e2a5-4d8a-9eae-780a94bbd972&search_date=08012024&search_label=MEGO.travel&t=HU17075841001707656100SVOPEKHKT_ddf7090bf5831ac8266e1e7e5c28c37e_36721",
	// 	Date:            date,
	// 	TransfersAmount: 1,
	// 	Price:           36721,
	// }

	// // Replace generated UUID for test
	// replaceRe := regexp.MustCompile(`expected_price_uuid=[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)
	// cheapestTicket.SearchLink = replaceRe.ReplaceAllString(cheapestTicket.SearchLink, "expected_price_uuid=4539cd1b-e2a5-4d8a-9eae-780a94bbd972")

	// assert.Equal(t, cheapestTicket, expected, "unexpected result")
}

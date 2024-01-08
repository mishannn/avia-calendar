package tickets

import (
	_ "embed"
	"encoding/json"
	"testing"
	"time"

	"github.com/mishannn/avia-calendar/internal/dal/aviasales"
)

//go:embed cheapestticket_test_response.json
var testRespBody []byte

func TestFindCheapestTicketInChunk(t *testing.T) {
	var response aviasales.SearchResultsResponseBody
	err := json.Unmarshal(testRespBody, &response)
	if err != nil {
		t.Errorf("can't unmarshal test response body: %s", err)
		return
	}

	date := time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC)

	cheapestTicket, err := findCheapestTicketInChunk("MOW", "HKT", date, response[0])
	if err != nil {
		t.Errorf("can't find cheapest ticket in chunk: %s", err)
		return
	}

	t.Logf("%+v", cheapestTicket)
}

package resttypes

type Ticket struct {
	SearchLink      string   `json:"search_link"`
	Date            DateOnly `json:"date"`
	TransfersAmount int      `json:"transfers_amount"`
	Price           float64  `json:"price"`
	Flights         []Flight `json:"flights"`
}

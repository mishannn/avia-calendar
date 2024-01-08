package aviasales

// Search Start Request

type SearchStartRequestBody struct {
	SearchParams SearchParams `json:"search_params"`
	MarketCode   string       `json:"market_code"`
	Marker       string       `json:"marker"`
	Citizenship  string       `json:"citizenship"`
	CurrencyCode string       `json:"currency_code"`
}

type SearchParams struct {
	Directions []Direction `json:"directions"`
	Passengers Passengers  `json:"passengers"`
	TripClass  string      `json:"trip_class"`
}

type Passengers struct {
	Adults   int `json:"adults"`
	Children int `json:"children"`
	Infants  int `json:"infants"`
}

type Direction struct {
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
	Date        string `json:"date"`
}

// Search Start Response

type SearchStartResponseBody struct {
	SearchID          string `json:"search_id"`
	ResultsURL        string `json:"results_url"`
	PollingIntervalMS int    `json:"polling_interval_ms"`
}

// Search Results Request

type SearchResultsRequestBody struct {
	SearchID string         `json:"search_id"`
	Order    string         `json:"order"`
	Limit    int            `json:"limit"`
	Filters  map[string]any `json:"filters"`
}

type Filters struct {
	Baggage                       []string            `json:"baggage"`
	TransfersCount                []int               `json:"transfers_count"`
	TransfersDuration             []TransfersDuration `json:"transfers_duration"`
	TransfersWithoutAirportChange bool                `json:"transfers_without_airport_change"`
	TransfersWithoutVisa          bool                `json:"transfers_without_visa"`
	WithoutShortLayover           bool                `json:"without_short_layover"`
}

type TransfersDuration struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

// Search Results Response

type SearchResultsResponseBody []SearchResultsResponseChunk

type SearchResultsResponseChunk struct {
	Agents     map[string]Agent   `json:"agents"`
	Airlines   map[string]Airline `json:"airlines"`
	FlightLegs []FlightLeg        `json:"flight_legs"`
	Tickets    []Ticket           `json:"tickets"`
}

type Agent struct {
	ID             int             `json:"id"`
	GateName       string          `json:"gate_name"`
	Label          map[string]Name `json:"label"`
	PaymentMethods []string        `json:"payment_methods"`
	MobileVersion  bool            `json:"mobile_version"`
	HideProposals  bool            `json:"hide_proposals"`
	Assisted       bool            `json:"assisted"`
	AirlineIatas   []string        `json:"airline_iatas,omitempty"`
	MobileType     *string         `json:"mobile_type,omitempty"`
}

type Airline struct {
	Iata      string          `json:"iata"`
	IsLowcost bool            `json:"is_lowcost"`
	Name      map[string]Name `json:"name"`
	SiteName  string          `json:"site_name"`
}

type Name struct {
	Default string `json:"default"`
}

type FlightLeg struct {
	Origin                     string            `json:"origin"`
	Destination                string            `json:"destination"`
	LocalDepartureDateTime     string            `json:"local_departure_date_time"`
	LocalArrivalDateTime       string            `json:"local_arrival_date_time"`
	DepartureUnixTimestamp     int64             `json:"departure_unix_timestamp"`
	ArrivalUnixTimestamp       int64             `json:"arrival_unix_timestamp"`
	OperatingCarrierDesignator CarrierDesignator `json:"operating_carrier_designator"`
}

type CarrierDesignator struct {
	Carrier   string `json:"carrier"`
	AirlineID string `json:"airline_id"`
	Number    string `json:"number"`
}

type Ticket struct {
	ID        string           `json:"id"`
	Proposals []TicketProposal `json:"proposals"`
	Segments  []Segment        `json:"segments"`
}

type TicketProposal struct {
	ID    string `json:"id"`
	Price Price  `json:"price"`
	// PricePerPerson Price `json:"price_per_person"`
	AgentID int `json:"agent_id"`
}

type Price struct {
	CurrencyCode string  `json:"currency_code"`
	Value        float64 `json:"value"`
}

type Segment struct {
	Flights   []int      `json:"flights"`
	Transfers []Transfer `json:"transfers"`
}

type Transfer struct {
	VisaRules      VisaRules `json:"visa_rules"`
	RecheckBaggage bool      `json:"recheck_baggage"`
	NightTransfer  bool      `json:"night_transfer"`
}

type VisaRules struct {
	Required bool `json:"required"`
}

// Search IATA

type SearchIATAResponseBody []SearchResultsResponseBodyElement

type SearchResultsResponseBodyElement struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	CountryName string `json:"country_name"`
}

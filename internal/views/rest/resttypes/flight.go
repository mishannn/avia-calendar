package resttypes

type Flight struct {
	From          string `json:"from"`
	DerartureTime string `json:"departure_time"`
	To            string `json:"to"`
	ArrivalTime   string `json:"arrival_time"`
	AirlineName   string `json:"airline_name"`
}

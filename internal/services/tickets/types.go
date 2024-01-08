package tickets

import "time"

type CheapestPriceForDate struct {
	SearchLink      string
	Date            time.Time
	TransfersAmount int
	Price           float64
	Error           error
}

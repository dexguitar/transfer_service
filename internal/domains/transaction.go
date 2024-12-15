package domains

import "time"

type Transaction struct {
	ID        int
	Sender    string
	Receiver  string
	Amount    float64
	CreatedAt time.Time
}

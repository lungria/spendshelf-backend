package transaction

import "time"

type Transaction struct {
	ID          string
	Time        time.Time
	Description string
	MCC         int32
	Hold        bool
	Amount      int64
	AccountID   string
	CategoryID  int32
}

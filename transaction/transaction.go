package transaction

import "time"

// Transaction describes single user's transaction.
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

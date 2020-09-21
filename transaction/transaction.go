package transaction

import "time"

type Transaction struct {
	BankID      string
	Time        time.Time
	Description string
	MCC         int32
	Hold        bool
	Amount      int64
}

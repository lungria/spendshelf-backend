package privatbank

type Transaction struct {
	Description             string
	Amount                  float64
	Currency                string
	CardNumber              string
	BalanceAfterTransaction float64
}

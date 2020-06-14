package parser

import "errors"

type transaction struct {
	Description             string
	Amount                  float64
	Currency                string
	CardNumber              string
	BalanceAfterTransaction float64
}

func parseMessageBody(messageBody string) (*transaction, error) {
	if len(messageBody) == 0 {
		return nil, ErrBodyMustNotBeEmpty
	}
	return nil, nil
}

var ErrBodyMustNotBeEmpty = errors.New("message body must not be empty")

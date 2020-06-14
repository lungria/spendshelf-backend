package parser

import (
	"regexp"
	"strconv"

	"github.com/lungria/spendshelf-backend/src/privatbank"
)

const (
	amountSubexp                  = "amount"
	currencySubexp                = "currency"
	descriptionSubexp             = "description"
	cardNumberSubexp              = "cardNumber"
	balanceAfterTransactionSubexp = "balanceAfterTransaction"
)

type parser struct {
	regex *regexp.Regexp
}

func newParser() *parser {
	return &parser{
		regexp.MustCompile(`^(?P<amount>.+)(?P<currency>[A-Z]{3}) (?P<description>.+) (?P<cardNumber>\d+\*\d+) (?P<time>\d{2}:\d{2}).+ (?P<balanceAfterTransaction>\d+\.\d+)([A-Z]{3}).+$`),
	}
}

func (p *parser) parseMessageBody(messageBody string) (*privatbank.Transaction, error) {
	if len(messageBody) == 0 {
		return nil, newError(ErrEmptyMessageBody)
	}
	matches, ok := p.getParams(messageBody)
	if !ok {
		return nil, newError(ErrInvalidMessageBody)
	}

	amount, err := tryGetFloat64(matches, amountSubexp)
	if err != nil {
		return nil, newError(unableParseField(amountSubexp, err))
	}
	description, ok := matches[descriptionSubexp]
	if !ok {
		return nil, newError(unableParseField(descriptionSubexp, ErrFieldNotFound))
	}
	currency, ok := matches[currencySubexp]
	if !ok {
		return nil, newError(unableParseField(currencySubexp, ErrFieldNotFound))
	}
	cardNumber, ok := matches[cardNumberSubexp]
	if !ok {
		return nil, newError(unableParseField(cardNumberSubexp, ErrFieldNotFound))
	}
	balance, err := tryGetFloat64(matches, balanceAfterTransactionSubexp)
	if err != nil {
		return nil, newError(unableParseField(balanceAfterTransactionSubexp, err))
	}
	return &privatbank.Transaction{
		Description:             description,
		Amount:                  amount,
		Currency:                currency,
		CardNumber:              cardNumber,
		BalanceAfterTransaction: balance,
	}, nil
}

func (p *parser) getParams(messageBody string) (paramsMap map[string]string, ok bool) {
	match := p.regex.FindStringSubmatch(messageBody)
	if match == nil {
		return nil, false
	}

	paramsMap = make(map[string]string)
	for i, name := range p.regex.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	return paramsMap, true
}

func tryGetFloat64(matches map[string]string, fieldName string) (float64, error) {
	strAmount, ok := matches[fieldName]
	if !ok {
		return 0, ErrFieldNotFound
	}

	amount, err := strconv.ParseFloat(strAmount, 64)
	if err != nil {
		return 0, err
	}

	return amount, nil
}

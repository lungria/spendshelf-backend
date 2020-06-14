package parser

import (
	"errors"
	"testing"

	"github.com/lungria/spendshelf-backend/src/privatbank"
)

func Test_parseMessageBody_WhenMessageIsEmpty_ReturnsErrBodyMustNotBeEmpty(t *testing.T) {
	p := newParser()

	_, err := p.parseMessageBody("")

	if errors.Unwrap(err) != ErrEmptyMessageBody {
		t.Fatalf("expected error, received nil")
	}
}

func Test_parseMessageBody_WhenMessageIsValid_ReturnsTransaction(t *testing.T) {
	p := newParser()
	expected := privatbank.Transaction{
		Description: "Універмаг ПБ ЛП Покупка и доставка товаров",
		Amount:      843.59,
		Currency:    "UAH",
		CardNumber:  "5*45",
		//Time: time.Date(time.now),
		BalanceAfterTransaction: 9738.57,
	}
	transaction, err := p.parseMessageBody("843.59UAH Універмаг ПБ ЛП Покупка и доставка товаров 5*45 10:31 Бал. 9738.57UAH Кред. лiмiт 5000.0UAH")

	if err != nil {
		t.Fatalf("expected error: nil, received error: %s", err)
	}

	if *transaction != expected {
		t.Fatalf("expected transaction: %v, received transaction: %v", expected, transaction)
	}
}

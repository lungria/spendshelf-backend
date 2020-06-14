package parser

import (
	"testing"
)

func Test_parseMessageBody_WhenMessageIsEmpty_ReturnsErrBodyMustNotBeEmpty(t *testing.T) {
	_, err := parseMessageBody("")

	if err != ErrBodyMustNotBeEmpty {
		t.Fatalf("expected error: %s, received error: %s", ErrBodyMustNotBeEmpty, err)
	}
}

func Test_parseMessageBody_WhenMessageIsValid_ReturnsTransaction(t *testing.T) {
	expected := transaction{
		Description: "Універмаг ПБ ЛП Покупка и доставка товаров",
		Amount:      843.59,
		Currency:    "UAH",
		CardNumber:  "5*45",
		//Time: time.Date(time.now),
		BalanceAfterTransaction: 9738.57,
	}
	transaction, err := parseMessageBody("843.59UAH Універмаг ПБ ЛП Покупка и доставка товаров 5*45 10:31 Бал. 9738.57UAH Кред. лiмiт 5000.0UAH")

	if err != nil {
		t.Fatalf("expected error: nil, received error: %s", err)
	}

	if *transaction != expected {
		t.Fatalf("expected transaction: %v, received transaction: %v", expected, transaction)
	}
}

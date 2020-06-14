package parser

import (
	"errors"
	"fmt"
)

const (
	errorTemplate = "message body don't match expected format: %s"
)

var (
	ErrEmptyMessageBody   = errors.New("received empty message body")
	ErrInvalidMessageBody = errors.New("invalid message body format")
	ErrFieldNotFound      = errors.New("field not found")
)

type Error struct {
	message      string
	wrappedError error
}

func newError(reason error) error {
	return &Error{
		message:      fmt.Sprintf(errorTemplate, reason),
		wrappedError: reason,
	}
}

func (e *Error) Error() string {
	return e.message
}

func (e *Error) Unwrap() error {
	return e.wrappedError
}

func unableParseField(fieldName string, reason error) error {
	return fmt.Errorf("unable to parse %s. reason: %w", fieldName, reason)
}

package model

import (
	"fmt"
)

type ValidationErrorCode string

const (
	ErrInvalidValue ValidationErrorCode = "invalid_value"
	ErrMissingField ValidationErrorCode = "missing_field"

	ErrNotFound string = "not_found"
)

type ErrorOption func(*Error)

type ValidationError struct {
	Field string              `json:"field"`
	Code  ValidationErrorCode `json:"code"`
}

type Error struct {
	Code             string            `json:"code"`
	Message          string            `json:"message"`
	DetailedError    string            `json:"detailedError,omitempty"`
	ValidationErrors []ValidationError `json:"validationErrors,omitempty"`
}

func NewError(code string, message string, opts ...ErrorOption) *Error {
	e := Error{
		Code:    code,
		Message: message,
	}

	for _, opt := range opts {
		if opt == nil {
			continue
		}

		opt(&e)
	}

	return &e
}

func (e *Error) Error() string {
	msg := fmt.Sprintf("%s: %s", e.Code, e.Message)
	if e.DetailedError != "" {
		msg += fmt.Sprintf("(%s)", e.DetailedError)
	}
	return msg
}

func WithValidateError(field string, code ValidationErrorCode) ErrorOption {
	return func(err *Error) {
		err.ValidationErrors = append(err.ValidationErrors, ValidationError{field, code})
	}
}

func WithError(code string, err error) *Error {
	return &Error{
		Code:    code,
		Message: err.Error(),
	}
}

func NewNotFoundError(id, resource string) *Error {
	return NewError(ErrNotFound, fmt.Sprintf("%s %s not found", resource, id))
}

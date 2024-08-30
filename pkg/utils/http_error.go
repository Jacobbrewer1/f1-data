package utils

import (
	"errors"
)

type HttpError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewHttpError(code int, message string) *HttpError {
	return &HttpError{
		Code:    code,
		Message: message,
	}
}

func HttpErrorFromError(err error) *HttpError {
	httpErr := new(HttpError)
	ok := errors.As(err, &httpErr)
	if !ok {
		httpErr = NewHttpError(500, err.Error())
	}
	return httpErr
}

func (e *HttpError) Error() string {
	return e.Message
}

func (e *HttpError) Unwrap() error {
	return errors.New(e.Message)
}

func (e *HttpError) Is(target error) bool {
	t := new(HttpError)
	ok := errors.As(target, &t)
	if !ok {
		return false
	}
	return e.Code == t.Code
}

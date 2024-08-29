package http

import (
	"errors"
	"net/http"
)

func GenericErrorHandler(w http.ResponseWriter, _ *http.Request, err error) {
	SendErrorMessageWithStatus(w, http.StatusBadRequest, MsgBadRequest, err)
}

type httpError struct {
	// str is the string representation of the error.
	str string

	// status is the HTTP status code.
	status int
}

func (e *httpError) Error() string {
	return e.str
}

func (e *httpError) Status() int {
	return e.status
}

func NewHttpError(status int, str string) error {
	return &httpError{
		str:    str,
		status: status,
	}
}

func SendHttpError(w http.ResponseWriter, err error) {
	e := new(httpError)
	ok := errors.As(err, &e)
	if !ok {
		GenericErrorHandler(w, nil, err)
		return
	}

	SendMessageWithStatus(w, e.Status(), e.Error())
}

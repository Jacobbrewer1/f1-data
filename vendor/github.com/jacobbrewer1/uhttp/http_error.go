package uhttp

import (
	"net/http"

	"github.com/jacobbrewer1/uhttp/common"
)

type StatusCoder interface { // nolint:iface // This is an interface used by external libraries
	StatusCode() int
}

type HTTPError struct {
	error
	common.ErrorMessage
}

func (e *HTTPError) Error() string {
	return e.error.Error()
}

func (e *HTTPError) StatusCode() int {
	if e.Status == 0 {
		return http.StatusOK
	}
	return e.Status
}

func (e *HTTPError) SetRequestId(requestId string) {
	e.RequestId = requestId
}

// NewHTTPError creates a new HTTPError.
func NewHTTPError(code int, err error, details ...any) *HTTPError {
	errMsg := &common.ErrorMessage{
		Title:  http.StatusText(code),
		Detail: defaultHttpErrorDetail,
		Status: code,

		// RequestId will be populated at error write time
		RequestId: "",
	}

	if err != nil {
		errMsg.Detail = err.Error()
	}

	if len(details) > 0 {
		errMsg.Details = details
	}

	return &HTTPError{
		error:        err,
		ErrorMessage: *errMsg,
	}
}

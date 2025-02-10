package uhttp

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	errNotFound         = errors.New("not found")
	errMethodNotAllowed = errors.New("method not allowed")
	errUnauthorized     = errors.New("unauthorized")
)

// NotFoundHandler returns a handler that returns a 404 response.
func NotFoundHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw, ok := w.(*ResponseWriter)
		if !ok {
			rw = NewResponseWriter(w,
				WithDefaultStatusCode(http.StatusNotFound),
				WithDefaultHeader(HeaderRequestID, RequestIDFromContext(GenerateOrCopyRequestID(r.Context(), r))),
				WithDefaultHeader(HeaderContentType, ContentTypeJSON),
			)
		}

		details := []any{
			fmt.Sprintf("method: %s", r.Method),
			fmt.Sprintf("path: %s", r.URL.Path),
		}

		if r.URL.RawQuery != "" {
			details = append(details, fmt.Sprintf("query: %s", r.URL.RawQuery))
		}

		msg := NewHTTPError(http.StatusNotFound, errNotFound, details...)
		MustEncode(rw, http.StatusNotFound, msg)
	}
}

// MethodNotAllowedHandler returns a handler that returns a 405 response.
func MethodNotAllowedHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw, ok := w.(*ResponseWriter)
		if !ok {
			rw = NewResponseWriter(w,
				WithDefaultStatusCode(http.StatusMethodNotAllowed),
				WithDefaultHeader(HeaderRequestID, RequestIDFromContext(GenerateOrCopyRequestID(r.Context(), r))),
				WithDefaultHeader(HeaderContentType, ContentTypeJSON),
			)
		}

		details := []string{
			fmt.Sprintf("method: %s", r.Method),
			fmt.Sprintf("path: %s", r.URL.Path),
		}

		if r.URL.RawQuery != "" {
			details = append(details, fmt.Sprintf("query: %s", r.URL.RawQuery))
		}

		msg := NewHTTPError(http.StatusMethodNotAllowed, errMethodNotAllowed, details)
		MustEncode(rw, http.StatusMethodNotAllowed, msg)
	}
}

// UnauthorizedHandler returns a handler that returns a 401 response.
func UnauthorizedHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw, ok := w.(*ResponseWriter)
		if !ok {
			rw = NewResponseWriter(w,
				WithDefaultStatusCode(http.StatusUnauthorized),
				WithDefaultHeader(HeaderRequestID, RequestIDFromContext(GenerateOrCopyRequestID(r.Context(), r))),
				WithDefaultHeader(HeaderContentType, ContentTypeJSON),
			)
		}

		details := []string{
			fmt.Sprintf("method: %s", r.Method),
			fmt.Sprintf("path: %s", r.URL.Path),
		}

		if r.URL.RawQuery != "" {
			details = append(details, fmt.Sprintf("query: %s", r.URL.RawQuery))
		}

		msg := NewHTTPError(http.StatusUnauthorized, errUnauthorized, details)
		MustEncode(rw, http.StatusUnauthorized, msg)
	}
}

func GenericErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	rw, ok := w.(*ResponseWriter)
	if !ok {
		rw = NewResponseWriter(w,
			WithDefaultStatusCode(http.StatusBadRequest),
			WithDefaultHeader(HeaderRequestID, RequestIDFromContext(GenerateOrCopyRequestID(r.Context(), r))),
			WithDefaultHeader(HeaderContentType, ContentTypeJSON),
		)
	}

	msg := NewHTTPError(http.StatusBadRequest, err)
	MustEncode(rw, http.StatusBadRequest, msg)
}

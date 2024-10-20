package uhttp

import (
	"log/slog"
	"net/http"
)

// NotFoundHandler returns a handler that returns a 404 response.
func NotFoundHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw, ok := w.(*ResponseWriter)
		if !ok {
			rw = NewResponseWriter(w,
				WithDefaultStatusCode(http.StatusOK),
				WithDefaultHeader("X-Request-ID", RequestIDFromContext(GenerateOrCopyRequestID(r.Context(), r))),
				WithDefaultHeader(HeaderContentType, ContentTypeJSON),
			)
		}

		msg := NewMessage(MsgNotFound)
		err := EncodeJSON(rw, http.StatusNotFound, msg)
		if err != nil {
			slog.Error("Error encoding response", slog.String(loggingKeyError, err.Error()))
		}
	}
}

// MethodNotAllowedHandler returns a handler that returns a 405 response.
func MethodNotAllowedHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw, ok := w.(*ResponseWriter)
		if !ok {
			rw = NewResponseWriter(w,
				WithDefaultStatusCode(http.StatusOK),
				WithDefaultHeader("X-Request-ID", RequestIDFromContext(GenerateOrCopyRequestID(r.Context(), r))),
				WithDefaultHeader(HeaderContentType, ContentTypeJSON),
			)
		}

		msg := NewMessage(MsgMethodNotAllowed)
		err := EncodeJSON(rw, http.StatusMethodNotAllowed, msg)
		if err != nil {
			slog.Error("Error encoding response", slog.String(loggingKeyError, err.Error()))
		}
	}
}

// UnauthorizedHandler returns a handler that returns a 401 response.
func UnauthorizedHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw, ok := w.(*ResponseWriter)
		if !ok {
			rw = NewResponseWriter(w,
				WithDefaultStatusCode(http.StatusOK),
				WithDefaultHeader("X-Request-ID", RequestIDFromContext(GenerateOrCopyRequestID(r.Context(), r))),
				WithDefaultHeader(HeaderContentType, ContentTypeJSON),
			)
		}

		msg := NewMessage(MsgUnauthorized)
		err := EncodeJSON(rw, http.StatusUnauthorized, msg)
		if err != nil {
			slog.Error("Error encoding response", slog.String(loggingKeyError, err.Error()))
		}
	}
}

func GenericErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	rw, ok := w.(*ResponseWriter)
	if !ok {
		rw = NewResponseWriter(w,
			WithDefaultStatusCode(http.StatusOK),
			WithDefaultHeader("X-Request-ID", RequestIDFromContext(GenerateOrCopyRequestID(r.Context(), r))),
			WithDefaultHeader(HeaderContentType, ContentTypeJSON),
		)
	}

	msg := NewErrorMessage(MsgBadRequest, err)
	encErr := EncodeJSON(rw, http.StatusBadRequest, msg)
	if encErr != nil {
		slog.Error("Error encoding response", slog.String(loggingKeyError, encErr.Error()))
	}
}

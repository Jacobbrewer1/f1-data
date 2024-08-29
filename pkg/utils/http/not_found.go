package http

import (
	"log/slog"
	"net/http"

	"github.com/Jacobbrewer1/f1-data/pkg/logging"
)

// NotFoundHandler returns a handler that returns a 404 response.
func NotFoundHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg := NewMessage(MsgNotFound)
		err := Encode(w, http.StatusNotFound, msg)
		if err != nil {
			slog.Error("Error encoding response", slog.String(logging.KeyError, err.Error()))
		}
	}
}

// MethodNotAllowedHandler returns a handler that returns a 405 response.
func MethodNotAllowedHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg := NewMessage(MsgMethodNotAllowed)
		err := Encode(w, http.StatusMethodNotAllowed, msg)
		if err != nil {
			slog.Error("Error encoding response", slog.String(logging.KeyError, err.Error()))
		}
	}
}

// UnauthorizedHandler returns a handler that returns a 401 response.
func UnauthorizedHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg := NewMessage(MsgUnauthorized)
		err := Encode(w, http.StatusUnauthorized, msg)
		if err != nil {
			slog.Error("Error encoding response", slog.String(logging.KeyError, err.Error()))
		}
	}
}

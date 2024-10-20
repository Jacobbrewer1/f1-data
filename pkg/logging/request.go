package logging

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/jacobbrewer1/uhttp"
)

func LoggerFromRequest(r *http.Request) *slog.Logger {
	l := slog.Default()
	if r != nil {
		l = LoggerFromContext(r.Context())
	}
	return l
}

func LoggerFromContext(ctx context.Context) *slog.Logger {
	l := slog.Default()
	if ctx != nil {
		l = l.With(KeyRequestID, uhttp.RequestIDFromContext(ctx))
	}
	return l
}

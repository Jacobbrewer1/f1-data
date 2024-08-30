package http

import (
	"context"
	"net/http"
)

const (
	// requestIDHeader is the name of the HTTP header that contains a request ID added by nginx.
	requestIDHeader = "X-Request-ID"
)

var (
	// requestIDHeaderKey is the context key to the value of the request ID HTTP header.
	requestIDHeaderKey = ContextKey(requestIDHeader)
)

// RequestIDToContextMux returns a gorilla mux middleware which copies the request ID HTTP header into the provided
// context.
func RequestIDToContextMux() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			newCtx := RequestIDToContext(r.Context(), r)
			r = r.WithContext(newCtx)
			next.ServeHTTP(w, r)
		})
	}
}

// RequestIDToContext copies the request ID HTTP header into the provided context.
// This should be used as a param to the transport/http.ServerBefore() (go-kit) func.
func RequestIDToContext(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(ctx, requestIDHeaderKey, r.Header.Get(requestIDHeader))
}

// RequestIDFromContext returns the request ID HTTP header value from the provided context.
// If the header was not set it returns an empty string.
func RequestIDFromContext(ctx context.Context) string {
	v, ok := ctx.Value(requestIDHeaderKey).(string)
	if !ok {
		return ""
	}
	return v
}

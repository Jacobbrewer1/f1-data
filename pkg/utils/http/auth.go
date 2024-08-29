package http

import (
	"context"
	"net/http"
)

const (
	// authHeader is the Authorization header in a HTTP request.
	authHeader = "Authorization"
)

var (
	// authHeaderKey is the context key to the value of the Authorization HTTP request.
	authHeaderKey = ContextKey(authHeader)
)

// AuthHeaderToContext copies the Authorization HTTP header into the provided context.
func AuthHeaderToContext(ctx context.Context, r *http.Request) context.Context {
	return AuthToContext(ctx, r.Header.Get(authHeader))
}

// AuthHeaderToContextMux returns a gorilla mux middleware which copies the Authorization HTTP header into the
// provided context.
func AuthHeaderToContextMux() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			newCtx := AuthHeaderToContext(r.Context(), r)
			r = r.WithContext(newCtx)
			next.ServeHTTP(w, r)
		})
	}
}

// AuthToContext puts the authorisation into the context directly.
func AuthToContext(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, authHeaderKey, value)
}

// AuthHeaderFromContext returns the Authorization HTTP header value from the provided context.
// If the header was not set it returns an empty string.
func AuthHeaderFromContext(ctx context.Context) string {
	v, ok := ctx.Value(authHeaderKey).(string)
	if !ok {
		return ""
	}
	return v
}

// IsInternal returns true if the request is internal.
func IsInternal(r *http.Request) bool {
	return IsProxied(r)
}

// IsProxied returns true if the request is from kubernetes.
func IsProxied(r *http.Request) bool {
	// Kubernetes sets the X-Forwarded-For header when coming from the ingress, therefore we can check if the header is
	// set to determine if the request is from kubernetes.
	//
	// See: https://stackoverflow.com/questions/70164677/determine-if-http-request-to-a-service-is-within-or-outside-of-the-kubernetes-cl
	// The header will only be set from external traffic, so we can check if the header is set to determine if the
	// request is internal.

	h := r.Header.Get("X-Forwarded-For")
	// If the header is not set, then the request is internal.
	return h == ""
}

func InternalOnly(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !IsInternal(r) {
			SendMessageWithStatus(w, http.StatusForbidden, "Forbidden")
			return
		}
		next.ServeHTTP(w, r)
	}
}

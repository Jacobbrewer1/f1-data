package http

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsInternal_internal(t *testing.T) {
	internalRequest := httptest.NewRequest("GET", "/", nil)
	internalRequest.Header.Set("X-Your-Internal-Header", "true")

	got := IsInternal(internalRequest)
	require.True(t, got)
}

func TestIsInternal_external(t *testing.T) {
	externalRequest := httptest.NewRequest("GET", "/", nil)
	externalRequest.Header.Set("X-Forwarded-For", "test")

	got := IsInternal(externalRequest)
	require.False(t, got)
}

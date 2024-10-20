package uhttp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Encode defaults to JSON encoding
func Encode[T any](w http.ResponseWriter, status int, v T) error {
	return EncodeJSON(w, status, v)
}

// EncodeJSON encodes a response as JSON
func EncodeJSON[T any](w http.ResponseWriter, status int, v T) error {
	w.Header().Set(ContentTypeJSON, ContentTypeJSON)
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

func DecodeRequestJSON[T any](r *http.Request, v *T) error {
	return DecodeJSON[T](r.Body, v)
}

func DecodeJSON[T any](reader io.ReadCloser, v *T) error {
	if err := json.NewDecoder(reader).Decode(&v); err != nil {
		return fmt.Errorf("decode json: %w", err)
	}
	return nil
}

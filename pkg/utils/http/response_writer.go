package http

import (
	"net/http"
	"time"
)

// ClientWriter is a wrapper around http.ResponseWriter that provides additional information.
// It is used to capture the status code and the number of bytes written. It also provides a way to determine if the
// header has been written.
type ClientWriter struct {
	http.ResponseWriter
	statusCode      int
	bytesWritten    uint64
	isHeaderWritten bool
	startTime       time.Time
}

// NewClientWriter creates a new ClientWriter.
func NewClientWriter(w http.ResponseWriter) *ClientWriter {
	return &ClientWriter{
		ResponseWriter: w,
		startTime:      time.Now().UTC(),
	}
}

// Write writes the data to the connection as part of an HTTP reply.
func (c *ClientWriter) Write(p []byte) (bytes int, err error) {
	if !c.isHeaderWritten {
		c.SetJsonContentType()
		c.WriteHeader(http.StatusOK)
	}
	bytes, err = c.ResponseWriter.Write(p)
	c.bytesWritten += uint64(bytes)
	return
}

// WriteHeader sends an HTTP response header with the provided status code.
func (c *ClientWriter) WriteHeader(code int) {
	if c.isHeaderWritten {
		return
	}

	if c.Header().Get("Content-Type") == "" {
		c.SetJsonContentType()
	}

	c.ResponseWriter.WriteHeader(code)
	c.isHeaderWritten = true
	c.statusCode = code
}

// StatusCode returns the status code.
func (c *ClientWriter) StatusCode() int {
	if !c.isHeaderWritten || c.statusCode == 0 {
		// If the header has not been written or the status code has not been set, assume 200.
		return http.StatusOK
	}
	return c.statusCode
}

// BytesWritten returns the number of bytes written.
func (c *ClientWriter) BytesWritten() uint64 {
	return c.bytesWritten
}

// IsHeaderWritten returns true if the header has been written.
func (c *ClientWriter) IsHeaderWritten() bool {
	return c.isHeaderWritten
}

// SetStatus sets the status code.
func (c *ClientWriter) SetStatus(status int) {
	c.WriteHeader(status)
}

func (c *ClientWriter) SetContentType(contentType string) {
	c.Header().Set("Content-Type", contentType)
}

// SetJsonContentType sets the content type to JSON.
func (c *ClientWriter) SetJsonContentType() {
	c.SetContentType("application/json")
}

// SetXmlContentType sets the content type to XML.
func (c *ClientWriter) SetXmlContentType() {
	c.SetContentType("application/xml")
}

// GetRequestDuration gets the duration of the request
func (c *ClientWriter) GetRequestDuration() time.Duration {
	return time.Since(c.startTime)
}

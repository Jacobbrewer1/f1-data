package uhttp

type WriterOpt func(w *ResponseWriter)

// WithDefaultContentType sets the default content type for the response writer.
func WithDefaultContentType(contentType string) WriterOpt {
	return func(w *ResponseWriter) {
		if w.defaultHeaders == nil {
			w.defaultHeaders = make(map[string]string)
		}
		w.defaultHeaders[HeaderContentType] = contentType
	}
}

// WithDefaultStatusCode sets the default status code for the response writer.
func WithDefaultStatusCode(statusCode int) WriterOpt {
	return func(w *ResponseWriter) {
		w.defaultStatusCode = statusCode
	}
}

// WithDefaultHeader sets the default headers for the response writer.
func WithDefaultHeader(header, value string) WriterOpt {
	return func(w *ResponseWriter) {
		if w.defaultHeaders == nil {
			w.defaultHeaders = make(map[string]string)
		}
		w.defaultHeaders[header] = value
	}
}

func WithDefaultHeaders(headers map[string]string) WriterOpt {
	return func(w *ResponseWriter) {
		for k, v := range headers {
			WithDefaultHeader(k, v)(w) // This forces append to the defaultHeaders map with existing logic, rather than overwriting it.
		}
	}
}

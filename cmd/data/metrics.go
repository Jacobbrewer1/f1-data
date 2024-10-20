package main

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jacobbrewer1/f1-data/pkg/logging"
	"github.com/jacobbrewer1/uhttp"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// httpTotalRequests is the total number of http requests.
	httpTotalRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name:      "http_requests_total",
			Namespace: appName,
			Help:      "Total number of http requests",
		},
		[]string{"path", "method", "status_code"},
	)

	// httpRequestDuration is the duration of the http request.
	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:      "http_request_duration_seconds",
			Namespace: appName,
			Help:      "Duration of the http request",
		},
		[]string{"path", "method", "status_code"},
	)

	// httpRequestSize is the size of the http request.
	httpRequestSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:      "http_request_size",
			Namespace: appName,
			Help:      "Size of the http request",
		},
		[]string{"path", "method", "status_code"},
	)
)

// metricsMiddleware is run after the request is completed
func metricsMiddleware(w http.ResponseWriter, r *http.Request) {
	cw, ok := w.(*uhttp.ResponseWriter)
	if !ok {
		cw = uhttp.NewResponseWriter(w)
	}

	path := ""
	route := mux.CurrentRoute(r)
	if route != nil { // The route may be nil if the request is not routed.
		var err error
		path, err = route.GetPathTemplate()
		if err != nil {
			// An error here is only returned if the route does not define a path.
			slog.Error("Error getting path template", slog.String(logging.KeyError, err.Error()))
			path = r.URL.Path // If the route does not define a path, use the URL path.
		}
	} else {
		path = r.URL.Path // If the route is nil, use the URL path.
	}

	// Record the total number of requests.
	httpTotalRequests.WithLabelValues(path, r.Method, strconv.Itoa(cw.StatusCode())).Inc()

	// Record the request duration.
	httpRequestDuration.WithLabelValues(path, r.Method, strconv.Itoa(cw.StatusCode())).Observe(cw.GetRequestDuration().Seconds())

	// Record the request size.
	httpRequestSize.WithLabelValues(path, r.Method, strconv.Itoa(cw.StatusCode())).Observe(float64(r.ContentLength))
}

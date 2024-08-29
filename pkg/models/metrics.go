package models

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// DatabaseLatency is the duration of database queries.
var DatabaseLatency = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "database_latency",
		Help: "Duration of database queries",
	},
	[]string{"query"},
)

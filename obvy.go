/*

	Observability

*/

package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

// define prometheus metrics
var CFetchCount = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "CFetch_total",
	Help: "Total number of requests for DateTime.",
})

var pingCount = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "ping_total",
	Help: "Total number of requests for Readiness /ping.",
})

var apiDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
	Name: "CFetch_api_timer_seconds",
	Help: "Historgram for the total runtime of returning CFetch",
	// 50 Buckets, 10ms each, starting at 1ms
	Buckets: prometheus.LinearBuckets(0.001, 0.01, 50),
})

var dtCount = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "craque_dt_total",
	Help: "Total number of requests for DateTime.",
})

var CFetchDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
	Name: "CFetch_timer_seconds",
	Help: "Histogram for the runtime of each remote call to /fetch",
	// 50 Buckets, 10ms each, starting at 1ms
	Buckets: prometheus.LinearBuckets(0.001, 0.01, 50),
})

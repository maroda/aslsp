/*

	Main

*/

package main

import (
	"flag"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

func main() {
	// Prometheus outputs
	prometheus.MustRegister(CFetchCount)
	prometheus.MustRegister(pingCount)
	prometheus.MustRegister(apiDuration)
	prometheus.MustRegister(dtCount)
	prometheus.MustRegister(CFetchDuration)

	// Zerolog
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Runtime Flags
	debug := flag.Bool("debug", false, "Log Level: DEBUG")
	nofetch := flag.Bool("nofetch", false, "Data Backend, not a fetching job.")

	// Parse Flags
	flag.Parse()

	// Flag Options
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Info().Msg("Log level set to DEBUG")
	}

	// Fetching from a data backend is default behavior, displaying it to the user on port 8888.
	// The 'nofetch' flag turns this off, and this becomes the data backend to fetch from on port 9999.
	if *nofetch {
		log.Info().Msg("I am Bacque")
		// fetch ::: retrieive data for remote call
		http.HandleFunc("/fetch", CFetch)

		// ping ::: readiness check that returns 'pong'
		http.HandleFunc("/ping", ping)

		// metrics ::: prometheus metrics endpoint
		http.Handle("/metrics", promhttp.Handler())

		// start server
		if err := http.ListenAndServe(":9999", nil); err != nil {
			log.Fatal().Err(err).Msg("startup failed!")
		}
	} else {
		log.Info().Msg("I am Craque")
		// dt ::: gets the datetime from a remote service
		http.HandleFunc("/dt", dt)

		// ping ::: readiness check that returns 'pong'
		http.HandleFunc("/ping", ping)

		// metrics ::: prometheus metrics endpoint
		http.Handle("/metrics", promhttp.Handler())

		// start server
		if err := http.ListenAndServe(":8888", nil); err != nil {
			log.Fatal().Err(err).Msg("startup failed!")
		}
	}
}

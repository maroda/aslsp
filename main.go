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
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func main() {
	// datadog tracer is recommended to be placed around your main function.
	// These env settings are *maybe* overwritten by what's in the template yaml.
	// Doing a test to find out... these defaults will be overwritten by real values.
	tracer.Start(
		tracer.WithEnv("proto"),            // DD_ENV
		tracer.WithService("aslsp"),        // DD_SERVICE
		tracer.WithServiceVersion("0.0.0"), // DD_VERSION
	)
	defer tracer.Stop() // When stopped, the tracer flushes contents to the Agent.

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

	// ping ::: readiness check that returns 'pong'
	http.HandleFunc("/ping", ping)

	// metrics ::: prometheus metrics endpoint
	http.Handle("/metrics", promhttp.Handler())

	// Fetching from a data backend is default behavior, displaying it to the user on port 8888.
	// The 'nofetch' flag turns this off, and this becomes the data backend to fetch from on port 9999.
	//
	if *nofetch {
		log.Info().Msg("I am Bacque")

		// fetch ::: returns enriched date-time data to the client.
		http.HandleFunc("/fetch", CFetch)

		if err := http.ListenAndServe(":9999", nil); err != nil {
			log.Fatal().Err(err).Msg("startup failed!")
		}
	} else {
		log.Info().Msg("I am Craque")

		// dt ::: returns a date-time of some form to the user.
		http.HandleFunc("/dt", dt)

		if err := http.ListenAndServe(":8888", nil); err != nil {
			log.Fatal().Err(err).Msg("startup failed!")
		}
	}
}

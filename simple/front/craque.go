/*

	Craque

	/ - Hello
	/dt - retrieves the datetime from Bacque
	/ping - a readiness check
	/metrics - prometheus metrics

	Version = Cv009

*/

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// define prometheus metrics
var dtCount = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "craque_dt_total",
	Help: "Total number of requests for DateTime.",
})
var rootCount = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "craque_root_total",
	Help: "Total number of requests for /.",
})
var pingCount = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "craque_ping_total",
	Help: "Total number of requests for Readiness /ping.",
})
var fetchDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
	Name: "craque_fetch_timer_seconds",
	Help: "Histogram for the runtime of each remote call to /fetch",
	// 50 Buckets, 10ms each, starting at 1ms
	Buckets: prometheus.LinearBuckets(0.001, 0.01, 50),
})

// readiness check
func ping(w http.ResponseWriter, r *http.Request) {
	pingCount.Add(1)
	w.Write([]byte("pong\n"))
	zerolog.TimeFieldFormat = ""
	log.Info().
		Str("host", r.Host).
		Str("ref", r.RemoteAddr).
		Str("xref", r.Header.Get("X-Forwarded-For")).
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Str("proto", r.Proto).
		Str("agent", r.Header.Get("User-Agent")).
		Msg("")
}

// check the backend server for a datetimestamp
func dt(w http.ResponseWriter, r *http.Request) {
	// prometheus metrics
	dtCount.Add(1)
	dtTimer := prometheus.NewTimer(fetchDuration)
	defer dtTimer.ObserveDuration()

	// create client and run
	url := os.Getenv("BACQUE")
	craqueClient := http.Client{
		Timeout: time.Second * 2,
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal()
	}
	req.Header.Set("User-Agent", "craquego")

	res, getErr := craqueClient.Do(req)
	if getErr != nil {
		log.Fatal()
	}
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal()
	}
	fmt.Fprintf(w, "%s", body)

	zerolog.TimeFieldFormat = ""
	// log frontend access
	log.Info().
		Str("host", r.Host).
		Str("ref", r.RemoteAddr).
		Str("xref", r.Header.Get("X-Forwarded-For")).
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Str("proto", r.Proto).
		Str("agent", r.Header.Get("User-Agent")).
		Msg("")
	// log access to backend
	/*		currently disabled
	log.Info().
		Str("url", url).
		Str("status", res.Status).
		Msg("")
	*/
}

func main() {
	// Prometheus outputs
	prometheus.MustRegister(dtCount)
	prometheus.MustRegister(rootCount)
	prometheus.MustRegister(pingCount)
	prometheus.MustRegister(fetchDuration)

	// print Hello and the request path
	// if there is no valid endpoint, it will always default here
	/*
			disabled, not yet used

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			rootCount.Add(1)
			fmt.Fprintf(w, "Hello. %s\n", r.URL.Path)
			zerolog.TimeFieldFormat = ""
			log.Info().
				Str("host", r.Host).
				Str("ref", r.RemoteAddr).
				Str("xref", r.Header.Get("X-Forwarded-For")).
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Str("proto", r.Proto).
				Str("agent", r.Header.Get("User-Agent")).
				Msg("")
		})
	*/

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

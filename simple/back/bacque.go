/*

	Bacque

	/ - Hello
	/fetch - retrieves the localtime from the containing OS
	/ping - a liveliness check

	Version = Bv004

*/

package main

import (
	"fmt"
	"net/http"
	"os/exec"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// define prometheus metrics
var fetchCount = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "bacque_fetch_total",
	Help: "Total number of requests for DateTime.",
})
var rootCount = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "bacque_root_total",
	Help: "Total number of requests for /.",
})
var pingCount = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "bacque_ping_total",
	Help: "Total number of requests for Readiness /ping.",
})

// API call that returns local system datetime
func fetch(w http.ResponseWriter, r *http.Request) {
	// access a local command and return its output
	fetchCount.Add(1)
	app := "date"
	cmd := exec.Command(app)
	stdout, err := cmd.Output()
	if err != nil {
		log.Fatal()
	}
	fmt.Fprintf(w, "%s", stdout)
	zerolog.TimeFieldFormat = ""
	log.Info().
		Str("host", r.Host).
		Str("ref", r.RemoteAddr).
		Str("xref", r.Header.Get("X-Forwarded-For")).
		Str("path", r.URL.Path).
		Str("proto", r.Proto).
		Str("agent", r.Header.Get("User-Agent")).
		Msg("")
}

// liveliness check
func ping(w http.ResponseWriter, r *http.Request) {
	pingCount.Add(1)
	w.Write([]byte("pong\n"))
	zerolog.TimeFieldFormat = ""
	log.Info().
		Str("host", r.Host).
		Str("ref", r.RemoteAddr).
		Str("xref", r.Header.Get("X-Forwarded-For")).
		Str("path", r.URL.Path).
		Str("proto", r.Proto).
		Str("agent", r.Header.Get("User-Agent")).
		Msg("")
}

func main() {
	prometheus.MustRegister(fetchCount)
	prometheus.MustRegister(rootCount)
	prometheus.MustRegister(pingCount)

	// print Hello and the request path
	// if there is no valid endpoint, it will always default here
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rootCount.Add(1)
		fmt.Fprintf(w, "Hello. %s\n", r.URL.Path)
		zerolog.TimeFieldFormat = ""
		log.Info().
			Str("host", r.Host).
			Str("ref", r.RemoteAddr).
			Str("xref", r.Header.Get("X-Forwarded-For")).
			Str("path", r.URL.Path).
			Str("proto", r.Proto).
			Str("agent", r.Header.Get("User-Agent")).
			Msg("")
	})

	// fetch local command output
	http.HandleFunc("/fetch", fetch)

	// ping ::: liveliness check that returns 'pong'
	http.HandleFunc("/ping", ping)

	// m ::: prometheus metrics endpoint (can this be logged?)
	http.Handle("/m", promhttp.Handler())

	// start server
	if err := http.ListenAndServe(":9999", nil); err != nil {
		log.Fatal().Err(err).Msg("startup failed!")
	}
}

/*
	Craque Lib
*/

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// define prometheus metrics
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

// check the backend server for a datetimestamp
func dt(w http.ResponseWriter, r *http.Request) {
	// prometheus metrics
	dtCount.Add(1)
	dtTimer := prometheus.NewTimer(CFetchDuration)
	defer dtTimer.ObserveDuration()

	zerolog.TimeFieldFormat = ""

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

	res, err := craqueClient.Do(req)
	if err != nil {
		// if BACQUE fails, fall back to a local timestamp only,
		// not the enriched IP addresses that BACQUE returns.
		http.Error(w, http.StatusText(418), 418)

		arg := "+%Y%m%d%H%S"
		app := "date"

		dtloc, err := exec.Command(app, arg).Output()
		if err != nil {
			log.Fatal()
		}
		fmt.Fprintf(w, "DateTime=%s", dtloc)

		// log service unresponsive
		log.Error().
			Str("host", r.Host).
			Str("ref", r.RemoteAddr).
			Str("xref", r.Header.Get("X-Forwarded-For")).
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("proto", r.Proto).
			Str("agent", r.Header.Get("User-Agent")).
			Str("response", "418").
			Err(err).
			Msg("service unresponsive local dt returned")
		return
	}
	// don't waste resources
	defer res.Body.Close()

	// TODO: status code should probably be checked to move on
	// fmt.Printf("StatusCode: %d: %q", res.StatusCode, res.Request.URL) // debug

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal()
	}
	fmt.Fprintf(w, "%s", body)

	// log frontend access
	log.Info().
		Str("host", r.Host).
		Str("ref", r.RemoteAddr).
		Str("xref", r.Header.Get("X-Forwarded-For")).
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Str("proto", r.Proto).
		Str("agent", r.Header.Get("User-Agent")).
		Str("response", "200").
		Msg("")
}

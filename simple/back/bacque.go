/*

	Bacque

	/ - Hello
	/fetch - returns three 'dynamic' actions:
			- retrieves local timestamp from the container OS
			- displays the client Request IP address
			- reports Local IP based on default egress
	/ping - a readiness check
	/metrics - prometheus metrics

	Version = Bv008

*/

package main

import (
	"fmt"
	"net"
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
var apiDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
	Name: "bacque_api_timer_seconds",
	Help: "Historgram for the total runtime of returning /fetch",
	// 50 Buckets, 10ms each, starting at 1ms
	Buckets: prometheus.LinearBuckets(0.001, 0.01, 50),
})

// API call that returns local system datetime
func fetch(w http.ResponseWriter, r *http.Request) {
	// access a local command and return its output
	fetchCount.Add(1)
	dtTimer := prometheus.NewTimer(apiDuration)
	defer dtTimer.ObserveDuration()

	arg := "+%Y%m%d%H%S"
	app := "date"
	stdout, err := exec.Command(app, arg).Output()
	if err != nil {
		log.Fatal()
	}

	fmt.Fprintf(w, "DateTime=%s", stdout)

	// grab just the IP of the requestor
	rHost, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Fatal()
	}

	userIP := net.ParseIP(rHost)
	if userIP == nil {
		log.Fatal()
	}

	// display request IP
	fmt.Fprintf(w, "RequestIP=%s\n", rHost)

	// an outgoing UDP connection reveals the egress IP
	extAddr := "8.8.8.8:80"
	conn, err := net.Dial("udp", extAddr)
	if err != nil {
		log.Fatal()
	}
	defer conn.Close()

	// grab just the local IP (byte -> string conversion required)
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	lHost, _, err := net.SplitHostPort(localAddr.String())
	if err != nil {
		log.Fatal()
	}

	// display local IP
	fmt.Fprintf(w, "LocalIP=%s\n", lHost)

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

// readiness check
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
	prometheus.MustRegister(apiDuration)

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

	// ping ::: readiness check that returns 'pong'
	http.HandleFunc("/ping", ping)

	// metrics ::: prometheus metrics endpoint
	http.Handle("/metrics", promhttp.Handler())

	// start server
	if err := http.ListenAndServe(":9999", nil); err != nil {
		log.Fatal().Err(err).Msg("startup failed!")
	}
}

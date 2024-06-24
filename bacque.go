/*

	Bacque

	/fetch - returns three 'dynamic' actions
			 that rely on several I/O subsystems:
			- retrieves local timestamp from the container OS
			- displays the client Request IP address
			- reports Local IP based on default egress
	/ping - a readiness check
	/metrics - prometheus metrics

*/

package main

import (
	"fmt"
	"net"
	"net/http"
	"os/exec"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// CFetch ::: API call that returns local system datetime
func CFetch(w http.ResponseWriter, r *http.Request) {
	// prometheus tracing
	CFetchCount.Add(1)
	dtTimer := prometheus.NewTimer(apiDuration)
	defer dtTimer.ObserveDuration()

	// access a local command and return its output
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
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Str("proto", r.Proto).
		Str("agent", r.Header.Get("User-Agent")).
		Str("response", "200").
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
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Str("proto", r.Proto).
		Str("agent", r.Header.Get("User-Agent")).
		Str("response", "200").
		Msg("")
}

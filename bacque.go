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
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"os/exec"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// LocalCMD. Takes an output destination, the OS app to use, and the app args
// The reason for this is to guarantee an interaction with the operating system
// as a "dependency" for when this service is not available.
func LocalCMD(w io.Writer, app, arg string) error {
	cmdOut, err := exec.Command(app, arg).Output()
	if err != nil {
		log.Fatal()
		return err
	}

	fmt.Fprintf(w, "%s", bytes.Trim(cmdOut, "\n\r"))
	return err
}

// IPFinder. Determines the local machine IP address.
type IPFinder interface {
	EgressIP() (string, error)
}

// IPConfig. Contains IP addresses used for gathering data.
type IPConfig struct {
	ExtIPwPort string
}

// EgressIP. Uses a network call to extract the local IP address.
// An outgoing UDP connection reveals the egress IP.
func (ipc *IPConfig) EgressIP() (string, error) {
	conn, err := net.Dial("udp", ipc.ExtIPwPort)
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

	return lHost, err
}

// FindIP. Locate the local system's network egress IP address.
func FindIP(i IPFinder) (string, error) {
	// this will work putting in something mnaually
	extIP, err := i.EgressIP()
	if err != nil {
		log.Fatal()
	}
	return extIP, err
}

// EgIP. Takes an IP:PORT string and returns the local IP address
// used for outbound network traffic.
func EgIP(e string) (string, error) {
	// an outgoing UDP connection reveals the egress IP
	conn, err := net.Dial("udp", e)
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

	return lHost, err
}

// CFetch. API call that returns local system datetime
func CFetch(w http.ResponseWriter, r *http.Request) {
	// prometheus tracing
	CFetchCount.Add(1)
	dtTimer := prometheus.NewTimer(apiDuration)
	defer dtTimer.ObserveDuration()

	// access a local command and return its output
	app := "date"
	arg := "+%Y%m%d%H%S"
	lcB := bytes.Buffer{}
	lcerr := LocalCMD(&lcB, app, arg)
	if lcerr != nil {
		log.Fatal()
	}
	fmt.Fprintf(w, "DateTime=%q\n", lcB.String())

	// grab just the IP of the requestor
	//	TODO: Add to IPFinder
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

	// Determine the local IP address
	// ExtIPwPort can be anything reachable
	lHost, err := FindIP(&IPConfig{ExtIPwPort: "8.8.8.8:80"})
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

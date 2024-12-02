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

// LocalCMD takes an output destination, the OS app to use, and the app args
// The reason for this is to guarantee an interaction with the operating system
// as a "dependency" for when this service is not available.
func LocalCMD(w io.Writer, app, arg string) error {
	cmdOut, err := exec.Command(app, arg).Output()
	if err != nil {
		log.Error().Err(err).Msg("Could not execute command")
		return err
	}

	_, err = fmt.Fprintf(w, "%s", bytes.Trim(cmdOut, "\n\r"))
	if err != nil {
		log.Error().Err(err).Msg("Could not write output")
	}

	return err
}

// IPFinder determines the local machine IP address.
type IPFinder interface {
	EgressIP() (string, error)
}

// IPConfig contains IP addresses used for gathering data.
type IPConfig struct {
	ExtIPwPort string
}

// EgressIP uses a network call to extract the local IP address.
// An outgoing UDP connection reveals the egress IP.
func (ipc *IPConfig) EgressIP() (string, error) {
	conn, err := net.Dial("udp", ipc.ExtIPwPort)
	if err != nil {
		log.Error().Err(err).Msg("Could not connect to external IP")
	}
	defer conn.Close()

	// grab just the local IP (byte -> string conversion required)
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	lHost, _, err := net.SplitHostPort(localAddr.String())
	if err != nil {
		log.Error().Err(err).Msg("Could not split host and port")
	}

	return lHost, err
}

// FindIP to locate the local system's network egress IP address.
func FindIP(i IPFinder) (string, error) {
	// this will work putting in something manually
	extIP, err := i.EgressIP()
	if err != nil {
		log.Error().Err(err).Msg("Could not find egress IP")
	}
	return extIP, err
}

// EgIP takes an IP:PORT string and returns the local IP address
// used for outbound network traffic.
func EgIP(e string) (string, error) {
	var lHost string

	// validate /e/
	rHost, _, err := net.SplitHostPort(e)
	if err != nil {
		log.Error().Err(err).Msg("Could not split host and port")
		return rHost, err
	}

	// an outgoing UDP connection reveals the egress IP
	conn, err := net.Dial("udp", e)
	if err != nil {
		log.Error().Err(err).Msg("Could not connect to external IP")
	}
	defer conn.Close()

	// grab just the local IP (byte -> string conversion required)
	localAddr, ok := conn.LocalAddr().(*net.UDPAddr)
	if ok {
		lHost, _, err = net.SplitHostPort(localAddr.String())
		if err != nil {
			log.Error().Err(err).Msg("Could not split host and port")
		}
	}

	return lHost, err
}

// RequestIP takes the active connection and writes the source IP.
// This is a bit different from other functions,
// because it writes its result instead of returning one.
func RequestIP(w http.ResponseWriter, r *http.Request) error {
	// grab just the IP of the requester
	rHost, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Error().Err(err).Msg("Could not split host and port")
		return err
	}

	// display request IP
	_, err = fmt.Fprintf(w, "RequestIP=%s\n", rHost)
	if err != nil {
		log.Error().Err(err).Msg("Could not write output")
		return err
	}

	return nil
}

// CFetch is an API call that returns local system datetime
func CFetch(w http.ResponseWriter, r *http.Request) {
	var err error
	var lHost string

	// prometheus tracing
	CFetchCount.Add(1)
	dtTimer := prometheus.NewTimer(apiDuration)
	defer dtTimer.ObserveDuration()

	// access a local command and return its output
	app := "date"         // Shell command
	arg := "+%Y%m%d%H%S"  // Command args
	lcB := bytes.Buffer{} // LocalCMD buffer
	err = LocalCMD(&lcB, app, arg)
	if err != nil {
		log.Error().Err(err).Msg("Could not execute command")
	}

	// TODO: this should be done by LocalCMD?
	_, err = fmt.Fprintf(w, "DateTime=%q\n", lcB.String())
	if err != nil {
		log.Error().Err(err).Msg("Could not write output")
	}

	// access the request IP and return its output
	err = RequestIP(w, r)
	if err != nil {
		log.Error().Err(err).Msg("Could not fetch IP")
	}

	// Determine the local IP address
	// ExtIPwPort can be anything reachable
	lHost, err = FindIP(&IPConfig{ExtIPwPort: "8.8.8.8:80"})
	if err != nil {
		log.Error().Err(err).Msg("Could not find IP")
	}

	// TODO: this should be done by FindIP?
	// display local IP
	_, err = fmt.Fprintf(w, "LocalIP=%s\n", lHost)
	if err != nil {
		log.Error().Err(err).Msg("Could not read output")
	}

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

// ping is a simple readiness check that answers with 'pong'
func ping(w http.ResponseWriter, r *http.Request) {
	pingCount.Add(1)
	_, err := w.Write([]byte("pong"))
	if err != nil {
		log.Error().Err(err).Msg("Could not write output")
	}

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

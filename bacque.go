/*

	Bacque

	/fetch - returns three 'dynamic' actions
			 that rely on several I/O subsystems:
			- retrieves local timestamp from the container OS
			- displays the client Request IP address
			- reports Local IP based on default egress
	/ping - a readiness check

*/

package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"os/exec"

	"github.com/rs/zerolog/log"
)

/*
CFetch is the primary handler,
returning underlying system data that gets displayed to the client.

Individual OS tasks are spawned from here.
*/
func CFetch(w http.ResponseWriter, r *http.Request) {
	var err error
	var lHost string

	// Get a local datetime
	datetime, err := getDate(osDateCmd, osDateArg)
	if err != nil {
		log.Error().Err(err).Msg("Could not get datetime")
	}

	// Print datetime to the output
	_, err = fmt.Fprintf(w, "DateTime=%q\n", datetime)
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

// IPFinder determines IP addresses
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
	// For now all we need is to call EgIP
	return EgIP(ipc.ExtIPwPort)
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
// This is the non-interface version of EgressIP
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
		log.Error().Err(err).Msg("Could not create connection")
		return rHost, err
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Error().Err(err).Msg("Could not close external IP connection")
			return
		}
	}()

	// grab the local IP
	localAddr, ok := conn.LocalAddr().(*net.UDPAddr)
	if ok {
		// Only split off the port if we get an IPv4 address
		if ipv4 := localAddr.IP.To4(); ipv4 != nil {
			lHost, _, err = net.SplitHostPort(localAddr.String())
			if err != nil {
				log.Error().Err(err).Msg("Could not split host and port")
			}
		} else {
			lHost = localAddr.IP.String()
		}
	}

	return lHost, err
}

// RequestIP takes the active connection and writes the source IP.
// This is a bit different from other functions,
// because it writes its result instead of returning one.
func RequestIP(w http.ResponseWriter, r *http.Request) error {
	// grab the IP of the requester
	rHost, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Error().Err(err).Msg("Could not split host and port")
		return err
	}

	// display request IP
	if rHost == "::1" {
		rHost = "localhost-ipv6"
	}
	_, err = fmt.Fprintf(w, "RequestIP=%s\n", rHost)
	if err != nil {
		log.Error().Err(err).Msg("Could not write output")
		return err
	}

	return nil
}

// ping is a simple readiness check that answers with 'pong'
func ping(w http.ResponseWriter, r *http.Request) {
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

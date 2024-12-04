/*
	Craque

	/dt - attempts a call to bacque/fetch and prints the results

*/

package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	osDateCmd = "date"
	osDateArg = "+%Y%m%d%H%S"
)

func getDate(cmd, arg string) (string, error) {
	appbuf := bytes.Buffer{}
	err := LocalCMD(&appbuf, cmd, arg)
	if err != nil {
		log.Error().Err(err).Msg("Could not get date")
	}

	// t := time.Now().Format("200601021504") // Full layout date: "20060102150405"
	t := appbuf.String()
	return t, err
}

func failBacque(w io.Writer) error {
	// Get a local datetime
	datetime, err := getDate(osDateCmd, osDateArg)
	if err != nil {
		log.Error().Err(err).Msg("Could not get local datetime")
	}

	_, err = fmt.Fprintf(w, "DateTime=%q\n", datetime)
	if err != nil {
		log.Error().Err(err).Msg("Could not write DateTime output")
	}

	return err
}

/*
dt is the primary handler,
checking a remote host for endpoint data
and substituting a fall-back value when it's unavailable.
*/
func dt(w http.ResponseWriter, r *http.Request) {
	// this value expects the full url,
	// i.e.: export BACQUE="http://localhost:9999/fetch"
	url := os.Getenv("BACQUE")
	craqueClient := http.Client{Timeout: time.Second * 2}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Error().Err(err).Msg("Could not create new request")
	}
	req.Header.Set("User-Agent", "craquego")

	res, err := craqueClient.Do(req)
	if err != nil {
		http.Error(w, http.StatusText(418), 418)

		// if BACQUE fails, fall back to a local timestamp only,
		// not the enriched IP addresses that BACQUE returns.
		// failBacque()

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

	}
	defer func() {
		err := res.Body.Close()
		if err != nil {
			log.Error().Err(err).Msg("Could not close response body")
			return
		}
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error().Err(err).Msg("Could not read response body")
	}

	_, err = fmt.Fprintf(w, "%s", body)
	if err != nil {
		log.Error().Err(err).Msg("Could not write output")
	}

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

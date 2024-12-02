/*
	Craque
*/

package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// check the backend server for a datetimestamp
func dt(w http.ResponseWriter, r *http.Request) {
	// prometheus tracing
	dtCount.Add(1)
	dtTimer := prometheus.NewTimer(CFetchDuration)
	defer dtTimer.ObserveDuration()

	zerolog.TimeFieldFormat = ""

	// this value expects the full url,
	// i.e.: http://localhost:9999/fetch
	url := os.Getenv("BACQUE")
	craqueClient := http.Client{Timeout: time.Second * 2}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Error().Err(err).Msg("Could not create request")
	}
	req.Header.Set("User-Agent", "craquego")

	res, err := craqueClient.Do(req)
	if err != nil {
		// if BACQUE fails, fall back to a local timestamp only,
		// not the enriched IP addresses that BACQUE returns.
		http.Error(w, http.StatusText(418), 418)
		app := "date"
		arg := "+%Y%m%d%H%S"
		lcB := bytes.Buffer{}
		lcerr := LocalCMD(&lcB, app, arg)
		if lcerr != nil {
			log.Fatal()
		}
		_, err := fmt.Fprintf(w, "DateTime=%q\n", lcB.String())
		if err != nil {
			log.Error().Err(err).Msg("Could not write output")
		}

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
	defer res.Body.Close()

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

/*

	Craque

	/ - Hello
	/dt - retrieves the datetime from Bacque
	/ping - a liveliness check

*/

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// liveliness check
func ping(w http.ResponseWriter, r *http.Request) {
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

// check the backend server for a datetimestamp
// for now hardcoded, with kube it will need modification
func dt(w http.ResponseWriter, r *http.Request) {
	// url := "http://localhost:9999/fetch"
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
		Str("path", r.URL.Path).
		Str("proto", r.Proto).
		Str("agent", r.Header.Get("User-Agent")).
		Msg("")
	// log access to backend
	log.Info().
		Str("url", url).
		Str("status", res.Status).
		Msg("")
}

// print Hello and the request path
// if there is no valid endpoint, it will always default here
func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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

	// valid endpoints:
	// dt ::: gets the datetime from a remote service
	http.HandleFunc("/dt", dt)
	// ping ::: liveliness check that returns 'pong'
	http.HandleFunc("/ping", ping)

	// start server
	if err := http.ListenAndServe(":8888", nil); err != nil {
		log.Fatal().Err(err).Msg("startup failed!")
	}
}

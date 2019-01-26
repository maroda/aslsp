/*

	Bacque

	/ - Hello
	/fetch - retrieves the localtime from the containing OS
	/ping - a liveliness check

	Version = Bv003

*/

package main

import (
	"fmt"
	"net/http"
	"os/exec"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// API call that returns local system datetime
func fetch(w http.ResponseWriter, r *http.Request) {
	// access a local command and return its output
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

	// fetch local command output
	http.HandleFunc("/fetch", fetch)
	// ping ::: liveliness check that returns 'pong'
	http.HandleFunc("/ping", ping)

	// start server
	if err := http.ListenAndServe(":9999", nil); err != nil {
		log.Fatal().Err(err).Msg("startup failed!")
	}
}

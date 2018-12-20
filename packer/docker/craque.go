// Simple web thing for creating some JSON logs to be ingested by SumoLogic

package main

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// All this does is return a simple string when "/" is requested, and then logs it.
// There is a special version of zerolog called logh that interfaces with net/http,
// but the exercise here is just to get simple JSON log output so I'm not using it.

// Just a little alternate method for log variance
func admin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello. %s\n", r.URL.Path)
	zerolog.TimeFieldFormat = ""
	log.Info().
		Str("host", r.Host).
		Str("path", r.URL.Path).
		Str("proto", r.Proto).
		Msg("")
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
	zerolog.TimeFieldFormat = ""
	log.Info().
		Str("host", r.Host).
		Str("path", r.URL.Path).
		Str("proto", r.Proto).
		Msg("")
}

func main() {
	// just print Hello and the request path
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello. %s\n", r.URL.Path)
		//
		// This currently logs to StdOut (I think)
		// To get them to Sumo, will have to investigate further...
		zerolog.TimeFieldFormat = ""
		log.Info().
			Str("host", r.Host).
			Str("path", r.URL.Path).
			Str("proto", r.Proto).
			Msg("")
	})
	// other handlers for variety
	http.HandleFunc("/ping", ping)
	http.HandleFunc("/admin", admin)

	if err := http.ListenAndServe(":8888", nil); err != nil {
		log.Fatal().Err(err).Msg("startup failed!")
	}
}

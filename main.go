/*

	Main

*/

package main

import (
	"flag"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Zerolog
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Runtime Flags
	debug := flag.Bool("debug", false, "Log Level: DEBUG")
	nofetch := flag.Bool("nofetch", false, "Data Backend, not a fetching job.")

	// Simple Parse Flags
	flag.Parse()

	// Flag Options
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Info().Msg("Log level set to DEBUG")
	}

	// ping ::: readiness check that returns 200 OK 'pong'
	http.HandleFunc("/ping", ping)

	// Fetching from a data backend is default behavior, displaying it to the user on port 8888.
	// The 'nofetch' flag turns this off, and this becomes the data backend to fetch from on port 9999.
	if *nofetch {
		// fetch ::: returns enriched date-time data to the client.
		log.Info().Msg("I am Bacque, listening for /fetch on port 9999")
		http.HandleFunc("/fetch", CFetch)

		if err := http.ListenAndServe(":9999", nil); err != nil {
			log.Fatal().Err(err).Msg("startup failed!")
		}
	} else {
		// dt ::: returns a date-time of some form to the user.
		log.Info().Msg("I am Craque, listening for /dt on port 8888")
		http.HandleFunc("/dt", dt)

		if err := http.ListenAndServe(":8888", nil); err != nil {
			log.Fatal().Err(err).Msg("startup failed!")
		}
	}
}

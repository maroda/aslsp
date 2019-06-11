package main

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type intHandle struct{}

func (h intHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	l, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	for _, f := range l {
		fmt.Fprintf(w, "int: %s\n", f.Name)
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

// Spin up a consumer if feature is set
func consumeK() {

	featK := os.Getenv("BACQUE_KAFKA")

	if featK == "on" {

		broker := os.Getenv("BACQUE_KAFKA_BROKER")
		group := os.Getenv("BACQUE_KAFKA_GROUP")

		// this isn't working yet
		// it's accepting the slice,
		// but the consumer is getting empty events
		var topics []string
		for _, n := range os.Getenv("BACQUE_KAFKA_TOPIC") {
			topics = append(topics, string(n))
		}
		// the following shows that this slice is correct
		// so something is wrong with how it's being used
		// fmt.Printf("%#v\n", topics)
		// []string{"c", "r", "a", "q", "K", "T"}
		//
		// the fix here may be to run this with a commandline arg
		// and duplicate *exactly* how the barebones example works
		// maybe that will clue me into what's fucked

		c, err := kafka.NewConsumer(&kafka.ConfigMap{
			"bootstrap.servers":     broker,
			"broker.address.family": "v4", // avoids v6 broker resolution on OSX
			"group.id":              group,
			"session.timeout.ms":    6000,
			"auto.offset.reset":     "earliest"})
		if err != nil {
			log.Error().Err(err).Msg("Failed to create consumer!")
			return
		}

		log.Info().Msg("Created consumer")

		// need to look into this to see why it's puking on my slice
		// it expects a []string but i'm not seeing it working
		err = c.SubscribeTopics(topics, nil)

		run := true

		for run == true {
			ev := c.Poll(100)
			if ev == nil {
				// log.Debug().Msg("ev is nil")
				continue
			}

			// the confluent lib is receiving my slice
			// as if it were individual elements
			// instead of a concatenated slice

			switch e := ev.(type) {
			case *kafka.Message:
				log.Info().
					Str("Value", string(e.Value)).
					Msg("")
			case kafka.Error:
				log.Info().
					Str("error", string(e.Code())).
					Msg("Broker Error")
			default:
				log.Warn().Msg("Ignored event")
			}
		}

		log.Info().Msg("Closing consumer")
		c.Close()
	}
}

func main() {

	go consumeK()

	err := http.ListenAndServe(":7777", intHandle{})
	if err != nil {
		log.Fatal()
	}
}

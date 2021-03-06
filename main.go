package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	dbPath = "/tmp/miniqueue_db"

	tlsCertPath = "./testdata/localhost.pem"
	tlsKeyPath  = "./testdata/localhost-key.pem"
)

func main() {
	// If the binary is run with ENV=PRODUCTION, use JSON formatted logging.
	if os.Getenv("ENV") != "PRODUCTION" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	port := flag.String("port", "8080", "port used to run the server")
	flag.Parse()

	// Start the server
	srv := newServer(newBroker(newStore(dbPath)))
	p := fmt.Sprintf(":%s", *port)

	log.Info().Str("port", p).Msg("starting miniqueue")

	if err := http.ListenAndServeTLS(p, tlsCertPath, tlsKeyPath, srv); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal().Err(err).Msg("server closed")
	}
}

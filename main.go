package main

import (
	"html/template"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

const service = "rndir.cc"

var (
	logger zerolog.Logger

	tmpl = template.Must(template.ParseGlob("*.html"))
	port = "8080"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	logger = zerolog.New(os.Stderr).With().
		Caller().Timestamp().Logger()

	if serverPort := os.Getenv("SERVER_PORT"); serverPort != "" {
		port = serverPort
	}
}

func main() {
	logger.Info().Msgf("starting %s on port:%s", service, port)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Lookup("index.html").Execute(w, nil)
	})

	fileServer := http.FileServer(http.Dir("assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fileServer))

	srv := &http.Server{
		Addr:         net.JoinHostPort("", port),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		IdleTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		logger.Error().Err(err).Msgf("%s is down", service)
	}
}

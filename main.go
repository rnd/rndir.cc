package main

import (
	"context"
	_ "embed"
	"html/template"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

const service = "rndir.cc"

var (
	ver = "devel"

	logger zerolog.Logger

	port = "8080"

	tpl        *template.Template
	assetsRoot = "assets"

	//go:embed index.html
	indexHTML []byte
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	logger = zerolog.New(os.Stderr).With().
		Caller().Timestamp().Logger()

	if serverPort := os.Getenv("SERVER_PORT"); serverPort != "" {
		port = serverPort
	}
	tpl = template.Must(
		template.New("index.html").Parse(string(indexHTML)),
	)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tpl.ExecuteTemplate(w, "index.html", nil)
	})

	fileServer := http.FileServer(http.Dir(assetsRoot))
	mux.Handle("/a/", http.StripPrefix("/a/", fileServer))

	srv := &http.Server{
		Addr:         net.JoinHostPort("", port),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		IdleTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	idleConnsClosed := make(chan struct{})
	// watch for os.Interrupt signal and gracefully shutdown
	// the server.
	go func() {
		const shutdownTimeout = 10 * time.Second

		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		ctx, cancel := context.WithTimeout(
			context.Background(),
			shutdownTimeout,
		)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			logger.Error().Err(err).Msg("HTTP server shutdown")
		}
		logger.Info().Msg("HTTP server shutdown")
		close(idleConnsClosed)
	}()

	logger.Info().Msgf("starting %s on port:%s", service, port)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		logger.Error().Err(err).Msg("HTTP server ListenAndServe")
	}
	<-idleConnsClosed
}

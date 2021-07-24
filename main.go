package main

import (
	"html/template"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

const service = "rndir.cc"

var (
	logger zerolog.Logger

	port = "8080"

	tmpl       *template.Template
	assetsRoot = "assets"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	logger = zerolog.New(os.Stderr).With().
		Caller().Timestamp().Logger()

	if serverPort := os.Getenv("SERVER_PORT"); serverPort != "" {
		port = serverPort
	}

	tmplRoot := "templates"
	if dir := os.Getenv("KO_DATA_PATH"); dir != "" {
		assetsRoot = filepath.Join(dir, assetsRoot)
		tmplRoot = filepath.Join(dir, tmplRoot)
	}
	tmpl = template.Must(
		template.ParseGlob(filepath.Join(tmplRoot, "*.html")),
	)
}

func main() {
	logger.Info().Msgf("starting %s on port:%s", service, port)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Lookup("index.html").Execute(w, nil)
	})

	fileServer := http.FileServer(http.Dir(assetsRoot))
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

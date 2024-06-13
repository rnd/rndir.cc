package server

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/rnd/site/config"
	"github.com/rnd/site/logging"
	"github.com/rnd/site/web/home"
)

const (
	defaultIdleTimeout  = 30 * time.Second
	defaultReadTimeout  = 5 * time.Second
	defaultWriteTimeout = 10 * time.Second
)

type Server struct {
	*http.Server
}

func New(cfg config.ServerCfg) *Server {
	return &Server{
		&http.Server{
			Handler:      mux(),
			Addr:         net.JoinHostPort(cfg.Host, cfg.Port),
			ReadTimeout:  defaultIdleTimeout,
			WriteTimeout: defaultWriteTimeout,
			IdleTimeout:  defaultIdleTimeout,
		},
	}
}

func mux() *http.ServeMux {
	mux := http.NewServeMux()
	home.NewHandler(mux)
	return mux
}

// Run runs server processes if only all dependencies are resolved.
func (s *Server) Run(ctx context.Context) error {
	logging.FromContext(ctx).
		InfoContext(ctx, "starting server", "url", s.Addr)

	if err := s.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}
	return nil
}

// Shutdown gracefully shutdown server processes.
func (s *Server) Shutdown(ctx context.Context) error {
	logging.FromContext(ctx).
		InfoContext(ctx, "shutting down server")

	return s.Server.Shutdown(ctx)
}

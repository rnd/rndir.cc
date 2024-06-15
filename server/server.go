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

	cfg config.ServerCfg
}

func New(cfg config.ServerCfg) *Server {
	return &Server{
		Server: &http.Server{
			Addr:         net.JoinHostPort(cfg.Host, cfg.Port),
			ReadTimeout:  defaultIdleTimeout,
			WriteTimeout: defaultWriteTimeout,
			IdleTimeout:  defaultIdleTimeout,
		},
		cfg: cfg,
	}
}

func (s *Server) register() {
	mux := http.NewServeMux()
	home.Register(s.cfg, mux)
	s.Handler = mux
}

// Run runs server processes if only all dependencies are resolved.
func (s *Server) Run() error {
	logging.New().Info("starting server", "url", s.Addr)

	s.register()
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

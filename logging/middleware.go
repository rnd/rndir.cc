package logging

import (
	"log/slog"
	"net/http"

	"github.com/rnd/site/config"
)

func Middleware(cfg config.ServerCfg, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := New()
		ctx := WithContext(
			r.Context(), logger,
			slog.String("service", cfg.Name),
			slog.String("version", cfg.Version),
		)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

package home

import (
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/a-h/templ"

	"github.com/rnd/site/config"
	"github.com/rnd/site/logging"
)

var (
	once sync.Once

	path string
)

type Handler struct {
	mux *http.ServeMux
}

func Register(cfg config.ServerCfg, mux *http.ServeMux) {
	handler, pat := mux.Handler(&http.Request{URL: &url.URL{Path: "/"}})
	if _, ok := handler.(*Handler); ok && pat == "/" {
		return
	}
	h := &Handler{mux: mux}
	mux.Handle("/", logging.Middleware(cfg, h))
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logging.FromContext(ctx).DebugContext(ctx, "incoming request")
	path := r.URL.Path
	switch {
	case path == "/": // index
		templ.Handler(home(nixStorePath())).ServeHTTP(w, r)
	case strings.HasPrefix(r.URL.Path, "/static"): // static assets
		http.StripPrefix("/static", http.FileServer(http.Dir("./static"))).ServeHTTP(w, r)
	default: // 404
		errorHandler(w, r, http.StatusNotFound)
	}
}

func errorHandler(w http.ResponseWriter, r *http.Request, statusCode int) {
	ctx := r.Context()

	logging.FromContext(ctx).DebugContext(ctx, "incoming request", "path", r.URL.Path)

	switch statusCode {
	case http.StatusNotFound:
		templ.Handler(
			errorP(nixStorePath(), statusCode, "not found"),
			templ.WithStatus(statusCode),
		).ServeHTTP(w, r)
	}
}

func nixStorePath() string {
	once.Do(func() {
		path, _ = os.Executable()
	})
	return path
}

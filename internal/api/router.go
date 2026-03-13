package api

import (
	"io/fs"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ptmind/piadmin/internal/auth"
	"github.com/ptmind/piadmin/internal/config"
	"github.com/ptmind/piadmin/internal/monitor"
)

func NewRouter(cfg *config.Config, a *auth.Auth, collector *monitor.Collector, assets fs.FS) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))

	ah := &authHandler{auth: a}
	sh := &systemHandler{collector: collector}
	wh := &wsHandler{collector: collector}
	ph := &processHandler{}
	svh := &servicesHandler{}
	nh := &networkHandler{}
	th := &terminalHandler{}
	fh := &filesHandler{}
	gh := &gpioHandler{}

	// public
	r.Post("/api/auth/login", ah.Login)

	// protected
	r.Group(func(r chi.Router) {
		r.Use(a.Middleware)

		// System monitoring
		r.Get("/api/system/snapshot", sh.GetSnapshot)
		r.Get("/api/ws/monitor", wh.Monitor)

		// System management
		r.Get("/api/processes", ph.List)
		r.Delete("/api/processes", ph.Kill)
		r.Get("/api/services", svh.List)
		r.Get("/api/services/{name}", svh.Status)
		r.Post("/api/services/{name}", svh.Action)
		r.Get("/api/network/interfaces", nh.Interfaces)

		// Web terminal
		r.Get("/api/ws/terminal", th.Handle)

		// File manager
		r.Get("/api/files", fh.List)
		r.Get("/api/files/download", fh.Download)
		r.Post("/api/files/upload", fh.Upload)
		r.Delete("/api/files", fh.Delete)
		r.Post("/api/files/mkdir", fh.Mkdir)
		r.Put("/api/files/rename", fh.Rename)

		// GPIO
		r.Get("/api/gpio/available", gh.Available)
		r.Get("/api/gpio/pins", gh.List)
		r.Post("/api/gpio/pins", gh.SetPin)
		r.Post("/api/gpio/export", gh.Export)
		r.Post("/api/gpio/unexport", gh.Unexport)
	})

	// SPA static files
	if assets != nil {
		r.Handle("/*", spaHandler(assets))
	}

	return r
}

func spaHandler(assets fs.FS) http.HandlerFunc {
	fileServer := http.FileServer(http.FS(assets))

	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}

		f, err := assets.Open(path)
		if err != nil {
			r.URL.Path = "/index.html"
		} else {
			f.Close()
		}
		fileServer.ServeHTTP(w, r)
	}
}

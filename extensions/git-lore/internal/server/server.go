package server

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/joshmcarthur/git-lore/extensions/git-lore/internal/api"
	loregit "github.com/joshmcarthur/git-lore/extensions/git-lore/internal/git"
)

//go:embed all:webdist
var webDist embed.FS

// New creates the root HTTP handler: API + embedded static UI.
func New(repo *loregit.Repo) (http.Handler, error) {
	mux := http.NewServeMux()
	apiHandler := &api.Handler{Repo: repo}
	apiHandler.Mount(mux)

	sub, err := fs.Sub(webDist, "webdist")
	if err != nil {
		return nil, err
	}
	fileServer := http.FileServer(http.FS(sub))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}
		if _, err := fs.Stat(sub, path); err != nil {
			// SPA fallback
			r.URL.Path = "/"
			fileServer.ServeHTTP(w, r)
			return
		}
		fileServer.ServeHTTP(w, r)
	})

	return mux, nil
}

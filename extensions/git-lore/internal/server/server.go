package server

import (
	"embed"
	"io/fs"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/joshmcarthur/git-lore/extensions/git-lore/internal/api"
	loregit "github.com/joshmcarthur/git-lore/extensions/git-lore/internal/git"
)

//go:embed all:webdist
var webDist embed.FS

// Options configures the served UI.
type Options struct {
	// Editable enables document editing from the browser (serve --edit).
	Editable bool
	// Addr is the listen address; its host is accepted alongside loopback
	// names when guarding writes.
	Addr string
}

// New creates the root HTTP handler: API + embedded Vue SPA.
func New(repo *loregit.Repo, opts Options) (http.Handler, error) {
	mux := http.NewServeMux()
	apiHandler := &api.Handler{Repo: repo, Editable: opts.Editable}
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

	return guardWrites(mux, opts.Addr), nil
}

// guardWrites rejects state-changing requests whose Host or Origin is not
// this server, so other sites open in the browser (including via DNS
// rebinding) cannot modify the repository.
func guardWrites(next http.Handler, addr string) http.Handler {
	allowed := map[string]bool{"localhost": true, "127.0.0.1": true, "::1": true}
	if host, _, err := net.SplitHostPort(addr); err == nil && host != "" {
		allowed[host] = true
	}
	hostOK := func(hostport string) bool {
		host := hostport
		if h, _, err := net.SplitHostPort(hostport); err == nil {
			host = h
		}
		return allowed[strings.Trim(host, "[]")]
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			if !hostOK(r.Host) {
				http.Error(w, "forbidden host", http.StatusForbidden)
				return
			}
			if origin := r.Header.Get("Origin"); origin != "" {
				u, err := url.Parse(origin)
				if err != nil || u.Host != r.Host {
					http.Error(w, "forbidden origin", http.StatusForbidden)
					return
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}

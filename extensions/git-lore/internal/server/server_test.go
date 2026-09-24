package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGuardWrites(t *testing.T) {
	ok := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	h := guardWrites(ok, "127.0.0.1:9473")

	cases := []struct {
		name, method, host, origin string
		want                       int
	}{
		{"get from anywhere", "GET", "evil.example:9473", "", 200},
		{"same-origin put", "PUT", "127.0.0.1:9473", "http://127.0.0.1:9473", 200},
		{"localhost put no origin", "PUT", "localhost:9473", "", 200},
		{"rebound host", "PUT", "evil.example:9473", "http://evil.example:9473", 403},
		{"cross-site origin", "POST", "127.0.0.1:9473", "https://evil.example", 403},
	}
	for _, c := range cases {
		req := httptest.NewRequest(c.method, "/api/works/demo/files/plan.md", nil)
		req.Host = c.host
		if c.origin != "" {
			req.Header.Set("Origin", c.origin)
		}
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code != c.want {
			t.Errorf("%s: status %d, want %d", c.name, rec.Code, c.want)
		}
	}
}

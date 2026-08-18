package api

import (
	"encoding/json"
	"net/http"
	"strings"

	loregit "github.com/joshmcarthur/git-lore/extensions/git-lore/internal/git"
)

// Handler serves the lore-explorer JSON API.
type Handler struct {
	Repo *loregit.Repo
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func isNotFound(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "not found") ||
		strings.Contains(msg, "Needed a single revision") ||
		strings.Contains(msg, "does not exist") ||
		strings.Contains(msg, "is not in work") ||
		(strings.Contains(msg, "path") && strings.Contains(msg, "exists") && strings.Contains(msg, "does not")) ||
		(strings.Contains(msg, "fatal: Path") && strings.Contains(msg, "exists on disk"))
}

func isBadRequest(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "invalid") ||
		strings.Contains(msg, "empty") ||
		strings.Contains(msg, "must not") ||
		strings.Contains(msg, "must be")
}

func mapGitError(w http.ResponseWriter, err error) {
	switch {
	case err == nil:
		return
	case isBadRequest(err):
		writeError(w, http.StatusBadRequest, err.Error())
	case isNotFound(err):
		writeError(w, http.StatusNotFound, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, err.Error())
	}
}

// Mount registers API routes on mux.
func (h *Handler) Mount(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/repo", h.handleRepo)
	mux.HandleFunc("GET /api/branch-associations", h.handleBranchAssociations)
	mux.HandleFunc("GET /api/works", h.handleListWorks)
	mux.HandleFunc("GET /api/works/{id}", h.handleGetWork)
	mux.HandleFunc("GET /api/works/{id}/files/{path...}", h.handleShowFile)
	mux.HandleFunc("GET /api/works/{id}/log", h.handleLog)
	mux.HandleFunc("GET /api/works/{id}/commits/{sha}", h.handleCommit)
	mux.HandleFunc("GET /api/remote/status", h.handleRemoteStatus)
	mux.HandleFunc("POST /api/remote/fetch", h.handleRemoteFetch)
}

func (h *Handler) handleRepo(w http.ResponseWriter, r *http.Request) {
	info, err := h.Repo.Info()
	if err != nil {
		mapGitError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, info)
}

func (h *Handler) handleBranchAssociations(w http.ResponseWriter, r *http.Request) {
	assocs, err := h.Repo.BranchAssociations()
	if err != nil {
		mapGitError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, assocs)
}

func (h *Handler) handleListWorks(w http.ResponseWriter, r *http.Request) {
	works, err := h.Repo.ListWorks()
	if err != nil {
		mapGitError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, works)
}

func (h *Handler) handleGetWork(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	work, err := h.Repo.GetWork(id)
	if err != nil {
		mapGitError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, work)
}

func (h *Handler) handleShowFile(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	path := r.PathValue("path")
	content, err := h.Repo.ShowFile(id, path)
	if err != nil {
		mapGitError(w, err)
		return
	}
	ct := "text/plain; charset=utf-8"
	if strings.HasSuffix(path, ".md") {
		ct = "text/markdown; charset=utf-8"
	}
	w.Header().Set("Content-Type", ct)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(content))
}

func (h *Handler) handleLog(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	limit, offset := loregit.ParseLimitOffset(r.URL.Query().Get("limit"), r.URL.Query().Get("offset"))
	commits, err := h.Repo.Log(id, limit, offset)
	if err != nil {
		mapGitError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, commits)
}

func (h *Handler) handleCommit(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	sha := r.PathValue("sha")
	detail, err := h.Repo.CommitDetail(id, sha)
	if err != nil {
		mapGitError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, detail)
}

func (h *Handler) handleRemoteStatus(w http.ResponseWriter, r *http.Request) {
	remote := r.URL.Query().Get("remote")
	if remote == "" {
		remote = "origin"
	}
	status, err := h.Repo.RemoteStatus(remote)
	if err != nil {
		mapGitError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, status)
}

type fetchRequest struct {
	Remote string `json:"remote"`
	WorkID string `json:"workId"`
}

func (h *Handler) handleRemoteFetch(w http.ResponseWriter, r *http.Request) {
	var req fetchRequest
	if r.Body != nil && r.ContentLength != 0 {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}
	}
	if req.Remote == "" {
		req.Remote = "origin"
	}
	result, err := h.Repo.FetchLore(req.Remote, req.WorkID)
	if err != nil {
		mapGitError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

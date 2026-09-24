package git

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ErrConflict means the Work moved on since the caller loaded it.
var ErrConflict = errors.New("work has changed since it was loaded")

// WriteFile commits content at filePath onto refs/lore/<id>, keeping the rest
// of the tree. Mirrors lib/git-lore/edit.sh (isolated index, commit-tree,
// update-ref). baseCommit must be the Work's current commit; if the ref has
// moved, ErrConflict is returned and nothing is written. Returns the Work's
// commit after the write (unchanged when content is identical).
func (r *Repo) WriteFile(id, filePath, content, message, baseCommit string) (string, error) {
	if err := ValidateWorkID(id); err != nil {
		return "", err
	}
	if err := ValidateFilePath(filePath); err != nil {
		return "", err
	}
	if err := ValidateSHA(baseCommit); err != nil {
		return "", err
	}
	ref := "refs/lore/" + id
	out, err := r.Output("rev-parse", "--verify", ref)
	if err != nil {
		return "", fmt.Errorf("work %q not found: %w", id, err)
	}
	current := strings.TrimSpace(out)
	if !strings.HasPrefix(current, strings.ToLower(baseCommit)) {
		return "", ErrConflict
	}
	if message == "" {
		message = "lore: update Work " + id
	}

	blob, err := r.OutputWith(nil, content, "hash-object", "-w", "--stdin")
	if err != nil {
		return "", err
	}

	dir, err := os.MkdirTemp("", "git-lore-index-")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(dir)
	env := []string{"GIT_INDEX_FILE=" + filepath.Join(dir, "index")}

	if _, err := r.OutputWith(env, "", "read-tree", current); err != nil {
		return "", err
	}
	cacheinfo := "100644," + strings.TrimSpace(blob) + "," + filePath
	if _, err := r.OutputWith(env, "", "update-index", "--add", "--cacheinfo", cacheinfo); err != nil {
		return "", err
	}
	tree, err := r.OutputWith(env, "", "write-tree")
	if err != nil {
		return "", err
	}
	tree = strings.TrimSpace(tree)
	oldTree, err := r.Output("rev-parse", current+"^{tree}")
	if err != nil {
		return "", err
	}
	if tree == strings.TrimSpace(oldTree) {
		return current, nil
	}

	commit, err := r.Output("commit-tree", tree, "-p", current, "-m", message)
	if err != nil {
		return "", err
	}
	commit = strings.TrimSpace(commit)
	// Passing the old value makes the ref update a compare-and-swap, so a
	// concurrent edit-lore commit is never overwritten.
	if _, err := r.Output("update-ref", "-m", message, ref, commit, current); err != nil {
		return "", ErrConflict
	}
	return commit, nil
}

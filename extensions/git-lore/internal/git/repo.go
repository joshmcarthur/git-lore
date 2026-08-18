package git

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Repo runs git commands against a repository root.
type Repo struct {
	Root string
	Git  string
}

// Open discovers the git repository containing dir (or dir itself).
func Open(dir string) (*Repo, error) {
	if dir == "" {
		var err error
		dir, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}
	abs, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	r := &Repo{Root: abs, Git: "git"}
	out, err := r.Output("rev-parse", "--show-toplevel")
	if err != nil {
		return nil, fmt.Errorf("not a git repository: %s: %w", abs, err)
	}
	r.Root = strings.TrimSpace(out)
	return r, nil
}

// Output runs git -C <root> with args and returns stdout.
func (r *Repo) Output(args ...string) (string, error) {
	cmd := exec.Command(r.Git, append([]string{"-C", r.Root}, args...)...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return "", fmt.Errorf("git %s: %s", strings.Join(args, " "), msg)
	}
	return stdout.String(), nil
}

// Run runs git and returns combined stdout for success; on failure returns stderr.
func (r *Repo) Run(args ...string) (string, error) {
	return r.Output(args...)
}

package git

import (
	"errors"
	"os/exec"
	"strings"
	"testing"
)

func newLoreRepo(t *testing.T) (*Repo, string) {
	t.Helper()
	dir := t.TempDir()
	for _, args := range [][]string{
		{"init", "-q"},
		{"config", "user.name", "Test"},
		{"config", "user.email", "test@example.com"},
	} {
		if out, err := exec.Command("git", append([]string{"-C", dir}, args...)...).CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v: %s", args, err, out)
		}
	}
	r, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	plan, _ := r.OutputWith(nil, "# Plan\n", "hash-object", "-w", "--stdin")
	notes, _ := r.OutputWith(nil, "notes\n", "hash-object", "-w", "--stdin")
	tree, err := r.OutputWith(nil,
		"100644 blob "+strings.TrimSpace(plan)+"\tplan.md\n"+
			"100644 blob "+strings.TrimSpace(notes)+"\tnotes.md\n",
		"mktree")
	if err != nil {
		t.Fatal(err)
	}
	commit, err := r.Output("commit-tree", strings.TrimSpace(tree), "-m", "lore: init")
	if err != nil {
		t.Fatal(err)
	}
	commit = strings.TrimSpace(commit)
	if _, err := r.Output("update-ref", "refs/lore/demo", commit); err != nil {
		t.Fatal(err)
	}
	return r, commit
}

func TestWriteFileCommitsOnTopOfWork(t *testing.T) {
	r, base := newLoreRepo(t)

	commit, err := r.WriteFile("demo", "plan.md", "# Plan\n\nUpdated.\n", "lore: edit plan", base)
	if err != nil {
		t.Fatal(err)
	}
	if commit == base {
		t.Fatal("expected a new commit")
	}
	got, _ := r.ShowFile("demo", "plan.md")
	if got != "# Plan\n\nUpdated.\n" {
		t.Errorf("plan.md = %q", got)
	}
	if got, _ := r.ShowFile("demo", "notes.md"); got != "notes\n" {
		t.Errorf("notes.md not preserved: %q", got)
	}
	parent, _ := r.Output("rev-parse", commit+"^")
	if strings.TrimSpace(parent) != base {
		t.Errorf("parent = %s, want %s", parent, base)
	}
	subject, _ := r.Output("log", "-1", "--format=%s", "refs/lore/demo")
	if strings.TrimSpace(subject) != "lore: edit plan" {
		t.Errorf("subject = %q", subject)
	}
}

func TestWriteFileAddsNewFile(t *testing.T) {
	r, base := newLoreRepo(t)
	if _, err := r.WriteFile("demo", "docs/decisions.md", "# Decisions\n", "", base); err != nil {
		t.Fatal(err)
	}
	if got, _ := r.ShowFile("demo", "docs/decisions.md"); got != "# Decisions\n" {
		t.Errorf("docs/decisions.md = %q", got)
	}
	subject, _ := r.Output("log", "-1", "--format=%s", "refs/lore/demo")
	if strings.TrimSpace(subject) != "lore: update Work demo" {
		t.Errorf("default subject = %q", subject)
	}
}

func TestWriteFileUnchangedIsNoop(t *testing.T) {
	r, base := newLoreRepo(t)
	commit, err := r.WriteFile("demo", "plan.md", "# Plan\n", "", base)
	if err != nil {
		t.Fatal(err)
	}
	if commit != base {
		t.Errorf("commit = %s, want unchanged %s", commit, base)
	}
}

func TestWriteFileRejectsStaleBase(t *testing.T) {
	r, base := newLoreRepo(t)
	if _, err := r.WriteFile("demo", "plan.md", "first\n", "", base); err != nil {
		t.Fatal(err)
	}
	_, err := r.WriteFile("demo", "plan.md", "second\n", "", base)
	if !errors.Is(err, ErrConflict) {
		t.Fatalf("err = %v, want ErrConflict", err)
	}
	if got, _ := r.ShowFile("demo", "plan.md"); got != "first\n" {
		t.Errorf("plan.md = %q, stale write should not land", got)
	}
}

func TestWriteFileValidates(t *testing.T) {
	r, base := newLoreRepo(t)
	if _, err := r.WriteFile("demo", "../escape.md", "x", "", base); err == nil {
		t.Error("expected path error")
	}
	if _, err := r.WriteFile("Bad_ID", "plan.md", "x", "", base); err == nil {
		t.Error("expected work id error")
	}
	if _, err := r.WriteFile("missing", "plan.md", "x", "", base); err == nil {
		t.Error("expected not found error")
	}
}

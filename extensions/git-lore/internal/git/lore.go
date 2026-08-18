package git

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// WorkSummary is a local Lore Work at refs/lore/<id>.
type WorkSummary struct {
	ID      string `json:"id"`
	Ref     string `json:"ref"`
	Commit  string `json:"commit"`
	Subject string `json:"subject"`
}

// WorkDetail extends WorkSummary with document paths.
type WorkDetail struct {
	WorkSummary
	Files []string `json:"files"`
}

// CommitInfo is a lore history entry.
type CommitInfo struct {
	SHA     string `json:"sha"`
	Subject string `json:"subject"`
	Author  string `json:"author"`
	Date    string `json:"date"`
}

// CommitDetail includes parents and unified diff.
type CommitDetail struct {
	CommitInfo
	Parents []string `json:"parents"`
	Body    string   `json:"body"`
	Diff    string   `json:"diff"`
}

// BranchAssociation maps a local branch to a work-id.
type BranchAssociation struct {
	Branch string `json:"branch"`
	WorkID string `json:"workId"`
}

// RepoInfo is high-level repository metadata.
type RepoInfo struct {
	Root          string   `json:"root"`
	CurrentBranch string   `json:"currentBranch"`
	Remotes       []string `json:"remotes"`
}

// documentPriority matches read-lore structured order.
var documentPriority = map[string]int{
	"plan.md":          1,
	"decisions.md":     2,
	"questions.md":     3,
	"spec.md":          4,
	"investigation.md": 5,
	"architecture.md":  6,
	"handoff.md":       7,
}

func sortDocumentPaths(files []string) {
	sort.SliceStable(files, func(i, j int) bool {
		pi, oki := documentPriority[files[i]]
		pj, okj := documentPriority[files[j]]
		if !oki {
			pi = 100
		}
		if !okj {
			pj = 100
		}
		if pi != pj {
			return pi < pj
		}
		return files[i] < files[j]
	})
}

// Info returns repository root, branch, and remotes.
func (r *Repo) Info() (RepoInfo, error) {
	branch, err := r.Output("branch", "--show-current")
	if err != nil {
		return RepoInfo{}, err
	}
	remotesOut, err := r.Output("remote")
	if err != nil {
		return RepoInfo{}, err
	}
	var remotes []string
	for _, line := range strings.Split(strings.TrimSpace(remotesOut), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			remotes = append(remotes, line)
		}
	}
	if remotes == nil {
		remotes = []string{}
	}
	return RepoInfo{
		Root:          r.Root,
		CurrentBranch: strings.TrimSpace(branch),
		Remotes:       remotes,
	}, nil
}

// ListWorks returns local refs/lore/* Works, skipping *-remote side refs by default.
func (r *Repo) ListWorks() ([]WorkSummary, error) {
	out, err := r.Output(
		"for-each-ref",
		"--format=%(refname:short)%00%(objectname)%00%(contents:subject)",
		"refs/lore",
	)
	if err != nil {
		return nil, err
	}
	var works []WorkSummary
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "\x00", 3)
		if len(parts) < 2 {
			continue
		}
		short := parts[0]
		id := strings.TrimPrefix(short, "lore/")
		if id == short {
			// Unexpected format; skip
			continue
		}
		if strings.HasSuffix(id, "-remote") {
			continue
		}
		subject := ""
		if len(parts) > 2 {
			subject = parts[2]
		}
		works = append(works, WorkSummary{
			ID:      id,
			Ref:     "refs/lore/" + id,
			Commit:  parts[1],
			Subject: subject,
		})
	}
	if works == nil {
		works = []WorkSummary{}
	}
	sort.Slice(works, func(i, j int) bool { return works[i].ID < works[j].ID })
	return works, nil
}

// GetWork returns detail for a single Work.
func (r *Repo) GetWork(id string) (WorkDetail, error) {
	if err := ValidateWorkID(id); err != nil {
		return WorkDetail{}, err
	}
	ref := "refs/lore/" + id
	commit, err := r.Output("rev-parse", ref)
	if err != nil {
		return WorkDetail{}, fmt.Errorf("work %q not found: %w", id, err)
	}
	subject, err := r.Output("log", "-1", "--format=%s", ref)
	if err != nil {
		return WorkDetail{}, err
	}
	filesOut, err := r.Output("ls-tree", "-r", "--name-only", ref)
	if err != nil {
		return WorkDetail{}, err
	}
	var files []string
	for _, f := range strings.Split(strings.TrimSpace(filesOut), "\n") {
		if f != "" {
			files = append(files, f)
		}
	}
	if files == nil {
		files = []string{}
	}
	sortDocumentPaths(files)
	return WorkDetail{
		WorkSummary: WorkSummary{
			ID:      id,
			Ref:     ref,
			Commit:  strings.TrimSpace(commit),
			Subject: strings.TrimSpace(subject),
		},
		Files: files,
	}, nil
}

// ShowFile returns the content of a path in a Lore Work.
func (r *Repo) ShowFile(id, filePath string) (string, error) {
	if err := ValidateWorkID(id); err != nil {
		return "", err
	}
	if err := ValidateFilePath(filePath); err != nil {
		return "", err
	}
	return r.Output("show", fmt.Sprintf("refs/lore/%s:%s", id, filePath))
}

// Log returns paginated lore history.
func (r *Repo) Log(id string, limit, offset int) ([]CommitInfo, error) {
	if err := ValidateWorkID(id); err != nil {
		return nil, err
	}
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	ref := "refs/lore/" + id
	out, err := r.Output(
		"log",
		fmt.Sprintf("--skip=%d", offset),
		fmt.Sprintf("-%d", limit),
		"--format=%H%x00%s%x00%an%x00%aI",
		ref,
	)
	if err != nil {
		return nil, err
	}
	var commits []CommitInfo
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "\x00", 4)
		if len(parts) < 4 {
			continue
		}
		commits = append(commits, CommitInfo{
			SHA:     parts[0],
			Subject: parts[1],
			Author:  parts[2],
			Date:    parts[3],
		})
	}
	if commits == nil {
		commits = []CommitInfo{}
	}
	return commits, nil
}

// CommitDetail returns message, parents, and unified diff for a commit.
func (r *Repo) CommitDetail(id, sha string) (CommitDetail, error) {
	if err := ValidateWorkID(id); err != nil {
		return CommitDetail{}, err
	}
	if err := ValidateSHA(sha); err != nil {
		return CommitDetail{}, err
	}
	// Ensure sha is in this work's history.
	_, err := r.Output("merge-base", "--is-ancestor", sha, "refs/lore/"+id)
	if err != nil {
		return CommitDetail{}, fmt.Errorf("commit %q is not in work %q", sha, id)
	}

	meta, err := r.Output("log", "-1", "--format=%H%x00%s%x00%an%x00%aI%x00%P%x00%b", sha)
	if err != nil {
		return CommitDetail{}, err
	}
	parts := strings.SplitN(strings.TrimSpace(meta), "\x00", 6)
	if len(parts) < 5 {
		return CommitDetail{}, fmt.Errorf("unexpected log format for %s", sha)
	}
	var parents []string
	if parts[4] != "" {
		parents = strings.Fields(parts[4])
	}
	if parents == nil {
		parents = []string{}
	}
	body := ""
	if len(parts) > 5 {
		body = parts[5]
	}

	diffOut, err := r.Output("show", "--format=", "--patch", sha)
	if err != nil {
		return CommitDetail{}, err
	}
	diff := diffOut

	return CommitDetail{
		CommitInfo: CommitInfo{
			SHA:     parts[0],
			Subject: parts[1],
			Author:  parts[2],
			Date:    parts[3],
		},
		Parents: parents,
		Body:    body,
		Diff:    diff,
	}, nil
}

// BranchAssociations returns local branch.*.lore mappings.
func (r *Repo) BranchAssociations() ([]BranchAssociation, error) {
	out, err := r.Output("config", "--get-regexp", `^branch\..*\.lore$`)
	if err != nil {
		// git config exits 1 when no matches
		if strings.Contains(err.Error(), "exit status 1") || strings.TrimSpace(out) == "" {
			return []BranchAssociation{}, nil
		}
		return nil, err
	}
	var assocs []BranchAssociation
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		if line == "" {
			continue
		}
		fields := strings.SplitN(line, " ", 2)
		if len(fields) != 2 {
			continue
		}
		key := fields[0] // branch.<name>.lore
		workID := fields[1]
		branch := strings.TrimSuffix(strings.TrimPrefix(key, "branch."), ".lore")
		assocs = append(assocs, BranchAssociation{Branch: branch, WorkID: workID})
	}
	if assocs == nil {
		assocs = []BranchAssociation{}
	}
	return assocs, nil
}

// ParseLimitOffset parses pagination query values with defaults.
func ParseLimitOffset(limitStr, offsetStr string) (limit, offset int) {
	limit = 10
	offset = 0
	if limitStr != "" {
		if n, err := strconv.Atoi(limitStr); err == nil && n > 0 {
			limit = n
			if limit > 100 {
				limit = 100
			}
		}
	}
	if offsetStr != "" {
		if n, err := strconv.Atoi(offsetStr); err == nil && n >= 0 {
			offset = n
		}
	}
	return limit, offset
}

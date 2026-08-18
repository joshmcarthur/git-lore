package git

import (
	"fmt"
	"path"
	"regexp"
	"strings"
)

var workIDPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

// ValidateWorkID checks that work-id matches the lore protocol slug rules.
func ValidateWorkID(id string) error {
	if id == "" {
		return fmt.Errorf("work id is empty")
	}
	if !workIDPattern.MatchString(id) {
		return fmt.Errorf("invalid work id %q: must be lowercase letters, digits, and hyphens", id)
	}
	return nil
}

// ValidateFilePath rejects path traversal and absolute paths.
func ValidateFilePath(p string) error {
	if p == "" {
		return fmt.Errorf("file path is empty")
	}
	if strings.Contains(p, `\`) {
		return fmt.Errorf("file path must use forward slashes: %q", p)
	}
	if strings.HasPrefix(p, "/") {
		return fmt.Errorf("file path must be relative: %q", p)
	}
	for _, part := range strings.Split(p, "/") {
		if part == ".." {
			return fmt.Errorf("file path must not contain ..: %q", p)
		}
	}
	cleaned := path.Clean("/" + p) // absolute clean avoids "." surprises
	cleaned = strings.TrimPrefix(cleaned, "/")
	if cleaned == ".." || strings.HasPrefix(cleaned, "../") {
		return fmt.Errorf("file path must not contain ..: %q", p)
	}
	return nil
}

// ValidateSHA accepts abbreviated or full hex object names.
func ValidateSHA(sha string) error {
	if sha == "" {
		return fmt.Errorf("sha is empty")
	}
	if len(sha) < 4 || len(sha) > 40 {
		return fmt.Errorf("invalid sha length: %q", sha)
	}
	for _, c := range sha {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return fmt.Errorf("invalid sha: %q", sha)
		}
	}
	return nil
}

// ValidateRemoteName rejects empty or unsafe remote names.
func ValidateRemoteName(name string) error {
	if name == "" {
		return fmt.Errorf("remote name is empty")
	}
	if strings.ContainsAny(name, " \t\n\r/") {
		return fmt.Errorf("invalid remote name: %q", name)
	}
	return nil
}

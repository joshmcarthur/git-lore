package git

import "testing"

func TestValidateWorkID(t *testing.T) {
	ok := []string{"git-lore", "a", "oauth2", "create-lore"}
	for _, id := range ok {
		if err := ValidateWorkID(id); err != nil {
			t.Errorf("ValidateWorkID(%q): %v", id, err)
		}
	}
	bad := []string{"", "Git-Lore", "has_underscore", "-leading", "trailing-", "a/b", "spaces here"}
	for _, id := range bad {
		if err := ValidateWorkID(id); err == nil {
			t.Errorf("ValidateWorkID(%q): expected error", id)
		}
	}
}

func TestValidateFilePath(t *testing.T) {
	ok := []string{"plan.md", "docs/spec.md", "a/b/c.md"}
	for _, p := range ok {
		if err := ValidateFilePath(p); err != nil {
			t.Errorf("ValidateFilePath(%q): %v", p, err)
		}
	}
	bad := []string{"", "/abs", "../escape", "a/../../b", `..\win`}
	for _, p := range bad {
		if err := ValidateFilePath(p); err == nil {
			t.Errorf("ValidateFilePath(%q): expected error", p)
		}
	}
}

func TestClassifySync(t *testing.T) {
	cases := []struct {
		local, remote, base string
		want                SyncState
	}{
		{"aaa", "aaa", "", StateSynced},
		{"aaa", "", "", StateLocalOnly},
		{"", "bbb", "", StateRemoteOnly},
		{"ccc", "bbb", "bbb", StateAhead},
		{"bbb", "ccc", "bbb", StateBehind},
		{"aaa", "bbb", "zzz", StateDiverged},
		{"aaa", "bbb", "", StateDiverged},
	}
	for _, tc := range cases {
		got := ClassifySync(tc.local, tc.remote, tc.base)
		if got != tc.want {
			t.Errorf("ClassifySync(%q,%q,%q)=%q want %q", tc.local, tc.remote, tc.base, got, tc.want)
		}
	}
}

func TestSortDocumentPaths(t *testing.T) {
	files := []string{"z.md", "spec.md", "plan.md", "questions.md", "decisions.md"}
	sortDocumentPaths(files)
	want := []string{"plan.md", "decisions.md", "questions.md", "spec.md", "z.md"}
	for i := range want {
		if files[i] != want[i] {
			t.Fatalf("sortDocumentPaths = %v want %v", files, want)
		}
	}
}

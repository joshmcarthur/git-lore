package git

import (
	"fmt"
	"strings"
)

// SyncState classifies local vs remote lore ref relationship.
type SyncState string

const (
	StateSynced     SyncState = "synced"
	StateLocalOnly  SyncState = "local-only"
	StateRemoteOnly SyncState = "remote-only"
	StateAhead      SyncState = "ahead"
	StateBehind     SyncState = "behind"
	StateDiverged   SyncState = "diverged"
)

// WorkRemoteStatus is sync classification for one work-id.
type WorkRemoteStatus struct {
	ID         string    `json:"id"`
	LocalSHA   string    `json:"localSha,omitempty"`
	RemoteSHA  string    `json:"remoteSha,omitempty"`
	State      SyncState `json:"state"`
}

// RemoteStatusResult is the response for remote status.
type RemoteStatusResult struct {
	Remote           string             `json:"remote"`
	RefspecConfigured bool              `json:"refspecConfigured"`
	Works            []WorkRemoteStatus `json:"works"`
}

// ClassifySync determines the relationship between local and remote SHAs.
// Empty local or remote means that side is absent.
func ClassifySync(localSHA, remoteSHA, mergeBase string) SyncState {
	switch {
	case localSHA == "" && remoteSHA == "":
		return StateSynced // nothing to compare
	case localSHA == "" && remoteSHA != "":
		return StateRemoteOnly
	case localSHA != "" && remoteSHA == "":
		return StateLocalOnly
	case localSHA == remoteSHA:
		return StateSynced
	case mergeBase == remoteSHA:
		return StateAhead
	case mergeBase == localSHA:
		return StateBehind
	default:
		return StateDiverged
	}
}

// HasLoreFetchRefspec reports whether remote.<name>.fetch includes refs/lore/.
func (r *Repo) HasLoreFetchRefspec(remote string) (bool, error) {
	if err := ValidateRemoteName(remote); err != nil {
		return false, err
	}
	out, err := r.Output("config", "--get-all", "remote."+remote+".fetch")
	if err != nil {
		// No fetch config configured
		if strings.Contains(err.Error(), "exit status 1") {
			return false, nil
		}
		return false, err
	}
	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, "refs/lore/") {
			return true, nil
		}
	}
	return false, nil
}

// EnsureLoreFetchRefspec adds refs/lore/*:refs/lore/* (no +) if missing.
func (r *Repo) EnsureLoreFetchRefspec(remote string) (added bool, err error) {
	ok, err := r.HasLoreFetchRefspec(remote)
	if err != nil {
		return false, err
	}
	if ok {
		return false, nil
	}
	_, err = r.Output("config", "--add", "remote."+remote+".fetch", "refs/lore/*:refs/lore/*")
	if err != nil {
		return false, err
	}
	return true, nil
}

// RemoteLoreRefs returns work-id → SHA from ls-remote.
func (r *Repo) RemoteLoreRefs(remote string) (map[string]string, error) {
	if err := ValidateRemoteName(remote); err != nil {
		return nil, err
	}
	out, err := r.Output("ls-remote", remote, "refs/lore/*")
	if err != nil {
		return nil, err
	}
	refs := make(map[string]string)
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		sha, ref := fields[0], fields[1]
		id := strings.TrimPrefix(ref, "refs/lore/")
		if id == ref || strings.HasSuffix(id, "-remote") {
			continue
		}
		refs[id] = sha
	}
	return refs, nil
}

// RemoteStatus compares local and remote lore refs.
func (r *Repo) RemoteStatus(remote string) (RemoteStatusResult, error) {
	if err := ValidateRemoteName(remote); err != nil {
		return RemoteStatusResult{}, err
	}
	configured, err := r.HasLoreFetchRefspec(remote)
	if err != nil {
		return RemoteStatusResult{}, err
	}

	localWorks, err := r.ListWorks()
	if err != nil {
		return RemoteStatusResult{}, err
	}
	remoteRefs, err := r.RemoteLoreRefs(remote)
	if err != nil {
		return RemoteStatusResult{}, err
	}

	ids := make(map[string]struct{})
	localSHAs := make(map[string]string)
	for _, w := range localWorks {
		ids[w.ID] = struct{}{}
		localSHAs[w.ID] = w.Commit
	}
	for id := range remoteRefs {
		ids[id] = struct{}{}
	}

	var works []WorkRemoteStatus
	for id := range ids {
		local := localSHAs[id]
		remoteSHA := remoteRefs[id]
		var mergeBase string
		if local != "" && remoteSHA != "" && local != remoteSHA {
			mb, err := r.Output("merge-base", local, remoteSHA)
			if err == nil {
				mergeBase = strings.TrimSpace(mb)
			}
			// If merge-base fails (unrelated histories), treat as diverged with empty base.
		}
		works = append(works, WorkRemoteStatus{
			ID:        id,
			LocalSHA:  local,
			RemoteSHA: remoteSHA,
			State:     ClassifySync(local, remoteSHA, mergeBase),
		})
	}

	// Stable sort by id
	for i := 0; i < len(works); i++ {
		for j := i + 1; j < len(works); j++ {
			if works[j].ID < works[i].ID {
				works[i], works[j] = works[j], works[i]
			}
		}
	}
	if works == nil {
		works = []WorkRemoteStatus{}
	}

	return RemoteStatusResult{
		Remote:            remote,
		RefspecConfigured: configured,
		Works:             works,
	}, nil
}

// FetchResult reports fetch outcome.
type FetchResult struct {
	Remote      string `json:"remote"`
	WorkID      string `json:"workId,omitempty"`
	RefspecAdded bool  `json:"refspecAdded"`
	Output      string `json:"output"`
}

// FetchLore fetches lore refs from remote. If workID is non-empty, fetches that Work only.
func (r *Repo) FetchLore(remote, workID string) (FetchResult, error) {
	if err := ValidateRemoteName(remote); err != nil {
		return FetchResult{}, err
	}
	if workID != "" {
		if err := ValidateWorkID(workID); err != nil {
			return FetchResult{}, err
		}
	}

	added, err := r.EnsureLoreFetchRefspec(remote)
	if err != nil {
		return FetchResult{}, err
	}

	var refspec string
	if workID != "" {
		refspec = fmt.Sprintf("refs/lore/%s:refs/lore/%s", workID, workID)
	} else {
		refspec = "refs/lore/*:refs/lore/*"
	}

	out, err := r.Output("fetch", remote, refspec)
	result := FetchResult{
		Remote:       remote,
		WorkID:       workID,
		RefspecAdded: added,
		Output:       strings.TrimSpace(out),
	}
	if err != nil {
		return result, err
	}
	return result, nil
}

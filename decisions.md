# Decisions

## Read-only v1

**Decision:** lore-explorer is read + fetch only. Create, edit, and push stay with agent skills.

**Why:** Matches the protocol decision to prove workflow before building infrastructure. Writing from a UI risks noisy lore without the edit-lore curation gate. Fetch is included so remote-only Works can be discovered and materialised locally (sync-lore semantics).

## Shell out to git, do not use libgit2

**Decision:** All lore operations run via `os/exec` against the user's `git` binary (`git -C <repo> …`).

**Why:** Lore is defined as ordinary Git refs and objects. Using the same CLI as the skills keeps behaviour identical (refspecs, merge-base classification, orphan commits) and avoids a native dependency.

## Embed Vue build in the Go binary

**Decision:** Vite builds into `internal/server/webdist`; Go embeds that tree with `go:embed`. One binary serves API + UI.

**Why:** Distributing a single `lore-explorer` binary matches "standalone" without requiring Node at runtime. `make build` rebuilds the SPA then the binary; CI runs the same pipeline.

## Work id `lore-explorer`

**Decision:** Lore Work id and binary name are both `lore-explorer`.

**Why:** Clear, semantic, and matches the create-lore preference for short ids over branch-derived slugs (`docs-skill-install-protocol-refs`).

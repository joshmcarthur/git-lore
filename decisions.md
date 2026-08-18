# Decisions

## Read-only serve UI (v1)

**Decision:** The `explore` UI is read + fetch only. Create, edit, and push stay with agent skills until dedicated CLI subcommands exist.

**Why:** Matches the protocol decision to prove workflow before building infrastructure. Writing from a UI risks noisy lore without the edit-lore curation gate. Fetch is included so remote-only Works can be discovered and materialised locally (sync-lore semantics).

## Shell out to git, do not use libgit2

**Decision:** All lore operations run via `os/exec` against the user's `git` binary (`git -C <repo> …`).

**Why:** Lore is defined as ordinary Git refs and objects. Using the same CLI as the skills keeps behaviour identical (refspecs, merge-base classification, orphan commits) and avoids a native dependency. Future CLI subcommands should reuse the same `internal/git` wrappers.

## Embed Vue build in the Go binary

**Decision:** Vite builds into `internal/server/webdist`; Go embeds that tree with `go:embed`. One binary serves API + UI for `explore`.

**Why:** Distributing a single `git-lore` binary matches "standalone" without requiring Node at runtime. `make build` rebuilds the SPA then the binary; CI runs the same pipeline.

## Lore Work id remains `lore-explorer`

**Decision:** Keep Lore Work id `refs/lore/lore-explorer` even though the binary/extension is now `git-lore`.

**Why:** Lore refs should not be renamed casually; history and branch association already point here. The Work covers the explorer UI origin and the broader CLI direction.

## Live under `extensions/`, not repository root

**Decision:** Place the Go module under `extensions/git-lore/`. Keep repository root skills-first.

**Why:** git-lore is primarily an agent-skills repository. The CLI is an optional extension for humans (and future non-skill users). Namespacing under `extensions/` keeps clones focused on `skills/` for install paths.

## Binary and CLI name: `git-lore`

**Decision:** Produce `bin/git-lore` with subcommands. Command: `explore` (aliases `ui`, `serve`). Extension path: `extensions/git-lore/`.

**Why:** Positions the tool as a general Lore CLI that can grow beyond a browser — abstracting the same Git conventions the skills use for people who do not want agent skills. `lore-explorer` as a binary name locked the tool to one UI; `git-lore explore` leaves room for `list` / `show` / `create` / `sync` without a second binary.

## Gitignore compiled SPA output

**Decision:** Do not commit `internal/server/webdist/` (or `bin/`). Regenerate with `make web` / `make build` before compile.

**Why:** Hashed Vite assets churn on every frontend change and are not source. `go:embed` requires the directory at compile time — Makefile `ensure-web` / `build` targets create it. CI path-filters on `extensions/git-lore/**`.

## Stdlib subcommands before a CLI framework

**Decision:** Dispatch commands with a small `switch` on `os.Args` and per-command `flag.FlagSet`s. Do not add Cobra/urfave yet.

**Why:** Only one real command exists. A framework can wait until subcommand count or flag complexity justifies the dependency.

## Command name: `serve`

**Decision:** The UI subcommand is `git-lore serve` only — no `explore` / `ui` aliases.

**Why:** One clear verb. Aliases add surface area without helping discovery.

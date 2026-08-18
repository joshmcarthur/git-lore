# Handoff

## What

Optional `git-lore` CLI under `extensions/git-lore/`. Full implementation exists on `feat/git-lore-cli` (and `origin/feat/git-lore-cli`). Monolithic PR https://github.com/joshmcarthur/git-lore/pull/4 was **closed** to split review.

## Why split

Backend (Go API/CLI) and frontend (Vue + release binary job) are separately reviewable. `go:embed` needs a placeholder until Vite `webdist` lands.

## Approved stack (not executed yet)

1. **PR1 → main:** Go CLI + `internal/git` + API + `serve` with committed plain HTML placeholder; Go-only CI/Makefile; no Vue; no release-please binary upload
2. **PR2 → PR1:** `web/` Vue SPA, Makefile `web`/`dist`, release-please `release-binaries`, fuller CI + README

Backup: do not delete `feat/git-lore-cli` until the stack is open and verified.

## Decided (durable)

See `decisions.md` — especially extensions/ layout, `serve` command name, gitignore webdist, release-please binaries, and stacked PRs.

## Next agent action

Execute the split: create PR1/PR2 branches from the backup commit, push, open PRs, leave `feat/git-lore-cli` intact.

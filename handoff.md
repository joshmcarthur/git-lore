# Handoff

## What

Optional `git-lore` CLI under `extensions/git-lore/`. First command: `serve` (embedded Vue Lore browser, read + fetch). Release-please uploads cross-platform archives on each GitHub Release.

## Why

Skills remain primary. The CLI is for humans who prefer commands; `serve` proves browsing without agent skills. Same Git conventions as the skills (`os/exec` to `git`).

## Decided

- Live under `extensions/`, not repo root
- Binary name `git-lore`; UI command `serve` only
- Read-only serve UI; writes stay skill-driven for now
- Gitignore `webdist`; `make build` / `make dist` regenerate
- Release binaries via release-please follow-on job (not separate goreleaser)

## Next

- Merge source PR
- Sync lore refs if not already on origin
- Later: `list` / `show` CLI before create/edit/sync

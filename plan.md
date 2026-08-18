# Plan

## Objective

Build a standalone Go binary (`lore-explorer`) that serves an embedded Vue SPA for exploring local and remote Lore Works. The tool uses ordinary `git` commands — no libgit2, no lore files in the working tree — so humans can browse `refs/lore/*` the same way agents use read-lore and sync-lore.

## Status

**Implemented.** `cmd/lore-explorer` embeds a Vue SPA and exposes a read-only HTTP API over `git` for listing Works, viewing documents, history/diffs, remote status, and fetch.

## Scope

**In scope:**

- List local Lore Works and branch→work associations
- View lore documents (Markdown) in read-lore priority order
- Browse lore commit history and per-commit diffs
- Remote status (synced / ahead / behind / diverged / local-only / remote-only)
- Fetch lore refs from a remote (configure refspec if needed)

**Out of scope:**

- Create or edit lore from the UI
- Push lore refs
- AI handoff synthesis (remains read-lore agent-side)
- MCP server, hosted service, or database

## Next steps

1. Dogfood the UI against this repository's Lore Works
2. Consider release binaries (goreleaser) if local `make build` proves useful
3. Revisit write/push UI only if skill-driven create/edit/sync remains painful

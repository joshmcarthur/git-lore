# Plan

## Objective

Build a standalone Go binary (`lore-explorer`) that serves an embedded Vue SPA for exploring local and remote Lore Works. The tool uses ordinary `git` commands — no libgit2, no lore files in the working tree — so humans can browse `refs/lore/*` the same way agents use read-lore and sync-lore.

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

1. Scaffold Go module, git wrapper, and HTTP API
2. Build Vue SPA (works list, document viewer, history/diff)
3. Add remote status + fetch endpoints
4. Makefile embed pipeline, CI build, smoke test against this repo

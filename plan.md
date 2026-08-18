# Plan

## Objective

Build an optional **`git-lore` CLI** that abstracts Lore operations (`refs/lore/*`) for people who prefer commands over agent skills. The first command is `serve`: eventually an embedded Vue SPA for browsing Lore via ordinary `git`. Later commands can wrap create / read / edit / sync without requiring an agent.

## Status

**Implemented on branch `feat/git-lore-cli` (backup); monolithic PR #4 closed.** Splitting into stacked PRs before merge:

1. Backend — CLI + API + placeholder HTML embed
2. Frontend — Vue UI + release binaries (stacked on 1)

## Scope

**In scope (now / near-term PRs):**

- CLI entrypoint (`serve`, `help`, `version`)
- Lore HTTP API over `git`
- Vue `serve` UI + gitignored `webdist`
- Release assets via release-please

**In scope (later):**

- CLI counterparts to skill workflows (`list`, `show`, `create`, `edit`, `sync`, …)

**Out of scope (for now):**

- Create/edit/push from the serve UI
- AI handoff synthesis
- MCP / hosted service
- Root-level product promotion (stays under `extensions/`)

## Layout

```
extensions/git-lore/   # self-contained Go module (+ Vue in frontend PR)
skills/                # primary repository surface
```

## Next steps

1. Open stacked PR1 (backend) and PR2 (frontend) from the approved split plan
2. Merge PR1 → then PR2; retire `feat/git-lore-cli` when stack lands
3. Dogfood `git-lore serve`; later add `list` / `show` before write commands

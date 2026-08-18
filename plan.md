# Plan

## Objective

Build an optional **`git-lore` CLI** that abstracts Lore operations (`refs/lore/*`) for people who prefer commands over agent skills. The first command is `serve`: an embedded Vue SPA for browsing local and remote Lore Works via ordinary `git`. Later commands can wrap create / read / edit / sync without requiring an agent.

## Status

**In progress — CLI shell + serve implemented; release binaries wired.** Lives under `extensions/git-lore/`. Binary: `bin/git-lore`. Cross-platform archives attach to GitHub Releases via release-please. Repository root remains skills-first.

## Scope

**In scope (now):**

- CLI entrypoint with subcommands (`serve`, `help`, `version`)
- `serve`: list Works, Markdown view, history/diffs, remote status, fetch
- Release assets for linux/darwin/windows on each GitHub Release

**In scope (later):**

- CLI counterparts to skill workflows (`list`, `show`, `create`, `edit`, `sync`, …)

**Out of scope (for now):**

- Create/edit/push from the serve UI
- AI handoff synthesis
- MCP server, hosted service, or database
- Promoting the CLI to a first-class root-level product surface

## Layout

```
extensions/git-lore/   # self-contained Go module + Vue (serve UI)
skills/                # primary repository surface
```

Build: `cd extensions/git-lore && make build` → `bin/git-lore`

## Next steps

1. Dogfood `git-lore serve` against this repository
2. Add read-only CLI commands (`list`, `show`) before write commands
3. Mirror skill semantics carefully when adding create/edit/sync

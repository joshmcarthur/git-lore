# Plan

## Objective

Build an optional **`git-lore` CLI** that abstracts Lore operations (`refs/lore/*`) for people who prefer commands over agent skills. The first command is `serve`: an embedded Vue SPA for browsing Lore via ordinary `git`. Later commands can wrap create / read / edit / sync without requiring an agent.

## Status

**Stacked PRs open.** Monolithic PR #4 closed. Backup: `feat/git-lore-cli` / `refs/backup/feat-git-lore-cli-*`.

1. [#5](https://github.com/joshmcarthur/git-lore/pull/5) — backend CLI + API + placeholder HTML → `main`
2. [#6](https://github.com/joshmcarthur/git-lore/pull/6) — Vue UI + release binaries → #5

## Scope

Unchanged: CLI under `extensions/git-lore/`; skills remain primary.

## Next steps

1. Review and merge #5, then #6
2. Retire `feat/git-lore-cli` after stack lands
3. Dogfood `serve`; later add `list` / `show` before write commands

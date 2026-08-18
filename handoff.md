# Handoff

## What

Stacked PRs for the optional `git-lore` CLI:

| PR | Branch | Base | Contents |
|----|--------|------|----------|
| [#5](https://github.com/joshmcarthur/git-lore/pull/5) | `feat/git-lore-cli-backend` | `main` | Go CLI, lore API, placeholder HTML, Go CI |
| [#6](https://github.com/joshmcarthur/git-lore/pull/6) | `feat/git-lore-cli-ui` | backend | Vue SPA, `make web`/`dist`, release-please binaries |

Closed: [#4](https://github.com/joshmcarthur/git-lore/pull/4) (monolith). Backup branch: `feat/git-lore-cli`.

## Merge order

Merge #5 first. Then merge #6 (retarget to `main` after #5 lands if GitHub does not auto-update the base).

## Next

Review PRs; after merge, delete backup branch when comfortable.

# Plan

## Objective

Implement the sync-lore agent skill: transport Lore Works between local
repository and Git remotes via ordinary Git fetch/push on `refs/lore/*`.
An agent following this skill should be able to configure refspecs, fetch
and push lore refs, detect transport problems, and report sync status —
without an executable, MCP server, or proprietary service.

sync-lore is the fourth and final core git-lore skill. create-lore,
read-lore, and edit-lore all defer remote transport to this skill.

## Status

**Done.** `skills/sync-lore/SKILL.md` implemented per spec.

## What we're trying to prove

A developer can share Lore Works with teammates through the same Git remote
used for source code, using only standard Git transport. Lore refs survive
push/fetch when refspecs are configured. A fresh clone can obtain lore refs
after one-time fetch configuration.

The test: configure refspecs, push local lore refs, clone elsewhere,
configure refspecs, fetch, and read lore with `git show refs/lore/<id>:plan.md`.

## Scope

Implement:

- `skills/sync-lore/SKILL.md` with YAML frontmatter and step-by-step workflow
- One-time per-remote fetch refspec configuration
- Fetch, push, and combined sync modes
- Status/discovery: list local vs remote lore refs
- Divergence detection when push is rejected
- Clone and post-clone setup documentation
- Link to `skills/protocol.md` and Lore Work `refs/lore/sync-lore`

Do not implement in this work:

- Automatic lore history merge (detect and document; manual reconciliation)
- GitHub-specific integrations
- Executable / CLI
- Lore ref deletion or archival
- Transport of non-lore refs

## Success criterion

A fresh agent with only sync-lore and no conversational history should be able to:

1. Configure lore fetch refspec for a remote
2. Push all local lore refs to remote
3. Fetch lore refs from remote on another clone
4. Detect when local and remote lore refs have diverged
5. Confirm working tree remains clean throughout

## Next steps

1. Dogfood: push lore refs when a remote exists
2. Smoke test clone + fetch in a temp directory
3. Update `refs/lore/git-lore` plan to reflect sync-lore complete

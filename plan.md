# Plan

## Objective

Implement the edit-lore agent skill: the primary way agents curate durable
context in an existing Lore Work. edit-lore modifies documents at
`refs/lore/<work-id>` via a new lore commit, without polluting the working
tree or source-branch history.

edit-lore is the skill that makes Lore useful over time. create-lore starts
a Work; read-lore consumes it; edit-lore maintains it.

## Status

**Done.** `skills/edit-lore/SKILL.md` implemented per spec.md.

## What we're trying to prove

The git-lore UX risk (agents writing too little or too much) is addressed
primarily by edit-lore's curation gate and document-emergence rules.

The test: after a session producing durable conclusions (decisions,
discoveries, plan changes), an agent can update Lore in one curated commit
that a fresh read-lore invocation surfaces clearly — without transcript
noise.

## Scope

Implement:

- `skills/edit-lore/SKILL.md` with YAML frontmatter and step-by-step workflow
- Work-id resolution (same order as read-lore)
- Mandatory curation gate before writing
- Temp work-tree export/edit/commit workflow
- Lore commit with parent linkage and `lore: <what changed>` message
- Document emergence rules (when to create decisions.md, questions.md, etc.)
- Post-edit verification and diff summary
- Link to `skills/protocol.md` and Lore Work `refs/lore/git-lore`

Do not implement in this work:

- sync-lore skill
- Automatic lore updates without explicit user/agent intent
- Merge/divergence resolution for lore refs (sync-lore + manual merge)
- Executable or helper scripts
- Lore deletion or archival workflows

## Working hypothesis

edit-lore should **read before write**: the agent must inspect current lore
state (via read-lore or `git show`) before editing, to avoid overwriting
curated content or duplicating decisions.

Batching related curation into a single lore commit is preferable to one
commit per sentence — Lore history should reflect meaningful curation
events, not activity telemetry.

## Success criterion

A fresh agent with only edit-lore should be able to:

1. Resolve work-id and load current lore tree into a temp directory
2. Add a decision to `decisions.md` (creating the file if absent)
3. Update `plan.md` next steps
4. Commit a single lore revision with a clear `lore:` message
5. Verify via `git diff` that working tree is clean and lore ref advanced

Dogfood: use edit-lore workflow to maintain `refs/lore/git-lore` and
per-skill Lore Works during implementation.

## Next steps

1. Commit `skills/edit-lore/SKILL.md` to the source branch
2. Run sync-lore when a remote Lore Work exists
3. Dogfood: update `refs/lore/git-lore` plan status via edit-lore workflow
4. Associate current branch via `git config branch.<name>.lore edit-lore`

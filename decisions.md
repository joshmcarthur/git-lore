# Decisions

## Temp work-tree for lore edits

**Decision:** edit-lore exports the lore tree to a temp directory via
`git --work-tree=<tmpdir> checkout <lore-commit> -- .`, edits there, then
`git add -A && write-tree && commit-tree -p <parent>`.

create-lore uses blob/mktree plumbing for the initial orphan commit (often
a single file). edit-lore routinely touches multiple files and benefits from
normal file editing in a temp directory.

Both approaches keep lore out of the repository working tree.

## Mandatory curation gate

**Decision:** edit-lore requires an explicit curation gate before any write.
The gate applies the protocol heuristic and category rules from
`skills/protocol.md`.

This is the primary control for the "too little / too much lore" UX risk
identified in `refs/lore/git-lore:investigation.md`. read-lore presents;
edit-lore curates.

## Read before write

**Decision:** Agents must inspect current lore content before editing.
Minimum: `git show` on target files. Prefer read-lore for complex edits.

Prevents overwriting curated content and duplicate decision entries.

## Batch related curation in one lore commit

**Decision:** One lore commit per meaningful curation event, not per line
or per conversation turn. Related changes (e.g. decision + plan update)
may share a commit.

Lore history should read as a curation log, not an activity feed.

## Document emergence, not preemption

**Decision:** edit-lore creates new artefacts (decisions.md, questions.md,
etc.) only when the work produces content that fits. Same rule as create-lore
but applied at edit time.

Matches protocol decision: no fixed document taxonomy.

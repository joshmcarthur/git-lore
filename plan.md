# Plan

## Objective

Implement the create-lore agent skill: the entry point for starting a new
Lore Work. An agent following this skill should be able to create an
independent Lore history at `refs/lore/<work-id>` without polluting the
working tree or source-branch history.

create-lore is the first skill in the git-lore workflow. Every other skill
(read-lore, edit-lore, sync-lore) assumes a Lore Work already exists.

## What we're trying to prove

A developer (or agent) can start a piece of work and immediately have a
durable, versioned context store that:

- survives branch rebases and renames
- does not clutter the repository working tree
- can be inspected with ordinary Git commands
- is associated with the current branch for convenient discovery

The test for create-lore specifically: given a new branch and a conversation
about what the work involves, an agent can create a Lore Work in under a
minute using only Git plumbing commands documented in the skill.

## Scope

Implement:

- `skills/create-lore/SKILL.md` with YAML frontmatter and step-by-step workflow
- Work ID resolution and validation
- Orphan lore commit creation via `git hash-object`, `git mktree`,
  `git commit-tree`, `git update-ref`
- Initial `plan.md` only (no other artefacts at creation time)
- Local branch association via `git config branch.<name>.lore <work-id>`
- Link to shared protocol reference at `skills/protocol.md`

Do not implement in this work:

- read-lore, edit-lore, sync-lore skills
- executable / CLI
- manifest or metadata files
- automatic work-id generation heuristics beyond branch-name default
- lore transport (sync-lore handles that separately)

## Initial document model

create-lore creates exactly one document: `plan.md`.

The plan should capture:

- what the work is trying to achieve
- what is in scope and out of scope
- immediate next steps

Additional artefacts (decisions.md, spec.md, questions.md, etc.) are created
later by edit-lore when the work produces durable conclusions. create-lore
should not preemptively create them.

## Working hypothesis

The hardest part of create-lore is not the Git plumbing — it is choosing a
good work-id and writing a useful initial plan.md from conversation context.

The skill must guide both:

1. **Work ID selection** — stable slug, lowercase, hyphens; defaults to
   sanitized branch name when not specified.
2. **Initial plan content** — concise, durable, not a transcript of the
   conversation that led to starting the work.

## Success criterion

A fresh agent with only the create-lore skill and no conversational history
should be able to:

1. Create a Lore Work on a new branch
2. Verify it exists with `git for-each-ref refs/lore` and
   `git show refs/lore/<work-id>:plan.md`
3. Confirm the Lore ref is independent of the branch's commit history
4. Confirm no lore files appear in the working tree

## Next steps

1. Write `skills/protocol.md` shared reference (work-id rules, git inspect
   commands, curation basics)
2. Write `skills/create-lore/SKILL.md` implementing this spec
3. Dogfood: associate `skills/create-lore` branch with this Lore Work
4. Verify with manual smoke test on a throwaway work-id in a temp clone

# Plan

## Objective

Implement the create-lore agent skill: the entry point for starting a new
Lore Work. An agent following this skill should be able to create an
independent Lore history at `refs/lore/<work-id>` without polluting the
working tree or source-branch history.

## Status

**Done.** `skills/create-lore/SKILL.md` and `skills/protocol.md` implemented
and smoke-tested.

## What was delivered

- `skills/create-lore/SKILL.md` — step-by-step workflow with git plumbing
- `skills/protocol.md` — distilled protocol reference for all skills
- Smoke test in temp clone confirms orphan commit, clean working tree,
  branch config with slash-containing branch names

## Remaining for this Work

- Commit skill files to source branch (`skills/create-lore`)
- sync-lore Lore Work (`refs/lore/sync-lore`) now planned — use sync-lore to push
  `refs/lore/create-lore` when a remote exists

## Success criterion

A fresh agent with only the create-lore skill and no conversational history
should be able to:

1. Create a Lore Work on a new branch
2. Verify it exists with `git for-each-ref refs/lore` and
   `git show refs/lore/<work-id>:plan.md`
3. Confirm the Lore ref is independent of the branch's commit history
4. Confirm no lore files appear in the working tree

All verified via smoke test.

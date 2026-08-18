# Decisions

## protocol.md lives in the working tree

**Decision:** Shared protocol reference at `skills/protocol.md`, distilled
from `refs/lore/git-lore`.

Project Lore (refs) holds full protocol context, decisions, and history.
`skills/protocol.md` is a concise working-tree reference agents can read
without git show commands. Skills link to both.

## Prefer semantic work-id over sanitized branch name

**Decision:** When the branch name is nested or long (e.g.
`skills/create-lore`), default to a shorter semantic id (`create-lore`)
when scope is clear, rather than the fully sanitized name
(`skills-create-lore`).

Confirm with the user when ambiguous. The sanitized form remains the
fallback when no semantic id is obvious.

## create-lore creates only plan.md at runtime

**Decision:** The skill creates exactly one document (`plan.md`) when
initialising a new Work.

This Work's own Lore contains `plan.md` and `spec.md` because documenting
the skill required both artefacts. That is dogfooding context, not the
default creation behaviour.

## Branch config key works with slash branch names

**Decision:** Use `git config branch.<branch-name>.lore <work-id>` including
when the branch name contains slashes (e.g.
`branch.skills/create-lore.lore`).

Verified: Git config handles this without escaping.

# Open Questions

## How much metadata is actually necessary?

The current protocol deliberately has no manifest.

Before adding machine-readable Work metadata, identify a concrete use
case that cannot be handled by:

- the Lore ref itself
- filenames
- Markdown
- YAML frontmatter

## How should Work discovery behave across worktrees?

The local branch → Work association remains unresolved.

This is likely to be an important part of the agent UX and should be
tested before introducing a more elaborate mechanism.

## How should Lore refs be transported?

Need to establish exact behaviour for:

- git clone
- git fetch
- git pull
- git push
- multiple remotes
- shallow clones
- worktree setups

The protocol must not assume arbitrary refs are transported automatically.

## What does "read Lore" actually mean?

**Partially resolved** by read-lore spec (`refs/lore/read-lore:spec.md`):

- Structured document order, not concatenation
- Handoff summary at read time (what / why / decided / rejected / next)
- Lore history preview (recent commits, not full diffs by default)

**Still open:**

- Default lore history depth (currently 10 commits — may be too shallow)
- Multi-work onboarding (read several Works in one session)
- Large document handling

## How much should agents automatically update Lore?

Too little:

    Lore becomes stale.

Too much:

    Lore becomes noisy.

**Partially addressed** by edit-lore spec (`refs/lore/edit-lore:spec.md`):

- Mandatory curation gate before any write
- Explicit edit-lore invocation (no automatic background updates)
- Batch meaningful curation, not per-turn activity

**Still open:** whether agents should proactively suggest edit-lore when durable conclusions emerge during implementation. See `refs/lore/edit-lore:questions.md`.

This is probably the most important UX question.

## Does the lack of a CLI become painful?

Do not predict this.

Dogfood the skills first and record concrete friction.

## Can a human understand Lore without the skills?

This is an important portability test.

A fresh clone containing Lore refs should remain understandable using
ordinary Git commands.

The skills should improve the experience, not be necessary to decode the
data.

## What is the right Work lifecycle?

We currently assume a Work can outlive a branch and potentially survive
the completion or deletion of its source branch.

Need to determine:

- when a Work is considered complete
- whether completed Lore should be archived
- whether Lore should ever be deleted
- how related or successor Works should be represented

Do not solve this until real usage provides evidence.

## Should handoff summaries be persisted?

read-lore generates summaries at read time. edit-lore can curate
handoff.md when a summary proves durable across sessions.

Open whether read-lore should suggest writing handoff.md automatically.
See `refs/lore/read-lore:questions.md`.

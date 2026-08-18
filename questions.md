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

Simply concatenating all Markdown is unlikely to be ideal.

A useful reader should distinguish:

- current state
- decisions
- unresolved questions
- historical context

Need to test this with real agents.

## How much should agents automatically update Lore?

Too little:

    Lore becomes stale.

Too much:

    Lore becomes noisy.

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

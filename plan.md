# Plan

## Objective

Create the smallest useful version of git-lore: a Git-native convention
for preserving the evolving context of software work, accessible to both
humans and AI agents.

The initial implementation should be a collection of agent skills rather
than an executable.

## What we're trying to prove

A developer can start a piece of work, have an AI agent maintain a small
set of durable artefacts describing the work, and later hand that work to
another human or AI agent with substantially less context reconstruction.

The key test is not whether we can store Markdown in Git. We obviously can.

The test is whether the resulting workflow is useful enough that people
actually want to maintain this context.

## Initial scope

Implement:

- create-lore
- read-lore
- edit-lore
- sync-lore

Define:

- Lore Work identity
- refs/lore/* storage
- Lore commit/history semantics
- branch/work association
- synchronisation semantics
- guidance for what belongs in Lore
- examples demonstrating useful Lore

Do not initially build:

- an executable
- an MCP server
- a database
- a GitHub application
- automatic transcript capture
- automatic AI summarisation
- GitHub-specific integrations
- a GUI/TUI

## Initial document model

A Lore Work contains ordinary files, primarily Markdown.

The initial creation flow should create:

- plan.md

Additional artefacts are created as needed.

Examples:

- spec.md
- decisions.md
- investigation.md
- architecture.md
- questions.md
- handoff.md

The protocol should not require a fixed taxonomy.

There is deliberately no required _lore.json or other manifest.

## Working hypothesis

Lore is most useful when it records durable conclusions rather than
transient activity.

A useful heuristic:

> Would a future developer or agent otherwise have to rediscover this?

If yes, it probably belongs in Lore.

## First experiment

Use git-lore itself as the first real project.

Record:

- why the project exists
- alternatives considered
- protocol decisions
- implementation discoveries
- failed approaches
- changes to the product model

Do not manufacture Lore merely to demonstrate the format.

The project should expose its own UX problems.

## Success criterion

After using git-lore to develop git-lore, hand the work to a fresh agent
with no conversational history.

The agent should be able to explain:

1. what git-lore is;
2. why it exists;
3. what has been decided;
4. what has been rejected and why;
5. what currently needs to happen next.

If this is not materially easier with Lore than without it, reconsider
the project.

## Implementation status

| Skill | Status | Lore Work |
|-------|--------|-----------|
| create-lore | Done | `refs/lore/create-lore` |
| read-lore | Done — `skills/read-lore/SKILL.md` | `refs/lore/read-lore` |
| edit-lore | Done — `skills/edit-lore/SKILL.md` | `refs/lore/edit-lore` |
| sync-lore | Done — `skills/sync-lore/SKILL.md` | `refs/lore/sync-lore` |
| protocol reference | Done | `skills/protocol.md` |

Per-skill Lore Works (`refs/lore/create-lore`, `refs/lore/read-lore`) hold
plan and spec for each skill. This Lore Work (`refs/lore/git-lore`) remains
the protocol authority.

## Next steps

All four core skills (create-lore, read-lore, edit-lore, sync-lore) are done.

Remaining:

1. Commit completed skills to the source branch
2. Dogfood transport when a remote exists (push/fetch lore refs)
3. README done (see refs/lore/readme); examples remain

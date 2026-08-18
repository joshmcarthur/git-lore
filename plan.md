# Plan

## Objective

Implement the read-lore agent skill: the primary way agents and humans
onboard to an existing Lore Work. read-lore resolves which Work to read,
presents documents in a useful order, and synthesizes a handoff summary so
a fresh agent can understand the work without conversational history.

## Status

**In progress.** Lore Work created; `skills/read-lore/SKILL.md` not yet
written.

## What we're trying to prove

The git-lore success criterion depends on read-lore:

> A fresh agent should explain what the project is, why it exists, what
> has been decided, what was rejected, and what needs to happen next.

read-lore must make that materially easier than manually running `git show`
on every lore file. The test is not whether we can read Markdown from a
Git ref — it is whether the **presentation** reduces context reconstruction.

## Scope

Implement:

- `skills/read-lore/SKILL.md` with YAML frontmatter and step-by-step workflow
- Work-id resolution (explicit argument → branch config → list and ask)
- Structured document read order (not raw concatenation)
- Lore history preview (`git log` — recent commits only)
- Handoff summary synthesis (what / why / decided / rejected / next)
- Support reading a specific document or all documents
- Link to `skills/protocol.md` and Lore Work `refs/lore/git-lore`

Do not implement in this work:

- edit-lore, sync-lore skills
- Automatic lore summarisation via AI beyond structured presentation
- Reading lore from remote without local ref (sync-lore first)
- Multi-work aggregation (read several Works in one invocation)
- Executable or helper scripts

## Working hypothesis

Simply concatenating all lore markdown is insufficient. A useful reader
must distinguish:

- current direction (plan.md)
- settled choices (decisions.md)
- unresolved items (questions.md)
- technical detail (spec.md, investigation.md)
- handoff context (handoff.md)

The agent synthesizes a brief handoff summary **after** reading structured
content — the summary is generated at read time, not stored as a transcript
in Lore (unless edit-lore later curates a handoff.md).

## Success criterion

A fresh agent invoked with only read-lore on `refs/lore/git-lore` should
produce a handoff summary covering all five git-lore success-criterion
points without reading conversation history.

Verification:

```bash
git for-each-ref refs/lore --format='%(refname:short)'
git show refs/lore/read-lore:plan.md
```

Manual smoke test: run read-lore workflow against `git-lore` and
`create-lore` Works; confirm structured output, not a wall of markdown.

## Next steps

1. Write `skills/read-lore/SKILL.md` per spec.md in this Lore Work
2. Dogfood: read `refs/lore/git-lore` using the skill workflow
3. Update `refs/lore/git-lore` plan.md status when read-lore is done
4. Associate `skills/read-lore` branch via branch config

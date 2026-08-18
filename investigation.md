# Investigation

## Existing AI memory systems

AI coding environments increasingly provide persistent memory and project
instructions.

These solve an adjacent problem but generally have an agent/product
orientation rather than being a portable shared history of a specific
piece of work.

The relevant distinction is:

    AI memory:
        What should this agent remember?

    Lore:
        What should this work remember?

Lore must therefore remain useful if the project changes AI providers.

## Existing repository documentation

Ordinary Markdown in the repository is the simplest competing solution.

The key question is therefore not whether Lore can store Markdown.

It is whether work-scoped storage outside the working tree provides enough
benefit to justify another convention.

Potential advantages:

- no repository clutter
- independent history
- no merge conflicts with source documentation
- survives branch rewriting
- work can have its own lifecycle
- shared through Git
- can persist after source branch deletion

These advantages need to be demonstrated in actual use.

## Git Notes

Git Notes are attractive because they already provide a Git-native
secondary information channel.

However, Notes are attached to Git objects, while Lore is attached to a
Work.

This distinction appears fundamental.

## Git garbage collection

A Lore ref is a normal Git ref and therefore keeps its reachable commit,
tree and blob objects alive.

This means Lore does not need special storage or garbage-collection
protection.

The important failure mode is transport rather than GC: arbitrary
refs/lore/* refs need explicit fetch/push handling.

## GitHub and hosted Git

Lore should not depend on GitHub-specific functionality.

The underlying model should work with any Git remote capable of
transporting the Lore refs.

GitHub integration may be useful later, but it should not be part of the
core protocol.

## Agent skills versus executable

The first implementation should be skills rather than a CLI.

The skills are primarily instructions for an agent to perform ordinary
Git operations.

This makes the initial implementation extremely small and lets real use
identify where a dedicated executable would provide value.

## Key UX risk

The protocol itself is easy.

The difficult question is whether agents will reliably recognise which
information deserves to become durable Lore.

If agents write too little, Lore becomes stale.

If they write too much, Lore becomes a transcript or bureaucratic log.

The skills therefore need strong guidance and examples around curation.

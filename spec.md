# Specification

## Core model

Lore is a persistent, versioned body of context associated with a piece
of software work.

The fundamental object is a Work, not a Git commit or branch.

A Work may span multiple branches during its lifetime.

## Storage

Each Work is identified by a Git ref:

    refs/lore/<work-id>

The ref points to the latest commit in an independent Git history.

The Lore tree contains ordinary files, primarily Markdown.

There is no required manifest or metadata file.

Example:

    plan.md
    decisions.md
    investigation.md

## Ref ownership

git-lore does not own or manage the refs/lore namespace as a database.

refs/lore/* is a Git convention.

The protocol defines how agents and humans can interpret and manipulate
those refs, but the underlying objects remain ordinary Git objects.

Lore tooling should use standard Git operations wherever practical.

## Why independent history?

Source history and reasoning history have different lifecycles.

Source commits may be:

- rebased
- squashed
- amended
- cherry-picked
- deleted

Lore should not be rewritten by these operations.

The Lore ref therefore must not depend on source commit ancestry.

## Garbage collection

A live Lore Work is kept alive by its refs/lore/<work-id> ref.

Lore must never rely on dangling objects or reflogs for persistence.

Deleting the Lore ref is the semantic operation for removing the Work.
Once unreachable, its objects may eventually be garbage collected by Git.

## Transport

Lore refs must be synchronisable through ordinary Git transport.

The protocol must support pushing and fetching:

    refs/lore/*

The exact refspec behaviour must be documented and tested.

The protocol must not require a proprietary service.

## Branch association

A local mapping may associate a current branch with a Lore Work.

This association is local metadata and should not itself be committed
to the repository.

The shared identity is the Lore ref.

## Documents

Lore documents are ordinary files in the Lore tree.

Markdown is the preferred format because it is:

- human-readable
- AI-readable
- Git-friendly
- portable
- easy to inspect without Lore tooling

There is no mandatory document taxonomy.

A Work may contain:

    plan.md
    spec.md
    decisions.md
    investigation.md
    questions.md
    handoff.md

or any other useful structure.

## Document metadata

Documents MAY use YAML frontmatter when metadata is useful.

For example:

    ---
    title: Topology generation investigation
    status: active
    ---

Metadata should only be introduced where it serves a concrete purpose.

The Work ID does not need to be repeated in documents because it is
already represented by the refs/lore/<work-id> ref.

## Curation

Lore is not a transcript.

It should preserve:

- decisions
- constraints
- discoveries
- rejected approaches
- important investigations
- current plans
- unresolved questions
- handoff information

It should not preserve routine conversation or activity telemetry.

## Agent independence

Lore must be usable by different AI agents and by humans.

No Lore data should require a particular AI provider to interpret.

## Git compatibility

Lore should remain useful without git-lore tooling.

A user should be able to inspect Lore with ordinary Git commands such as:

    git log refs/lore/<work-id>
    git show refs/lore/<work-id>:plan.md
    git diff refs/lore/<work-id>~1 refs/lore/<work-id>

git-lore provides conventions and agent workflows rather than hiding
the underlying Git representation.

## Future tooling

A CLI may eventually provide a convenient interface to this protocol.

Such a CLI is explicitly not part of the initial implementation.

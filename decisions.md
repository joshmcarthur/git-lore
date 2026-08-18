# Decisions

## Use Git refs rather than Git Notes

**Decision:** Lore will use independent refs under `refs/lore/*`.

Git Notes are designed to attach additional information to existing Git
objects. Lore instead represents an independent history belonging to a
piece of work.

This allows Lore to survive rebases, squashes and deletion of the source
branch.

## No executable initially

**Decision:** Implement the first version as agent skills operating
directly on Git.

The proposed executable would initially be a thin wrapper around Git
plumbing.

Building it before proving the workflow risks creating a product around
an unvalidated interaction model.

If repeated friction appears, those pain points can justify a future CLI.

## Lore is not AI memory

**Decision:** Position Lore as shared project/work context rather than
agent-specific memory.

AI products may maintain private or provider-specific memories. Lore is
intended to remain:

- portable
- shared
- Git-native
- explicit
- reviewable
- durable

## No transcript capture

**Decision:** Do not store full AI conversations.

Transcripts contain large amounts of low-value and ephemeral information.

Lore should preserve conclusions and useful reasoning extracted from those
conversations.

## Work rather than branch

**Decision:** The primary identity is a Work rather than a branch.

Branches can be rebased, renamed, split, merged or deleted.

The reasoning associated with the work should have a longer and more
stable identity.

## No fixed document taxonomy

**Decision:** Do not require plan.md, spec.md, decisions.md, etc.

Only plan.md is created initially.

Other artefacts should emerge according to the needs of the work.

This avoids turning Lore into documentation bureaucracy.

## No Lore manifest

**Decision:** Do not require _lore.json or another manifest in a Lore
tree.

The Lore ref already provides the Work identity. A manifest would add
structure and apparent ownership without solving a demonstrated problem.

Documents may use YAML frontmatter when machine-readable metadata is
actually useful.

We can introduce additional metadata structures later if a concrete use
case requires them.

## git-lore does not own the refs/lore namespace

**Decision:** Treat refs/lore/* as a convention rather than an owned
database namespace.

The underlying Git refs and objects remain ordinary Git primitives.

This keeps the protocol inspectable with standard Git tooling and avoids
creating an artificial distinction between "Lore data" and Git data.

## Markdown is the primary interchange format

**Decision:** Lore documents should normally be Markdown.

This maximises portability between humans, AI agents and development
tools.

Metadata should be expressed using YAML frontmatter where necessary,
rather than introducing a separate metadata database or manifest.

## Garbage collection

**Decision:** Every live Lore Work must be reachable from refs/lore/*.

This guarantees that normal Git garbage collection preserves Lore
history.

Dangling Lore objects are not a supported persistence mechanism.

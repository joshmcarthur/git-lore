# Open Questions

## When should agents proactively suggest edit-lore?

edit-lore is invoked explicitly today. Open whether agents should suggest
"record this in Lore?" when durable conclusions emerge during implementation.

Risk of noise if too proactive; risk of stale Lore if too passive. The
curation gate helps but does not solve the trigger problem.

## Lore commit granularity in practice

Decision: batch related curation. Real usage may reveal cases where
separate commits per decision are preferable for history readability.

Monitor during dogfooding.

## File deletion in lore

edit-lore supports deletion via `git add -A` in temp dir, but guidance says
prefer recording superseding decisions over removing old ones.

Open: whether explicit deletion workflow or guidance is needed when Lore
content is genuinely wrong (not merely outdated).

## Diverged lore refs (partially resolved)

edit-lore still does not merge diverged histories automatically. sync-lore
now documents how to handle divergence:

- **Detection** — status mode compares local and remote lore refs
- **Side-ref workflow** — remote tip held at `refs/lore/<id>-remote` while
  local stays at `refs/lore/<id>`
- **Reconciliation** — manual: inspect both histories, curate merged outcome
  with edit-lore, then push

See `refs/lore/sync-lore:spec.md` and `refs/lore/sync-lore:decisions.md`.
Open: whether sync-lore should grow assisted merge tooling beyond documenting
the manual path.

## handoff.md vs read-time summaries

read-lore generates handoff summaries at read time. edit-lore can curate
handoff.md for durable cross-session context.

Open: criteria for when a read-time summary should become handoff.md.

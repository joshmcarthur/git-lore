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

## Diverged lore refs

If local and remote lore histories diverge, edit-lore does not merge them.
sync-lore era problem. Record friction when encountered; do not solve in
edit-lore spec until sync-lore exists.

## handoff.md vs read-time summaries

read-lore generates handoff summaries at read time. edit-lore can curate
handoff.md for durable cross-session context.

Open: criteria for when a read-time summary should become handoff.md.

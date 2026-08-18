# Open Questions

## How much lore history should read-lore show by default?

Current default: 10 commits via `git log --oneline`. Full history and
per-commit diffs are available on request.

Risk: shallow history hides important earlier decisions that were never
copied into decisions.md. Monitor during dogfooding.

## Should read-lore aggregate multiple Works?

Example: onboarding to git-lore might need both `refs/lore/git-lore` and
`refs/lore/create-lore`.

Current decision: single Work per invocation. User or agent reads Works
sequentially. Revisit if multi-work onboarding is a common pain point.

## How to handle large lore documents?

No size limits yet. If a spec.md or investigation.md grows very large,
read-lore may need truncation guidance or "read section X" support.

Do not solve until observed in practice.

## Should handoff summaries be persisted?

read-lore generates handoff summaries at read time. edit-lore can curate
a durable `handoff.md` when a summary proves useful across sessions.

Open: whether read-lore should ever suggest writing handoff.md automatically.

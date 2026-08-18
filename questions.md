# Open Questions

## Should protocol.md duplicate refs/lore/git-lore?

Currently `skills/protocol.md` distills the protocol for quick agent
access. `refs/lore/git-lore` remains authoritative for decisions,
investigations, and history.

Risk: drift between working-tree reference and Lore. Mitigation for now:
protocol.md is minimal and links to the Lore ref. Revisit if duplication
becomes painful.

## Default work-id when branch naming is inconsistent

Semantic id vs sanitized branch name is a judgment call. Real usage may
reveal a better default heuristic or a need to always ask.

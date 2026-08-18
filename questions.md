# Open questions

## When (if ever) should extensions graduate?

If the `git-lore` CLI becomes widely used, should it stay under `extensions/` indefinitely, move to a separate repository, or join a broader tooling package? Prefer evidence from dogfooding over premature promotion to the root.

## CLI vs skills ownership of writes

When `create` / `edit` / `sync` subcommands land, should they enforce the same curation gate as edit-lore (and how — prompts, dry-run, required message flags)? Risk: a convenient CLI becomes a transcript dump without skill discipline.

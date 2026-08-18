# Decisions

## Git commands for Lore attribution, not GitHub blob URLs

**Decision:** README cites Lore excerpts with `git show refs/lore/<id>:<file>` commands.

GitHub does not provide stable browse URLs for custom refs under `refs/lore/*`.
Skill files link to normal paths in the tree; Lore refs use inspect commands.

## Dogfooding narrative is primary structure

**Decision:** "See it in action" is the central README section, not a sidebar example.

The repository's strongest evidence is the chain bootstrap → create-lore → read-lore
→ edit-lore → sync-lore with matching Lore Works. The README leads with that story.

## README Lore Work separate from protocol Lore

**Decision:** README authoring tracked in `refs/lore/readme`, not only in
`refs/lore/git-lore`.

Keeps protocol Lore focused on git-lore itself; readme Lore holds README-specific
planning and decisions without bloating the protocol Work.

# Open Questions

## Shallow clones

Does shallow clone (`--depth 1`) affect lore ref transport or history
inspection? Lore histories are typically short but may reference old
commits for `git diff` and `git log`. Test before documenting.

## Multiple remotes

When a repository has multiple remotes (origin, upstream), lore refspec
must be configured per remote. Is per-remote lore sync confusing? May need
clearer status reporting showing which remote was synced.

## Side ref cleanup

`refs/lore/<id>-remote` side refs accumulate during divergence resolution.
Should sync-lore always clean up, or leave them for inspection? Current
approach: suggest cleanup after reconciliation.

## Force-push policy

Default is no force-push. Is there a legitimate case for force-pushing lore
(e.g. accidental bad lore commit)? If so, require explicit user confirmation
and document risks. Do not add until real usage demonstrates need.

## git-lore as CI check

Could CI verify lore refs are pushed? Out of scope for sync-lore skill but
transport semantics may affect future tooling.

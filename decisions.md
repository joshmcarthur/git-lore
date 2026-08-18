# Decisions

## Fetch refspec without force prefix

**Decision:** Configure `refs/lore/*:refs/lore/*` without the `+` prefix
on fetch refspecs.

A `+` prefix would force-overwrite local lore refs on fetch, destroying
unpushed local lore commits. Lore histories are independently curated;
local unpushed curation must not be silently discarded by fetch.

## No automatic lore merge

**Decision:** sync-lore detects divergence but does not automatically merge
lore histories.

Lore documents are curated Markdown. Automatic merge produces conflict
markers and merged noise. Reconciliation is a curation task — read both
sides, reconcile via edit-lore.

Divergence workflow: fetch to side ref (`refs/lore/<id>-remote`), compare
histories, manual reconciliation, push.

## Transport is explicit, not ambient

**Decision:** Lore refs require explicit fetch/push or sync-lore invocation.
`git pull` alone does not sync lore refs to local even after refspec
configuration — users must `git fetch` (which sync-lore fetch mode runs).

This matches Git's default behaviour and keeps lore transport visible and
intentional.

## Refspec is per-remote local config

**Decision:** Fetch refspec lives in `remote.<name>.fetch` git config,
local to the clone, not committed to the repository.

Same pattern as branch.lore association. Each clone and each remote needs
one-time setup. Document in post-clone workflow.

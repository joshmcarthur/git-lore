# Specification

## Purpose

sync-lore is the agent skill that transports Lore Works between the local
repository and Git remotes using ordinary `git fetch` and `git push` on
`refs/lore/*`. It configures fetch refspecs (which are not set by default),
executes transport operations, and reports status.

## Prerequisites

- A Git remote exists (default: `origin`). If no remote, explain lore is
  local-only and stop.
- Lore Works exist locally and/or on remote.

## Inputs

| Input | Required | Default |
|-------|----------|---------|
| Remote name | No | `origin` |
| Mode | No | `sync` (fetch then push) |
| Work id | No | All lore refs |

### Modes

| Mode | Trigger terms | Behaviour |
|------|---------------|-----------|
| sync | "sync lore", default | Fetch then push lore refs |
| fetch | "fetch lore", "pull lore" | Fetch lore refs only |
| push | "push lore" | Push lore refs only |
| status | "lore status", "lore remote" | Compare local vs remote refs |

When user specifies a work-id, scope fetch/push/status to that Work only.

## Outputs

| Output | Description |
|--------|-------------|
| Transport report | What was fetched, pushed, or already up to date |
| Refspec confirmation | Whether fetch refspec is configured |
| Divergence warning | When local and remote lore histories disagree |

No files are written to the repository working tree.

## Fetch refspec configuration

Lore refs are **not** fetched or pushed by default. Configure once per
remote:

```bash
REMOTE=origin
git config --add "remote.${REMOTE}.fetch" 'refs/lore/*:refs/lore/*'
```

Properties:

- **No `+` prefix** on fetch refspec — avoids force-overwriting local lore
  refs that have unpushed commits when fetching
- Must be repeated for each remote that should transport lore
- Persists in `.git/config` (local, not committed)
- After configuration, `git fetch <remote>` includes lore refs alongside
  branch refs

### Post-clone setup

A fresh `git clone` does not include lore refs until configured and fetched:

```bash
git clone <url> <dir>
cd <dir>
git config --add remote.origin.fetch 'refs/lore/*:refs/lore/*'
git fetch origin
git for-each-ref refs/lore
```

`git pull` updates the current branch but does not replace explicit lore
fetch/push. Users who want lore after pull should run `git fetch origin`
(with refspec configured) or use sync-lore fetch mode.

## Fetch

```bash
REMOTE=origin
git fetch "$REMOTE" 'refs/lore/*:refs/lore/*'
```

For a single Work:

```bash
WORK_ID=<work-id>
git fetch "$REMOTE" "refs/lore/${WORK_ID}:refs/lore/${WORK_ID}"
```

Fetch without `+` will:

- Fast-forward local lore ref if remote is ahead
- Leave local ref unchanged if local is ahead (no overwrite)
- Fail or skip if histories have diverged (depending on Git version/config)

After fetch, suggest read-lore to review newly fetched Works.

## Push

Push all local lore refs:

```bash
REMOTE=origin
git push "$REMOTE" 'refs/lore/*'
```

Push a single Work:

```bash
WORK_ID=<work-id>
git push "$REMOTE" "refs/lore/${WORK_ID}"
```

If push is rejected (non-fast-forward), local and remote lore histories
have diverged. Do not force-push lore refs by default. See divergence
handling below.

## Sync mode

Default workflow:

```bash
REMOTE=origin
# 1. Ensure refspec configured (see above)
# 2. Fetch
git fetch "$REMOTE" 'refs/lore/*:refs/lore/*'
# 3. Push
git push "$REMOTE" 'refs/lore/*'
```

Report results of each step. If push fails after successful fetch, report
divergence.

## Status mode

Compare local and remote lore refs:

```bash
REMOTE=origin
# Local refs
git for-each-ref refs/lore --format='%(refname:short) %(objectname:short)'

# Remote refs
git ls-remote "$REMOTE" 'refs/lore/*'
```

For each work-id present on either side, report:

| State | Meaning |
|-------|---------|
| synced | Same commit SHA |
| local-only | Ref exists locally, not on remote |
| remote-only | Ref exists on remote, not locally |
| diverged | Both exist, different SHAs, neither is ancestor of other |
| ahead | Local has commits not on remote (FF push possible) |
| behind | Remote has commits not on local (FF fetch possible) |

Use `git rev-parse` and `git merge-base` to classify.

## Divergence handling

sync-lore does **not** automatically merge diverged lore histories. Lore
documents are curated Markdown — automatic merge produces conflict markers
that violate curation principles.

When push is rejected or status shows divergence:

1. Fetch remote ref to a side ref without overwriting local:

```bash
WORK_ID=<work-id>
REMOTE=origin
git fetch "$REMOTE" "refs/lore/${WORK_ID}:refs/lore/${WORK_ID}-remote"
```

2. Show both histories:

```bash
git log --oneline refs/lore/${WORK_ID}
git log --oneline refs/lore/${WORK_ID}-remote
```

3. Suggest manual reconciliation:
   - read-lore for local version
   - `git show refs/lore/${WORK_ID}-remote:plan.md` (etc.) for remote
   - edit-lore to curate a reconciled version
   - sync-lore push after reconciliation
4. Clean up side ref when done:

```bash
git update-ref -d "refs/lore/${WORK_ID}-remote"
```

Record persistent friction in Lore `questions.md` via edit-lore if the
workflow is painful.

## SKILL.md structure

The skill file at `skills/sync-lore/SKILL.md` must include:

- YAML frontmatter: `name`, `description` (third person, trigger terms)
- Step-by-step workflow matching this spec
- Checklist pattern consistent with other skills
- Post-operation verification (`git status` clean)
- Links to protocol.md and other skills
- Explicit non-goals

Trigger terms: "sync lore", "push lore", "fetch lore", "pull lore",
"lore status", "lore remote", "share lore".

## Verification

After any sync operation:

```bash
git status
git for-each-ref refs/lore --format='%(refname:short) %(objectname:short)'
```

Confirm working tree is clean. sync-lore never checks out lore files.

## Relationship to other skills

| Skill | Relationship |
|-------|--------------|
| create-lore | Creates refs that sync-lore transports |
| read-lore | Suggested after fetch; reads fetched Works |
| edit-lore | Suggested for divergence reconciliation |
| protocol.md | Shared transport summary |

sync-lore does not call other skills. It suggests them when appropriate.

## Explicit non-goals

- No force-push of lore refs (default)
- No automatic lore history merge
- No GitHub-specific ref handling
- No executable or helper scripts
- No lore ref deletion
- No working-tree checkout of lore files
- No transport of source branches (ordinary git push/pull handles those)

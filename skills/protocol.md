# Lore protocol reference

Canonical copy: `skills/protocol.md`. Each skill links to it at
`references/protocol.md` (symlink). The skills CLI dereferences symlinks on
install, so individual skill installs still get a real copy.

Distilled from the git-lore specification (`refs/lore/git-lore`). Read that
Lore Work (`git show refs/lore/git-lore:spec.md`) for full protocol context,
decisions, and open questions.

## Work identity

A Lore **Work** is identified by a Git ref:

    refs/lore/<work-id>

- `work-id`: lowercase slug — letters, digits, hyphens
- Examples: `git-lore`, `create-lore`, `auth-refactor`
- The ref points to the latest commit in an **independent** Git history
- Lore history must not parent off source-branch commits

## Lore tree

Ordinary files, primarily Markdown. No manifest (`_lore.json` etc.).

Common artefacts (none required except what create-lore starts with):

| File | Typical role |
|------|-------------|
| `plan.md` | Current direction and next steps |
| `spec.md` | Technical specification for the work |
| `decisions.md` | Settled choices |
| `questions.md` | Unresolved questions |
| `investigation.md` | Research and alternatives |
| `handoff.md` | Context for the next person or agent |

Documents may use YAML frontmatter when machine-readable metadata is useful.
The work-id does not need to be repeated in documents.

## Lore commits

- Initial commit: **orphan** (no parents) — create-lore uses `git hash-object` +
  `git mktree` + `git commit-tree`
- Message prefix: `lore:` (e.g. `lore: initialise Work create-lore`)
- Subsequent edits append with normal parent linkage — edit-lore parents off
  current lore HEAD
- Lore commits are separate from source-code commits on the branch

## Editing lore (edit-lore)

edit-lore exports via `git archive`, edits in a temp directory, and commits
with an **isolated index** (`GIT_INDEX_FILE`). Without index isolation,
`git add -A` in the temp dir pollutes the main repository index.

```bash
REF="refs/lore/<work-id>"
COMMIT=$(git rev-parse "$REF")
DIR=$(mktemp -d)
INDEX=$(mktemp)
GIT_DIR=$(git rev-parse --git-dir)

git archive "$COMMIT" | tar -x -C "$DIR"
# edit files in $DIR

export GIT_INDEX_FILE="$INDEX"
git --git-dir="$GIT_DIR" read-tree --empty
cd "$DIR"
git --git-dir="$GIT_DIR" add -A .
TREE=$(git --git-dir="$GIT_DIR" write-tree)
unset GIT_INDEX_FILE

NEW=$(git commit-tree "$TREE" -p "$COMMIT" -m "lore: <what changed>")
git update-ref "$REF" "$NEW"
rm -rf "$DIR" "$INDEX"
```

Apply the curation gate before writing. Read current lore before editing.
See `refs/lore/edit-lore:spec.md` for full edit-lore specification.

## Branch association (local)

Associate the current branch with a Work for discovery:

    git config branch.<branch-name>.lore <work-id>

- Local only — not committed to the repository
- Read by read-lore and edit-lore when no work-id is given
- Works with branch names containing slashes (e.g. `branch.skills/create-lore.lore`)

## Curation

Lore is not a transcript. Preserve durable conclusions:

- decisions, constraints, discoveries, rejected approaches
- current plans, unresolved questions, handoff information

Heuristic:

> Would a future developer or agent otherwise have to rediscover this?

Do not preserve routine conversation or activity telemetry.

## Inspect without skills

```bash
git for-each-ref refs/lore
git log refs/lore/<work-id>
git show refs/lore/<work-id>:plan.md
git diff refs/lore/<work-id>~1 refs/lore/<work-id>
```

## Transport

Lore refs are not fetched or pushed by default. Configure once per remote
(local `.git/config`, not committed):

```bash
git config --add remote.origin.fetch 'refs/lore/*:refs/lore/*'
```

No `+` prefix on the fetch refspec — avoids force-overwriting local lore
refs with unpushed commits on fetch.

```bash
git fetch origin 'refs/lore/*:refs/lore/*'
git push origin 'refs/lore/*'
```

Post-clone: add the refspec, then `git fetch origin` before expecting lore
locally.

Diverged histories are not merged automatically. sync-lore documents detection,
side refs (`refs/lore/<id>-remote`), and manual reconciliation via edit-lore.
See `refs/lore/sync-lore:spec.md` and `refs/lore/sync-lore:decisions.md`;
agent skill [skills/sync-lore/SKILL.md](sync-lore/SKILL.md).

Lore must remain reachable via `refs/lore/*` for Git GC to preserve it.

## Garbage collection

Every live Work must be reachable from `refs/lore/*`. Deleting the ref is
the semantic operation for removing a Work.

# Lore protocol reference

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

- Initial commit: **orphan** (no parents)
- Message prefix: `lore:` (e.g. `lore: initialise Work create-lore`)
- Subsequent edits append with normal parent linkage
- Lore commits are separate from source-code commits on the branch

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

Lore refs are not fetched or pushed by default. See sync-lore for refspec
setup. Lore must remain reachable via `refs/lore/*` for Git GC to preserve it.

## Garbage collection

Every live Work must be reachable from `refs/lore/*`. Deleting the ref is
the semantic operation for removing a Work.

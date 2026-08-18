# Specification

## Purpose

edit-lore is the agent skill that curates an existing Lore Work. It exports
the lore tree to a temporary directory, applies edits, commits a new lore
revision with parent linkage, and updates `refs/lore/<work-id>`.

edit-lore never writes lore files to the repository working tree.

## Inputs

| Input | Required | Default |
|-------|----------|---------|
| `work-id` | No | Branch config, else list and ask |
| Edit intent | Yes | Agent derives from conversation — what durable content to record |
| Target document(s) | No | Agent selects based on edit intent |

### Work-id resolution

Same as read-lore — resolve in order:

1. Explicit argument
2. `git config --get branch.<current-branch>.lore`
3. List `git for-each-ref refs/lore` and ask

If Work does not exist → suggest create-lore.
If Work may exist only on remote → suggest sync-lore first.

## Outputs

| Output | Description |
|--------|-------------|
| Updated lore ref | `refs/lore/<work-id>` advanced to new commit |
| Lore commit | Child of previous lore HEAD with `lore: <what changed>` message |
| Diff summary | `git diff <old> <new>` or `git log -1` for user review |

No files written to the repository working tree.

## Edit modes

| Mode | Trigger | Behaviour |
|------|---------|-----------|
| Curate | Default | Apply curation gate, edit appropriate document(s) |
| Update plan | Plan or scope changed | Edit `plan.md` |
| Record decision | Choice settled | Add to or create `decisions.md` |
| Record question | Uncertainty identified | Add to or create `questions.md` |
| Add investigation | Research worth preserving | Add to or create `investigation.md` |
| Curate handoff | Durable handoff needed | Create or update `handoff.md` |
| Create artefact | New document warranted | Add file to lore tree |

Modes are guidance for the agent, not separate code paths. One invocation
may combine modes (e.g. record decision + update plan next steps) in a
single lore commit.

## Curation gate (mandatory)

Before writing, the agent must pass the curation gate:

1. **Durable?** Apply heuristic from protocol.md:

   > Would a future developer or agent otherwise have to rediscover this?

2. **Category?** Lore should preserve:
   - decisions, constraints, discoveries, rejected approaches
   - current plans, unresolved questions, handoff information

3. **Exclude?** Do not record:
   - conversation transcripts or chat summaries
   - routine implementation activity ("fixed typo", "ran tests")
   - detail that belongs in source code, commit messages, or PR descriptions
   - duplicate of content already in Lore

If content fails the gate, do not edit Lore. Tell the user why.

## Read before write

Before editing, inspect current lore state:

- Prefer read-lore for full context on complex edits
- Minimum: `git ls-tree -r --name-only refs/lore/<id>` and
  `git show refs/lore/<id>:<target-file>` for files being modified

Never overwrite a document without reading its current content.

## Document emergence rules

Do not preemptively create documents. Create a new artefact only when the
work produces content that fits:

| File | Create when |
|------|-------------|
| `decisions.md` | First settled choice needs recording |
| `questions.md` | First unresolved question needs recording |
| `investigation.md` | Research or alternatives worth preserving |
| `spec.md` | Technical specification emerges for the work |
| `architecture.md` | System design decisions need recording |
| `handoff.md` | Durable handoff context needed across sessions |

`plan.md` always exists (created by create-lore). Other files are optional.

### decisions.md entry format

```markdown
## Short title

**Decision:** [What was decided]

[Why — durable reasoning, not conversation recap]
```

### questions.md entry format

```markdown
## Short title

[Question and any context needed to understand it]
```

## Lore commit semantics

Unlike create-lore (orphan), edit-lore commits **parent off the current
lore HEAD**:

```bash
COMMIT=$(git rev-parse "refs/lore/${WORK_ID}")
# ... build new tree ...
NEW=$(git commit-tree "$TREE" -p "$COMMIT" -m "lore: <what changed>")
git update-ref "refs/lore/${WORK_ID}" "$NEW"
```

Message format:

    lore: <concise description of what was curated>

Examples:

- `lore: record decision on temp work-tree for edits`
- `lore: update plan status — read-lore done`
- `lore: add open question on lore merge conflicts`

One lore commit per meaningful curation event. Batch related changes;
do not commit per-line edits.

## Git workflow (temp work-tree)

edit-lore uses a temporary directory and **isolated index** to avoid touching
the project working tree or the main Git index:

```bash
WORK_ID=<work-id>
REF="refs/lore/${WORK_ID}"
COMMIT=$(git rev-parse "$REF")
DIR=$(mktemp -d)
INDEX=$(mktemp)
GIT_DIR=$(git rev-parse --git-dir)

# 1. Export current lore tree
git archive "$COMMIT" | tar -x -C "$DIR"

# 2. Edit files in $DIR (agent applies curated changes)

# 3. Build new tree using isolated index
export GIT_INDEX_FILE="$INDEX"
cd "$DIR"
git --git-dir="$GIT_DIR" add -A .
TREE=$(git --git-dir="$GIT_DIR" write-tree)
unset GIT_INDEX_FILE

# 4. Commit and update ref
NEW=$(git commit-tree "$TREE" -p "$COMMIT" -m "lore: <what changed>")
git update-ref "$REF" "$NEW"

# 5. Cleanup
rm -rf "$DIR" "$INDEX"
```

Notes:

- `GIT_INDEX_FILE` isolation is **required** — without it, `git add -A` in the
  temp dir pollutes the main repository index
- `git archive` exports without checkout, avoiding index side effects
- `git add -A` in the temp dir stages creates, modifications, and deletions
- File deletion in lore is supported but rare — prefer recording superseding
  decisions over removing history
- The lore ref must always point to a commit; do not create an empty tree

## Workflow checklist

```
- [ ] Step 1: Resolve work-id
- [ ] Step 2: Verify Work exists
- [ ] Step 3: Curation gate — should this be recorded?
- [ ] Step 4: Read current lore state
- [ ] Step 5: Export to temp dir and apply edits
- [ ] Step 6: Commit lore revision
- [ ] Step 7: Verify, show diff, and report
- [ ] Step 8: Suggest read-lore or sync-lore if appropriate
```

### Step 7 verification

```bash
git log --oneline -3 "refs/lore/${WORK_ID}"
git diff "${COMMIT}" "${NEW}"
git status
git ls-tree -r --name-only "refs/lore/${WORK_ID}"
```

Confirm:

- Lore ref advanced to new commit
- Working tree clean
- Diff reflects intended curation only

### Step 8 suggestions

| Situation | Suggest |
|-----------|---------|
| User wants to review changes | read-lore |
| Lore should be shared | sync-lore |
| No Work exists | create-lore |

## SKILL.md structure

The skill file at `skills/edit-lore/SKILL.md` must include:

- YAML frontmatter: `name`, `description` (third person, trigger terms)
- Curation gate (prominent — this is the primary UX control)
- Read-before-write requirement
- Document emergence rules
- Temp work-tree workflow with full shell commands
- Edit modes table
- Entry formats for decisions.md and questions.md
- Link to `skills/protocol.md` and `refs/lore/git-lore`

Trigger terms: "update lore", "edit lore", "record decision", "curate lore",
"add to lore", "lore decision", "update plan in lore".

## Relationship to other skills

| Skill | Relationship |
|-------|-------------|
| create-lore | Creates Works that edit-lore modifies |
| read-lore | Preferred way to inspect before editing; surfaces gaps edit-lore fills |
| sync-lore | Pushes lore commits after editing |
| protocol.md | Curation heuristic and identity rules |

edit-lore does not call other skills. It may suggest read-lore after editing
so the user can review the result.

## Explicit non-goals

- No working-tree lore files
- No automatic/background lore updates
- No lore ref deletion or archival
- No lore history rewriting (no amend/rebase of lore commits)
- No merge of diverged lore refs (document friction, defer to sync-lore era)
- No executable or helper scripts
- No mixing lore commits into source-branch history

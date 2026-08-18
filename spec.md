# Specification

## Purpose

create-lore is the agent skill that initialises a new Lore Work. It creates
an independent Git history at `refs/lore/<work-id>` containing an initial
`plan.md`, and optionally associates the Work with the current branch.

## Inputs

| Input | Required | Default |
|-------|----------|---------|
| `work-id` | No | Sanitized current branch name |
| Plan content | No | Agent derives from conversation context |

### Work ID rules

- Lowercase slug: letters, digits, hyphens
- Examples: `create-lore`, `git-lore`, `auth-refactor`
- Maps to ref: `refs/lore/<work-id>`
- Must not already exist locally

If the proposed work-id already exists, stop and ask whether to read the
existing Work (read-lore) or choose a different id.

### Branch name sanitization

When defaulting from branch name:

- Lowercase
- Replace `/` and other non-alphanumeric characters with `-`
- Collapse consecutive hyphens
- Strip leading/trailing hyphens

Example: branch `skills/create-lore` → work-id `skills-create-lore` or,
if the agent and user agree on scope, a shorter id like `create-lore`.

## Outputs

| Output | Description |
|--------|-------------|
| Lore ref | `refs/lore/<work-id>` pointing to initial commit |
| `plan.md` | Single markdown file in the lore tree |
| Branch config | `branch.<current-branch>.lore = <work-id>` (local only) |

No files are written to the repository working tree.

## Lore commit semantics

The initial lore commit is an **orphan commit** with no parents. It must
not parent off any source-branch commit. This guarantees lore history
remains independent of source rebases, squashes, and amendments.

Commit message format:

    lore: initialise Work <work-id>

Subsequent lore edits (via edit-lore) append to this history with normal
parent linkage.

## Git plumbing workflow

create-lore uses Git plumbing to avoid touching the working tree:

```bash
WORK_ID=<work-id>

# 1. Guard: ref must not exist
git rev-parse "refs/lore/${WORK_ID}" 2>/dev/null && exit 1

# 2. Write plan.md blob
PLAN_BLOB=$(printf '%s' "$PLAN_CONTENT" | git hash-object -w --stdin)

# 3. Build tree
TREE=$(printf '100644 blob %s\tplan.md' "$PLAN_BLOB" | git mktree)

# 4. Create orphan commit
COMMIT=$(git commit-tree "$TREE" -m "lore: initialise Work ${WORK_ID}")

# 5. Update ref
git update-ref "refs/lore/${WORK_ID}" "$COMMIT"
```

For multiple files at creation (not the default — create-lore should only
create plan.md), build the tree with multiple entries piped to `git mktree`.

## Branch association

After creating the lore ref, associate with the current branch:

```bash
BRANCH=$(git branch --show-current)
git config "branch.${BRANCH}.lore" "${WORK_ID}"
```

Properties:

- Local metadata only — not committed to the repository
- Read by read-lore and edit-lore when no work-id is explicitly given
- Can be changed or cleared without affecting the Lore Work itself
- Experimental: Lore open questions note worktree/discovery behaviour is
  unresolved; this config key is the v1 approach

## plan.md content requirements

The initial plan.md should include:

- **Objective** — what this work is trying to achieve (1–3 sentences)
- **Scope** — what is in and out of scope for this work
- **Next steps** — immediate actionable items

The plan should be concise and durable. It is not a transcript of the
conversation that led to creating the work.

Heuristic (from git-lore protocol):

> Would a future developer or agent otherwise have to rediscover this?

If yes, include it. If it is ephemeral conversation detail, omit it.

## SKILL.md structure

The skill file at `skills/create-lore/SKILL.md` must include:

- YAML frontmatter: `name`, `description` (third person, with trigger terms)
- Step-by-step workflow matching this spec
- The git plumbing commands above
- Link to `skills/protocol.md` for shared protocol details
- Post-creation verification commands
- Explicit note: creates only plan.md; other artefacts emerge via edit-lore

Trigger terms for description: "create lore", "start lore", "new lore work",
"initialise lore".

## Verification

After creation, the agent should confirm:

```bash
git for-each-ref refs/lore --format='%(refname:short)'
git log --oneline refs/lore/<work-id>
git show refs/lore/<work-id>:plan.md
git config --get "branch.$(git branch --show-current).lore"
```

Confirm:

- Lore ref exists and points to a single commit
- Only `plan.md` is in the tree
- Lore commit is not an ancestor of any source-branch commit
- No lore files in the working tree (`git status` is clean)

## Relationship to other skills

| Skill | Relationship |
|-------|-------------|
| read-lore | Reads a Work created by create-lore |
| edit-lore | Modifies a Work created by create-lore |
| sync-lore | Transports a Work created by create-lore to/from remote |
| protocol.md | Shared reference for work-id rules and curation |

create-lore does not call other skills. It reports what was created and
suggests read-lore if the user wants to review the result.

## Explicit non-goals

- No executable, CLI, or shell scripts in the skill directory
- No `_lore.json` or other manifest
- No automatic creation of decisions.md, spec.md, questions.md, etc.
- No lore transport (push/fetch) — that is sync-lore
- No GitHub-specific behaviour
- No working-tree checkout of lore files

# Specification

## Purpose

read-lore is the agent skill that loads and presents an existing Lore Work.
It resolves the Work identity, reads documents in a structured order,
shows recent lore history, and synthesizes a handoff summary for onboarding.

read-lore does not modify Lore. It does not write to the working tree.

## Inputs

| Input | Required | Default |
|-------|----------|---------|
| `work-id` | No | Branch config, else list and ask |
| Document | No | All documents (structured order) |
| History depth | No | 10 commits |

### Work-id resolution

Resolve in order:

1. **Explicit argument** — user or caller provides `work-id`
2. **Branch config** — `git config --get branch.<current-branch>.lore`
3. **Disambiguation** — list Works and ask:

```bash
git for-each-ref refs/lore --format='%(refname:short)'
```

If branch config points to a ref that does not exist, report the stale
config and fall through to disambiguation.

If zero Works exist, suggest create-lore.

### Read modes

| Mode | Trigger | Behaviour |
|------|---------|-----------|
| Full | Default | Structured read order + handoff summary |
| Document | User names a file | `git show refs/lore/<id>:<path>` only |
| Summary | "summarize lore", "handoff" | Handoff summary without full document dump |
| History | "lore history", "lore log" | `git log` only, no document content |

## Outputs

| Output | Description |
|--------|-------------|
| Work metadata | work-id, ref, latest commit hash |
| Document list | Files in lore tree |
| Lore history | Recent commits (`git log --oneline`) |
| Document content | Per structured order (mode-dependent) |
| Handoff summary | Synthesized brief: what / why / decided / rejected / next |

The handoff summary is **generated at read time** from lore content. It is
not appended to Lore unless the user later curates it via edit-lore.

## Structured read order

When reading all documents, follow this priority. Skip missing files silently.

| Priority | File | Role |
|----------|------|------|
| 1 | `plan.md` | Current direction and next steps |
| 2 | `decisions.md` | Settled choices |
| 3 | `questions.md` | Unresolved items |
| 4 | `spec.md` | Technical specification |
| 5 | `investigation.md` | Research and alternatives |
| 6 | `architecture.md` | System design |
| 7 | `handoff.md` | Curated handoff (if present) |
| 8 | Other | Remaining files, alphabetical |

Rationale: addresses git-lore open question "What does read Lore actually
mean?" — distinguish current state, decisions, unresolved questions, and
historical context without requiring a fixed taxonomy (files absent are
skipped).

Do **not** concatenate all files into one block without structure. Present
each document under a clear heading or section.

## Git commands

### List documents

```bash
WORK_ID=<work-id>
git ls-tree -r --name-only "refs/lore/${WORK_ID}"
```

### Lore history

```bash
git log --oneline -10 "refs/lore/${WORK_ID}"
```

Default depth: 10. Increase only when user asks for full lore history.

Do not dump full commit diffs by default. Offer `git diff refs/lore/<id>~1
refs/lore/<id>` if user wants last change detail.

### Read document

```bash
git show "refs/lore/${WORK_ID}:plan.md"
```

Generalise path for any document in the lore tree.

### Verify Work exists

```bash
git rev-parse "refs/lore/${WORK_ID}"
```

Fails if Work does not exist locally. Suggest sync-lore (if remote may have
it) or create-lore.

## Handoff summary format

After reading structured content, synthesize a brief summary:

```markdown
## Lore handoff: <work-id>

**What:** [What this work / project is — 1–2 sentences]
**Why:** [Why it exists / problem being solved]
**Decided:** [Key settled choices — bullet list]
**Rejected:** [Approaches ruled out and why — bullet list, if any]
**Next:** [Immediate next steps from plan.md or obvious gaps]
**Open:** [Unresolved questions from questions.md — bullet list, if any]
```

Keep the summary concise. If a section has no content in Lore, say
"Not recorded in Lore" rather than inventing detail.

This format maps directly to the git-lore success criterion.

## SKILL.md structure

The skill file at `skills/read-lore/SKILL.md` must include:

- YAML frontmatter: `name`, `description` (third person, trigger terms)
- Step-by-step workflow matching this spec
- Work-id resolution flow
- Structured read order table
- Handoff summary template
- Git commands above
- Link to `skills/protocol.md` and `refs/lore/git-lore`

Trigger terms: "read lore", "load lore", "show lore", "what is the lore",
"onboard to work", "handoff", "summarize lore", "lore context".

## Workflow checklist

```
- [ ] Step 1: Resolve work-id
- [ ] Step 2: Verify Work exists
- [ ] Step 3: Show metadata and document list
- [ ] Step 4: Show recent lore history
- [ ] Step 5: Read documents (structured order or specific file)
- [ ] Step 6: Synthesize handoff summary
- [ ] Step 7: Suggest edit-lore or sync-lore if appropriate
```

### Step 7 suggestions

- Lore stale or missing context → suggest edit-lore
- Work not found locally but may exist on remote → suggest sync-lore
- User wants to change Lore → edit-lore (not read-lore)

## Relationship to other skills

| Skill | Relationship |
|-------|-------------|
| create-lore | Creates Works that read-lore reads |
| edit-lore | Modifies Works after read-lore surfaces gaps |
| sync-lore | Fetches Works not present locally |
| protocol.md | Shared curation and identity rules |

read-lore does not call other skills. It may suggest them in step 7.

## Verification

After reading, confirm commands work:

```bash
git rev-parse refs/lore/<work-id>
git ls-tree -r --name-only refs/lore/<work-id>
git log --oneline -3 refs/lore/<work-id>
git show refs/lore/<work-id>:plan.md
```

Confirm:

- No files written to working tree
- Handoff summary reflects lore content, not conversation
- Missing documents skipped without error

## Explicit non-goals

- No working-tree checkout of lore files
- No automatic AI summarisation beyond structured presentation
- No lore modification
- No multi-work reads in a single invocation
- No executable or helper scripts
- No fetching from remote (sync-lore)

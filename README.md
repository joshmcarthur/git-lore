# git-lore

**Code tells you what the software does. Git history tells you what changed. Lore tells you why the work became what it is.**

git-lore is a small collection of agent skills that use ordinary Git primitives to maintain durable context about software work. Lore lives in independent refs (`refs/lore/<work-id>`) as ordinary Markdown. There is no daemon, database, manifest, or proprietary format — just skills plus Git conventions.

```bash
git switch -c add-widget
# Ask your agent: "Create lore for this work"
# ... work normally ...
# Ask your agent: "Record this decision in lore"
# Later: "Read lore" or "Summarize lore for handoff"
```

The Lore ref survives branch rebases, squash merges, and branch deletion. Source history and reasoning history have different lifecycles.

---

## Why Lore exists

Most project context lives in places with the wrong lifecycle:

| Where context lives | Problem |
|---------------------|---------|
| Branch | Rebases, squashes, and deletions rewrite or erase the association |
| Commit messages | Too small; reasoning gets compressed or omitted |
| Git Notes | Annotate commits — when a commit is rebased away, the note does not follow the *work* |
| Docs in the tree | Merge conflicts; clutter; tied to tree layout |
| AI memory | Provider-specific; not shared project history |
| Chat transcripts | Ephemeral, noisy, not curated |

Lore is deliberately **not** a transcript archive, AI memory system, documentation generator, database, new VCS, or CLI that hides Git.

It is curated, work-scoped context with its own Git history — portable, inspectable with `git show`, and shareable through the same remote as your code.

---

## The Work model

A **Work** is identified by `refs/lore/<work-id>` — an independent Git history belonging to a piece of work, not to a branch or commit.

```
refs/heads/feature/foo
         │
         │ source history
         ▼
    code commits

refs/lore/oauth2
         │
         │ work history
         ▼
    plans / decisions /
    investigations / handoff
```

A Work can span branches. The source branch can be rebased, renamed, merged, squashed, or deleted. Lore keeps going.

Local discovery uses branch metadata (not committed):

```bash
git config branch.<branch-name>.lore <work-id>
```

The shared identity is always the Lore ref.

---

## See it in action — dogfooding git-lore

This repository was built by dogfooding git-lore. The Lore refs are the historical evidence.

```
Bootstrap Lore (refs/lore/git-lore)
         │
         ▼
    create-lore skill
         │
         ▼
    Lore for read-lore (refs/lore/read-lore)
         │
         ▼
    read-lore skill
         │
         ▼
    Lore for edit-lore (refs/lore/edit-lore)
         │
         ▼
    edit-lore skill
         │
         ▼
    Lore for sync-lore (refs/lore/sync-lore)
         │
         ▼
    sync-lore skill
```

Each skill was developed using the previous skill to create and maintain Lore for the next Work. This is not a hypothetical workflow — it is how the project was actually built.

| Stage | Source commit | Lore Work | Lore HEAD |
|-------|---------------|-----------|-----------|
| Protocol + bootstrap | — | [`refs/lore/git-lore`](https://github.com/joshmcarthur/git-lore/commit/d39214e39944d2bdbf012452d39ba3e8d55761e9) | `d39214e` |
| create-lore | [`7ec772d`](https://github.com/joshmcarthur/git-lore/commit/7ec772d) | [`refs/lore/create-lore`](https://github.com/joshmcarthur/git-lore/commit/ee5791e5a1ad6774e4098c3eee8d463cfd7c9930) | `ee5791e` |
| read-lore | [`1e85201`](https://github.com/joshmcarthur/git-lore/commit/1e85201) | [`refs/lore/read-lore`](https://github.com/joshmcarthur/git-lore/commit/1de52226bcefeab727430ced3c719f920251979f) | `1de5222` |
| edit-lore | [`a791d0d`](https://github.com/joshmcarthur/git-lore/commit/a791d0d) | [`refs/lore/edit-lore`](https://github.com/joshmcarthur/git-lore/commit/86c2f8c4b06824ec09637c0a2f0dd9c1f9fd92b8) | `86c2f8c` |
| sync-lore | [`3095a7b`](https://github.com/joshmcarthur/git-lore/commit/3095a7b) | [`refs/lore/sync-lore`](https://github.com/joshmcarthur/git-lore/commit/dcb7b8e5fdbc32e95b9241e52e391d85e0a95baa) | `dcb7b8e` |

Lore commit links browse the lore tree at that revision on GitHub (files live at the tree root, not under a source path). Lore refs do not appear in GitHub's branch picker — link by lore commit SHA.

### Bootstrap: what we were trying to prove

From [`refs/lore/git-lore` plan.md](https://github.com/joshmcarthur/git-lore/blob/d39214e39944d2bdbf012452d39ba3e8d55761e9/plan.md):

> The key test is not whether we can store Markdown in Git. We obviously can.
>
> The test is whether the resulting workflow is useful enough that people actually want to maintain this context.

### Decision that shaped everything: Work, not branch

From [`refs/lore/git-lore` decisions.md](https://github.com/joshmcarthur/git-lore/blob/d39214e39944d2bdbf012452d39ba3e8d55761e9/decisions.md):

> **Decision:** The primary identity is a Work rather than a branch.
>
> Branches can be rebased, renamed, split, merged or deleted. The reasoning associated with the work should have a longer and more stable identity.

### Discovery during edit-lore: isolated index

Dogfooding surfaced a real implementation constraint. From [`refs/lore/edit-lore` decisions.md](https://github.com/joshmcarthur/git-lore/blob/86c2f8c4b06824ec09637c0a2f0dd9c1f9fd92b8/decisions.md):

> **Decision:** edit-lore must set `GIT_INDEX_FILE` to a temp file when building the new lore tree. Using `git --work-tree=<tmpdir> add -A` against the shared index pollutes the main repository index.

This discovery in one Work's Lore informed the edit-lore skill that every subsequent Work uses.

### Open question carried forward

From [`refs/lore/git-lore` questions.md](https://github.com/joshmcarthur/git-lore/blob/d39214e39944d2bdbf012452d39ba3e8d55761e9/questions.md):

> **How much should agents automatically update Lore?**
>
> Too little: Lore becomes stale. Too much: Lore becomes noisy.
>
> This is probably the most important UX question.

edit-lore's mandatory curation gate is a partial answer; the trigger problem remains open in [`refs/lore/edit-lore` questions.md](https://github.com/joshmcarthur/git-lore/blob/86c2f8c4b06824ec09637c0a2f0dd9c1f9fd92b8/questions.md).

### read-lore: presentation matters

From [`refs/lore/read-lore` plan.md](https://github.com/joshmcarthur/git-lore/blob/1de52226bcefeab727430ced3c719f920251979f/plan.md):

> The test is not whether we can read Markdown from a Git ref — it is whether the **presentation** reduces context reconstruction.

read-lore presents documents in priority order and synthesizes a handoff summary at read time — not a wall of concatenated Markdown.

---

## The skills

Skills live in `skills/`. Each is a `SKILL.md` file instructing an agent to perform ordinary Git operations. Install by pointing your agent at this repository (see [Getting started](#getting-started)).

### [create-lore](skills/create-lore/SKILL.md)

Starts a new Work and establishes initial context (`plan.md` only).

**When:** beginning a piece of work that needs durable context.

**Say to your agent:** "Create lore for this work" / "Start lore" / "Initialise a lore work"

**Does:** creates orphan commit at `refs/lore/<work-id>`, sets `branch.<name>.lore`, leaves working tree clean.

### [read-lore](skills/read-lore/SKILL.md)

Loads Lore into an agent's working context — structured document order plus handoff summary.

**When:** onboarding a fresh agent or human; understanding what was decided.

**Say to your agent:** "Read lore" / "Summarize lore" / "Handoff for this work"

**Does:** `git show` on lore ref; presents plan → decisions → questions → spec → …; synthesizes what/why/decided/rejected/next/open.

### [edit-lore](skills/edit-lore/SKILL.md)

Curates durable discoveries, decisions, plan updates, and questions into a new lore commit.

**When:** a durable conclusion emerged that a future reader would otherwise rediscover.

**Say to your agent:** "Record this decision in lore" / "Update lore" / "Curate lore"

**Does:** curation gate → export via `git archive` → edit in temp dir → commit with isolated index → advance lore ref.

### [sync-lore](skills/sync-lore/SKILL.md)

Shares Lore refs through Git remotes.

**When:** pushing Lore to teammates; fetching Lore after clone.

**Say to your agent:** "Sync lore" / "Push lore" / "Fetch lore" / "Lore status"

**Does:** configures `refs/lore/*` fetch refspec; fetch/push/status; detects divergence without auto-merge.

Protocol reference shared by all skills: [skills/protocol.md](skills/protocol.md)

---

## What belongs in Lore

**Good Lore:**

- "We chose PostgreSQL because…"
- "The apparent race condition was actually caused by…"
- "We rejected approach X because…"
- "This API cannot be changed because…"
- "The remaining uncertainty is…"
- "A future implementer should know…"

**Bad Lore:**

- full conversation transcripts
- every command executed
- routine implementation details
- information already obvious from the code
- speculative filler
- status updates with no durable value

Heuristic (from the protocol):

> Would a future developer or agent otherwise have to rediscover this?

If yes, it probably belongs in Lore. edit-lore applies this as a mandatory gate before every write.

---

## How it works technically

Lore is just Git refs and Git objects:

```
refs/lore/<work-id>
         │
         ▼
       commit
         │
         ▼
   Markdown files
```

**Initial commit:** orphan (no parents) — lore history is independent of source-branch rebases and squashes.

**Subsequent commits:** parent off current lore HEAD — a curation log, not source commits.

**Messages:** `lore: <what changed>`

**Garbage collection:** the `refs/lore/*` ref keeps lore objects reachable. Deleting the ref is how you remove a Work.

**Inspect without skills:**

```bash
git for-each-ref refs/lore
git log refs/lore/git-lore
git show refs/lore/git-lore:plan.md
git diff refs/lore/git-lore~1 refs/lore/git-lore
```

**Transport** (not automatic — configure once per remote):

```bash
git config --add remote.origin.fetch 'refs/lore/*:refs/lore/*'
git fetch origin 'refs/lore/*:refs/lore/*'
git push origin 'refs/lore/*'
```

No `+` prefix on fetch — unpushed local lore curation must not be overwritten by fetch.

---

## Why not just use Markdown?

Ordinary Markdown in the repository is an excellent solution for many projects. Lore is interesting only when reasoning has a **different lifecycle** from the source tree.

Potential advantages:

- no source-tree clutter
- independent history
- no merge conflicts with docs
- survives branch rewriting
- can outlive the source branch
- shared through Git
- accessible to humans and different AI agents

Whether these justify the additional convention is an empirical question. git-lore exists to test it — by dogfooding.

---

## Why not Git Notes?

Git Notes annotate existing Git objects. Lore represents an independent history belonging to a Work.

When a source commit is rebased or squashed, its Notes do not automatically follow the work. Lore does — because it is not attached to source commits at all.

---

## Getting started

1. **Clone this repository** (or copy the `skills/` directory into your agent's skill path).

2. **Start work on a branch:**

   ```bash
   git switch -c my-feature
   ```

3. **Ask your agent to create Lore:**

   > Create lore for this work.

   The agent runs [create-lore](skills/create-lore/SKILL.md). Verify:

   ```bash
   git for-each-ref refs/lore
   git show refs/lore/my-feature:plan.md
   git status   # should be clean — no lore files in the tree
   ```

4. **Work normally.** When something durable emerges:

   > Record this decision in lore: we chose X because Y.

5. **Hand off to another agent or session:**

   > Read lore for this work. Summarize for handoff.

6. **Share Lore** (when you have a remote):

   > Sync lore.

   Or configure and run transport yourself (see [How it works technically](#how-it-works-technically)).

---

## Lore in this repository

| Lore Work | Role | Browse on GitHub |
|-----------|------|------------------|
| `refs/lore/git-lore` | Protocol authority — spec, decisions, investigations, open questions | [tree](https://github.com/joshmcarthur/git-lore/tree/d39214e39944d2bdbf012452d39ba3e8d55761e9) · [plan.md](https://github.com/joshmcarthur/git-lore/blob/d39214e39944d2bdbf012452d39ba3e8d55761e9/plan.md) · [spec.md](https://github.com/joshmcarthur/git-lore/blob/d39214e39944d2bdbf012452d39ba3e8d55761e9/spec.md) |
| `refs/lore/create-lore` | create-lore skill plan and spec | [tree](https://github.com/joshmcarthur/git-lore/tree/ee5791e5a1ad6774e4098c3eee8d463cfd7c9930) · [spec.md](https://github.com/joshmcarthur/git-lore/blob/ee5791e5a1ad6774e4098c3eee8d463cfd7c9930/spec.md) |
| `refs/lore/read-lore` | read-lore skill plan and spec | [tree](https://github.com/joshmcarthur/git-lore/tree/1de52226bcefeab727430ced3c719f920251979f) · [spec.md](https://github.com/joshmcarthur/git-lore/blob/1de52226bcefeab727430ced3c719f920251979f/spec.md) |
| `refs/lore/edit-lore` | edit-lore skill plan, spec, and implementation discoveries | [tree](https://github.com/joshmcarthur/git-lore/tree/86c2f8c4b06824ec09637c0a2f0dd9c1f9fd92b8) · [spec.md](https://github.com/joshmcarthur/git-lore/blob/86c2f8c4b06824ec09637c0a2f0dd9c1f9fd92b8/spec.md) |
| `refs/lore/sync-lore` | sync-lore skill plan, spec, and transport decisions | [tree](https://github.com/joshmcarthur/git-lore/tree/dcb7b8e5fdbc32e95b9241e52e391d85e0a95baa) · [spec.md](https://github.com/joshmcarthur/git-lore/blob/dcb7b8e5fdbc32e95b9241e52e391d85e0a95baa/spec.md) |
| `refs/lore/readme` | README authoring work | [tree](https://github.com/joshmcarthur/git-lore/tree/b763033816fb1b4db287664ae63a2744e08f4f97) |

Locally:

```bash
git for-each-ref refs/lore --format='%(refname:short) %(objectname:short)'
git show refs/lore/git-lore:plan.md
```

On GitHub, lore is browsable by **lore commit SHA** (`blob/<sha>/plan.md`). Lore refs are not listed in the branch/tag picker, but each lore commit is a normal commit object once pushed. Link by SHA for permalinks; use `git log refs/lore/<work-id>` to find the current HEAD.

Skill files in the working tree:

- [skills/create-lore/SKILL.md](skills/create-lore/SKILL.md)
- [skills/read-lore/SKILL.md](skills/read-lore/SKILL.md)
- [skills/edit-lore/SKILL.md](skills/edit-lore/SKILL.md)
- [skills/sync-lore/SKILL.md](skills/sync-lore/SKILL.md)
- [skills/protocol.md](skills/protocol.md)

---

## Status and philosophy

This project is **experimental**.

The protocol is intentionally tiny. The interesting question is not whether Git can store this — it obviously can. The interesting question is whether durable, curated work context is useful enough that developers and agents will actually maintain it.

What exists today:

- Four core skills (create, read, edit, sync)
- A distilled protocol reference in the working tree
- Lore refs documenting the protocol and each skill's implementation
- Dogfooding evidence that the workflow can carry context from one Work to the next

What is deliberately **not** here (yet):

- a daemon, database, hosted service, MCP server, or CLI
- a manifest or transcript store
- GitHub-specific integrations
- automatic lore updates or AI summarisation pipelines

That restraint is a feature. The project is trying to prove the workflow before building infrastructure around it.

Success criterion (from [`refs/lore/git-lore` plan.md](https://github.com/joshmcarthur/git-lore/blob/d39214e39944d2bdbf012452d39ba3e8d55761e9/plan.md)): hand this repository to a fresh agent with no conversational history. If Lore makes it materially easier to explain what git-lore is, why it exists, what was decided, what was rejected, and what needs to happen next — the experiment is working.

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

The Lore ref survives branch rebases, squash merges, and branch deletion. Source history and reasoning history have different lifecycles. For how that is implemented in Git (refs, objects, plumbing), see [How it works technically](#how-it-works-technically).

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
| Protocol + bootstrap | — | [`refs/lore/git-lore`](https://github.com/joshmcarthur/git-lore/tree/d39214e39944d2bdbf012452d39ba3e8d55761e9) | `d39214e` |
| create-lore | [`7ec772d`](https://github.com/joshmcarthur/git-lore/commit/7ec772d) | [`refs/lore/create-lore`](https://github.com/joshmcarthur/git-lore/tree/ee5791e5a1ad6774e4098c3eee8d463cfd7c9930) | `ee5791e` |
| read-lore | [`1e85201`](https://github.com/joshmcarthur/git-lore/commit/1e85201) | [`refs/lore/read-lore`](https://github.com/joshmcarthur/git-lore/tree/1de52226bcefeab727430ced3c719f920251979f) | `1de5222` |
| edit-lore | [`a791d0d`](https://github.com/joshmcarthur/git-lore/commit/a791d0d) | [`refs/lore/edit-lore`](https://github.com/joshmcarthur/git-lore/tree/86c2f8c4b06824ec09637c0a2f0dd9c1f9fd92b8) | `86c2f8c` |
| sync-lore | [`3095a7b`](https://github.com/joshmcarthur/git-lore/commit/3095a7b) | [`refs/lore/sync-lore`](https://github.com/joshmcarthur/git-lore/tree/dcb7b8e5fdbc32e95b9241e52e391d85e0a95baa) | `dcb7b8e` |

On GitHub, link lore by **tree** (`tree/<lore-commit-sha>`) to browse the file list — clearer than a commit page, which often shows only a status-update message. Lore files live at the tree root, not under a source path. Lore refs do not appear in GitHub's branch picker; use the lore commit SHA from `git log refs/lore/<work-id>`.

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

Lore is ordinary Git. The unusual part is *which* Git features it uses: named refs and plumbing commands, not files in the working tree.

### What a ref is

A **ref** is a name that points at a commit SHA. You already use them:

| Name you know | Ref | Points at |
|---------------|-----|-----------|
| Branch `main` | `refs/heads/main` | latest commit on that branch |
| Tag `v1.0` | `refs/tags/v1.0` | a specific commit |
| Lore Work `oauth2` | `refs/lore/oauth2` | latest lore commit for that work |

Git stores refs under `.git/refs/` (or in `.git/packed-refs`). Updating a branch is updating a ref. Creating lore is the same operation in a different namespace: write a commit, then point `refs/lore/<work-id>` at it.

That namespace is a convention, not a Git feature. Git does not treat `refs/lore/*` specially. Default `git fetch` / `git push` ignore it, which is why lore has to be fetched and pushed explicitly.

### What Git stores (objects)

Git's object database (`.git/objects`) holds three kinds of object lore cares about:

| Object | What it is |
|--------|------------|
| **Blob** | The bytes of one file (`plan.md`, `decisions.md`, …) |
| **Tree** | A directory listing: names → blobs (or nested trees) |
| **Commit** | A tree, an optional parent commit, a message, and author metadata |

A branch commit usually has a parent (the previous commit on that branch). The first lore commit is an **orphan**: it has no parent, so lore history is not attached to `main` or to your feature branch. Later lore commits parent off the previous *lore* commit. Rebases, squashes, and deleting the source branch do not move or rewrite that chain.

```
refs/heads/feature/foo          refs/lore/oauth2
         │                               │
         ▼                               ▼
    code commits                   lore commits
    (parented on                   (orphan root,
     the source branch)             then lore→lore)
```

The lore ref is what keeps those objects alive. `git gc` collects unreachable objects; as long as `refs/lore/<work-id>` exists, the lore commits and their trees/blobs stay. Deleting the ref is how you remove a Work.

### Why nothing appears in `git status`

Everyday Git writes files into the working tree, then `git add` / `git commit` snapshots them. Lore never does that.

The skills (and the optional CLI) build objects **directly** with plumbing:

- `git hash-object -w` — store file bytes as a blob
- `git mktree` / `git write-tree` — store a directory listing as a tree
- `git commit-tree` — store a commit that points at that tree
- `git update-ref` — move the lore ref to the new commit

`git show refs/lore/oauth2:plan.md` reads a blob out of the object database. No checkout, so the working tree stays clean and lore files cannot collide with source files of the same name.

Branch discovery is separate and local: `git config branch.<branch>.lore <work-id>` is not a ref and is not committed. The shared identity is always the lore ref.

### What the skills do in Git

#### Create (create-lore)

Builds the first lore commit from `plan.md` only, then names it:

```
plan.md bytes
    │  git hash-object -w
    ▼
  blob
    │  git mktree
    ▼
  tree (plan.md)
    │  git commit-tree   ← no parent (orphan)
    ▼
  commit  "lore: initialise Work <id>"
    │  git update-ref refs/lore/<id>
    ▼
  named pointer
```

Also sets `branch.<current>.lore` so later reads/edits can find the Work without asking for an id.

#### Read (read-lore)

Does not check anything out. It resolves `refs/lore/<id>` (from an argument, or from `branch.*.lore`), lists files with `git ls-tree`, reads each document with `git show <ref>:<path>`, and shows recent lore history with `git log` on that ref. The handoff summary is generated at read time; it is not stored unless someone later curates a `handoff.md`.

#### Change (edit-lore)

Advances the same ref with a new commit whose parent is the current lore HEAD:

```
refs/lore/<id>  ──►  existing commit
                         │
                         │  git archive → temp directory → edit Markdown
                         │  isolated GIT_INDEX_FILE → git write-tree
                         ▼
                       new tree
                         │  git commit-tree -p <old commit>
                         ▼
                       new commit  "lore: <what changed>"
                         │  git update-ref refs/lore/<id>
                         ▼
                   ref now points here
```

The isolated index matters: `git add` against the repo's real index would stage lore files as if they belonged to the source branch. Edits happen in a temp directory; that directory is deleted after the commit exists.

Lore commits are never mixed into the source branch. One lore commit per curation event; messages use the `lore:` prefix.

#### Store and share (sync-lore)

Create and edit already stored the objects in this clone's `.git`. "Storing" lore on a remote is ordinary fetch/push of those refs.

`refs/lore/*` is **not** in the default fetch refspec. After clone, lore is absent until you add one and fetch:

```bash
git config --add remote.origin.fetch 'refs/lore/*:refs/lore/*'
git fetch origin 'refs/lore/*:refs/lore/*'
git push origin 'refs/lore/*'
```

No `+` prefix on fetch — that would force-update local lore refs and discard unpushed curation. Diverged lore histories are not merged automatically (Markdown conflict markers are a bad outcome for curated notes). sync-lore fetches the remote side to `refs/lore/<id>-remote` so you can compare and reconcile by hand.

### Inspect without skills

```bash
git for-each-ref refs/lore
git log refs/lore/git-lore
git show refs/lore/git-lore:plan.md
git diff refs/lore/git-lore~1 refs/lore/git-lore
git ls-tree -r --name-only refs/lore/git-lore
```

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

### Install the skills

git-lore skills follow the [Agent Skills](https://agentskills.io) format — one `SKILL.md` per skill under `skills/`. Install them into your coding agent using any of the methods below.

#### skills.sh (recommended)

The [skills.sh](https://skills.sh) ecosystem is powered by the open [`skills` CLI](https://github.com/vercel-labs/skills). It works across Cursor, Claude Code, Codex, GitHub Copilot, Windsurf, and [70+ other agents](https://github.com/vercel-labs/skills#supported-agents).

```bash
# List available skills in this repo
npx skills add joshmcarthur/git-lore --list

# Install all four skills (interactive — detects your installed agents)
npx skills add joshmcarthur/git-lore

# Install all skills globally, non-interactively
npx skills add joshmcarthur/git-lore --all -g -y

# Install specific skills to specific agents
npx skills add joshmcarthur/git-lore \
  --skill create-lore --skill read-lore --skill edit-lore --skill sync-lore \
  -a cursor -a claude-code -a codex \
  -g -y
```

| Scope | Flag | Where skills land |
|-------|------|-------------------|
| Project (default) | — | Agent-specific project directory (e.g. `.agents/skills/` for Cursor and Codex, `.claude/skills/` for Claude Code) |
| Global | `-g` | User home directory (e.g. `~/.cursor/skills/`, `~/.claude/skills/`, `~/.codex/skills/`) |

After installing, restart your agent or reload the window so it picks up the new skills. In most agents you can verify with a skills command (e.g. `/skills` in Claude Code) or by asking your agent to "create lore for this work".

Update later with `npx skills update` or remove with `npx skills remove`.

#### GitHub CLI

If you use [GitHub's skill commands](https://docs.github.com/en/copilot/how-tos/use-copilot-agents/coding-agent/create-skills) (`gh skill install`):

```bash
gh skill install joshmcarthur/git-lore create-lore
gh skill install joshmcarthur/git-lore read-lore edit-lore sync-lore
```

Pin to a release tag for reproducible installs:

```bash
gh skill install joshmcarthur/git-lore create-lore --pin v0.1.0
gh skill install joshmcarthur/git-lore read-lore@v0.1.0
```

Without `--pin`, `gh skill install` resolves the latest release tag, then falls back to the default branch. See [Releases](https://github.com/joshmcarthur/git-lore/releases) and [CHANGELOG.md](CHANGELOG.md) for version history.

#### Manual install

Copy or symlink each skill directory into your agent's skills path:

```bash
git clone https://github.com/joshmcarthur/git-lore /tmp/git-lore

# Cursor (global)
ln -s /tmp/git-lore/skills/create-lore ~/.cursor/skills/create-lore
ln -s /tmp/git-lore/skills/read-lore  ~/.cursor/skills/read-lore
ln -s /tmp/git-lore/skills/edit-lore  ~/.cursor/skills/edit-lore
ln -s /tmp/git-lore/skills/sync-lore  ~/.cursor/skills/sync-lore

# Claude Code (global)
ln -s /tmp/git-lore/skills/create-lore ~/.claude/skills/create-lore
# ... repeat for read-lore, edit-lore, sync-lore

# Codex (global)
ln -s /tmp/git-lore/skills/create-lore ~/.codex/skills/create-lore
# ... repeat for read-lore, edit-lore, sync-lore
```

For project-scoped installs, use the agent's project path instead (e.g. `.agents/skills/` for Cursor and Codex, `.claude/skills/` for Claude Code). See the [supported agents table](https://github.com/vercel-labs/skills#supported-agents) for other tools.

Each skill links to [skills/protocol.md](skills/protocol.md) at `references/protocol.md` inside the skill directory (symlink in this repo; copied on install by the skills CLI).

### Listing on skills.sh

There is no separate submission form. [skills.sh](https://skills.sh) is a leaderboard built from anonymous install telemetry collected by the `skills` CLI. A skill appears once people install it from a public Git repository.

**What you need:**

1. A **public GitHub repository** (this repo qualifies).
2. Valid **`SKILL.md` files** with required YAML frontmatter (`name`, `description`) in a [discoverable location](https://github.com/vercel-labs/skills#skill-discovery) — git-lore uses `skills/<skill-name>/SKILL.md`.
3. **Installs via `npx skills add`** — ranking reflects install counts, not a manual review queue.

Optional: add an install badge to your README (replace `owner/repo` with your source):

```markdown
[![skills.sh](https://skills.sh/b/joshmcarthur/git-lore)](https://skills.sh/joshmcarthur/git-lore)
```

Browse individual skills at `https://skills.sh/<owner>/<repo>/<skill-name>` once installs accumulate. Skills are subject to routine security audits; report issues at [security.vercel.com](https://security.vercel.com).

### Use git-lore in your project

These steps assume the skills are already installed in your agent (see above). You do **not** need to clone this repository into your project — only the skills need to be available to your agent.

1. **Start work on a branch:**

   ```bash
   git switch -c my-feature
   ```

2. **Ask your agent to create Lore:**

   > Create lore for this work.

   The agent runs [create-lore](skills/create-lore/SKILL.md). Verify:

   ```bash
   git for-each-ref refs/lore
   git show refs/lore/my-feature:plan.md
   git status   # should be clean — no lore files in the tree
   ```

3. **Work normally.** When something durable emerges:

   > Record this decision in lore: we chose X because Y.

4. **Hand off to another agent or session:**

   > Read lore for this work. Summarize for handoff.

5. **Share Lore** (when you have a remote):

   > Sync lore.

   Or configure and run transport yourself (see [How it works technically](#how-it-works-technically)).

---

## Releasing

Source commits on `main` should follow [Conventional Commits](https://www.conventionalcommits.org/) (`feat:`, `fix:`, `docs:`, etc.). [release-please](https://github.com/googleapis/release-please) opens a Release PR that bumps [version.txt](version.txt) and [CHANGELOG.md](CHANGELOG.md). Merge that PR to create a git tag and GitHub Release.

When a release is published, CI also builds the optional [`git-lore` CLI](extensions/git-lore) for linux/darwin/windows and attaches archives (`bin/git-lore` + `lib/git-lore/{*.sh,serve}`) as `git-lore_<version>_<os>_<arch>.tar.gz` / `.zip`.

Lore commits (`lore: ...` on `refs/lore/*`) are separate from source history and are not included in release changelogs.

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
| `refs/lore/release-please` | Skill versioning and release automation | [tree](https://github.com/joshmcarthur/git-lore/tree/7b21e12e3adbd87e4d71c1ecc1db77d50b84a92b) · [spec.md](https://github.com/joshmcarthur/git-lore/blob/7b21e12e3adbd87e4d71c1ecc1db77d50b84a92b/spec.md) · [decisions.md](https://github.com/joshmcarthur/git-lore/blob/7b21e12e3adbd87e4d71c1ecc1db77d50b84a92b/decisions.md) |

Locally:

```bash
git for-each-ref refs/lore --format='%(refname:short) %(objectname:short)'
git show refs/lore/git-lore:plan.md
```

On GitHub, browse lore with **`tree/<lore-commit-sha>`** to see the file list; use **`blob/<sha>/<file>`** to link to a specific document. Lore refs are not listed in the branch/tag picker, but each lore commit is a normal commit object once pushed.

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
- Automated releases via release-please (semver tags and GitHub Releases)
- An optional **`git-lore` CLI** under [`extensions/git-lore`](extensions/git-lore) — Git external command (`git lore`) with shell wrappers for create/edit/sync and a Go `serve` binary for the browser UI

What is deliberately **not** here (yet):

- a daemon, database, hosted service, or MCP server
- a manifest or transcript store
- a GitHub App or formal marketplace submission (skills appear on [skills.sh](https://skills.sh) via `npx skills add` install telemetry)
- automatic lore updates or AI summarisation pipelines
- write/push from the serve UI (create/edit stay skill- or CLI-driven; the browser remains read-only)

That restraint is a feature. The project is trying to prove the workflow before building infrastructure around it. The repository root stays skills-first; executables live under `extensions/` so they do not crowd the primary agent-skills surface.

### git-lore CLI (optional extension)

Optional Git plugin in [`extensions/git-lore`](extensions/git-lore). Not required to use the skills. Install puts `git-lore` on your `PATH` so `git lore …` works (same convention as `git lfs`).

- **Shell scripts** — `list`, `show`, `create`, `edit`, `export`, `delete`, `branch`, `remote` (mirrors skill Git plumbing)
- **Go binary** — `serve` for the read-only Lore browser + JSON API

```bash
cd extensions/git-lore
make install PREFIX=~/.local   # or: make build && PATH=$PWD/bin:$PATH
git lore list
git lore serve --repo ../.. --open
# listens on http://127.0.0.1:9473 by default
```

Release archives include `bin/git-lore` plus `lib/git-lore/{*.sh,serve}` (`git-lore_<version>_<os>_<arch>.tar.gz` / `.zip`). On Windows, git-ops commands need Git Bash or WSL; `serve` runs natively.

Flags for `serve`: `--repo <path>`, `--addr host:port`, `--open`. See [`extensions/git-lore/README.md`](extensions/git-lore/README.md) for the full command list.

Success criterion (from [`refs/lore/git-lore` plan.md](https://github.com/joshmcarthur/git-lore/blob/d39214e39944d2bdbf012452d39ba3e8d55761e9/plan.md)): hand this repository to a fresh agent with no conversational history. If Lore makes it materially easier to explain what git-lore is, why it exists, what was decided, what was rejected, and what needs to happen next — the experiment is working.

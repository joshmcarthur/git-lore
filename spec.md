# Specification

## Purpose

Define how git-lore skill versions are released, tagged, and installed.

## Version model

- **One version for the whole repository** — all skills share the same tag
- Tag format: `v{semver}` (e.g. `v0.1.0`)
- Version source of truth: `version.txt` and `.github/.release-please-manifest.json`
- Changelog: root `CHANGELOG.md` (maintained by release-please)

Per-skill version tags are not used. `gh skill install` resolves versions
from repository git tags, not per-directory metadata.

## Release automation

release-please runs on every push to `main` via GitHub Actions.

1. Conventional commits on `main` accumulate in a Release PR
2. Merging the Release PR bumps `version.txt`, updates `CHANGELOG.md`,
   creates a git tag, and publishes a GitHub Release

Lore commits (`lore: ...` on `refs/lore/*`) are unaffected — release-please
only reads source-branch history.

### Configuration

| File | Role |
|------|------|
| `.github/release-please-config.json` | Release strategy (`simple` type) |
| `.github/.release-please-manifest.json` | Last-released version per package |
| `.github/workflows/release-please.yml` | CI trigger |
| `version.txt` | Human-readable current version |
| `CHANGELOG.md` | Release notes |

### Pre-1.0 policy

While version is `<1.0.0`, `feat:` commits bump the patch version
(`bump-minor-pre-major: true`). Matches experimental project status.

## Installing skills

```bash
# Latest release tag
gh skill install joshmcarthur/git-lore create-lore

# Pinned version
gh skill install joshmcarthur/git-lore create-lore --pin v0.1.0
gh skill install joshmcarthur/git-lore create-lore@v0.1.0
```

Available skills: `create-lore`, `read-lore`, `edit-lore`, `sync-lore`.

## Lore transport

Lore refs remain independent of skill releases. Transport via sync-lore
is unchanged — lore and skill versions have separate lifecycles.

## CLI release artifacts

After release-please publishes a GitHub Release, workflow job `release-binaries`
(in `.github/workflows/release-please.yml`) builds `extensions/git-lore` for:

- linux/amd64, linux/arm64
- darwin/amd64, darwin/arm64
- windows/amd64

Archives are named `git-lore_<version>_<os>_<arch>.tar.gz` (`.zip` on Windows)
and uploaded to the same release via `gh release upload`. The binary version
string is injected with `-ldflags -X main.version=<semver>`.

# Decisions

## Unified repo version

**Decision:** Tag one version for the entire skills collection, not
per-skill tags like `create-lore-v1.0.0`.

`gh skill install` resolves versions from repository git tags. A single
`v0.1.0` tag pins the whole skills snapshot at that commit.

## Simple release type

**Decision:** Use release-please `simple` release type with `version.txt`.

The repository has no `package.json` or other ecosystem manifest. The
`simple` type is designed for version-file-only repositories.

## No release skill

**Decision:** Do not create a new agent skill for releases.

Release automation is CI (GitHub Actions) plus README documentation.
Consistent with the git-lore decision to avoid executables before proving
the workflow (`refs/lore/git-lore:decisions.md`).

## No SKILL.md version fields

**Decision:** Do not add `version` to each skill's YAML frontmatter.

The git tag is the version pin. Duplicating version in five files adds
maintenance burden without enabling `gh skill install` per-skill pinning.

## GitHub Actions only

**Decision:** Use GitHub Actions for release automation, not a GitHub App
or hosted service.

This is release infrastructure, not a lore protocol extension. Lore refs
and skill distribution remain separate concerns.

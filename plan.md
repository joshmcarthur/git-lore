# Plan

## Objective

Automate semver releases for the git-lore skills collection so users can
install pinned versions via `gh skill install --pin vX.Y.Z`.

## Scope

In scope:

- release-please configuration and GitHub Actions workflow
- Unified repo version tags (not per-skill)
- GitHub Releases and CHANGELOG.md
- README documentation for pinned installs

Out of scope:

- Per-skill independent version tags
- `gh skill publish` / skill marketplace setup
- Version fields in SKILL.md frontmatter
- Versioning lore refs themselves
- A new agent skill for releases

## Next steps

1. Merge this PR (release-please setup)
2. Merge the first Release PR opened by CI on main
3. Verify `v0.1.0` tag and GitHub Release
4. Smoke test: `gh skill install joshmcarthur/git-lore create-lore --pin v0.1.0`

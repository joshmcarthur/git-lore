#!/usr/bin/env bats
# remote status classification (local-only, no network)

setup() {
  EXT="$(cd "$(dirname "$BATS_TEST_FILENAME")/.." && pwd)"
  export PATH="$EXT/bin:$PATH"
  REPO="$(mktemp -d "${TMPDIR:-/tmp}/git-lore-bats.XXXXXX")"
  git -C "$REPO" init -q
  git -C "$REPO" config user.email "test@example.com"
  git -C "$REPO" config user.name "Test"
  git -C "$REPO" commit --allow-empty -q -m "init"
  # Bare remote for local status (empty lore)
  REMOTE="$(mktemp -d "${TMPDIR:-/tmp}/git-lore-remote.XXXXXX")"
  git -C "$REMOTE" init -q --bare
  git -C "$REPO" remote add origin "$REMOTE"
  export REPO REMOTE
}

teardown() {
  rm -rf "$REPO" "$REMOTE"
}

@test "remote status shows local-only work" {
  printf '# Plan\n' | git-lore create local-only --repo "$REPO"
  run git-lore remote status --repo "$REPO"
  [ "$status" -eq 0 ]
  [[ "$output" == *"local-only"* ]]
  [[ "$output" == *"local-only"* ]] || true
  echo "$output" | grep -q $'local-only\tlocal-only'
}

@test "remote status without remote fails" {
  git -C "$REPO" remote remove origin
  run git-lore remote status --repo "$REPO"
  [ "$status" -ne 0 ]
  [[ "$output" == *"no such remote"* ]]
}

@test "branch list and set unset" {
  printf '# Plan\n' | git-lore create branch-work --repo "$REPO"
  git -C "$REPO" checkout -q -b feat
  git-lore branch --set branch-work --repo "$REPO"
  run git-lore branch --repo "$REPO"
  [ "$status" -eq 0 ]
  [[ "$output" == *"feat"* ]]
  [[ "$output" == *"branch-work"* ]]
  git-lore branch --unset --repo "$REPO"
  run git-lore branch --repo "$REPO"
  [ "$status" -eq 0 ]
  [[ "$output" == *"No branch lore associations"* ]]
}

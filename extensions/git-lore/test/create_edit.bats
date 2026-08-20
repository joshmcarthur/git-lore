#!/usr/bin/env bats
# create → show → edit → export → delete lifecycle

setup() {
  EXT="$(cd "$(dirname "$BATS_TEST_FILENAME")/.." && pwd)"
  export PATH="$EXT/bin:$PATH"
  REPO="$(mktemp -d "${TMPDIR:-/tmp}/git-lore-bats.XXXXXX")"
  git -C "$REPO" init -q
  git -C "$REPO" config user.email "test@example.com"
  git -C "$REPO" config user.name "Test"
  git -C "$REPO" commit --allow-empty -q -m "init"
  export REPO
}

teardown() {
  rm -rf "$REPO"
}

@test "create show edit export delete" {
  printf '# Plan\n\nObjective.\n' | git-lore create test-work --repo "$REPO"
  run git-lore show test-work plan.md --repo "$REPO"
  [ "$status" -eq 0 ]
  [[ "$output" == *"# Plan"* ]]
  [[ "$output" == *"Objective."* ]]

  printf '# Plan\n\nUpdated.\n' | git-lore edit test-work --file plan.md --message "lore: update plan" --repo "$REPO"
  run git-lore show test-work plan.md --repo "$REPO"
  [ "$status" -eq 0 ]
  [[ "$output" == *"Updated."* ]]

  export_dir="$(mktemp -d)"
  git-lore export test-work --dir "$export_dir" --repo "$REPO"
  [ -f "$export_dir/plan.md" ]
  grep -q "Updated" "$export_dir/plan.md"
  rm -rf "$export_dir"

  git-lore delete test-work --repo "$REPO"
  run git-lore show test-work --repo "$REPO"
  [ "$status" -ne 0 ]
}

@test "create rejects duplicate" {
  printf '# Plan\n' | git-lore create dup --repo "$REPO"
  run bash -c "printf '# Plan\n' | git-lore create dup --repo '$REPO'"
  [ "$status" -ne 0 ]
  [[ "$output" == *"already exists"* ]]
}

@test "create --associate-branch" {
  git -C "$REPO" checkout -q -b feature/test
  printf '# Plan\n' | git-lore create assoc-test --associate-branch --repo "$REPO"
  run git -C "$REPO" config --get branch.feature/test.lore
  [ "$status" -eq 0 ]
  [ "$output" = "assoc-test" ]
}

@test "flags after positional work for edit" {
  printf '# Plan\n' | git-lore create flag-test --repo "$REPO"
  printf '# Decisions\n' | git-lore edit flag-test --file decisions.md --repo "$REPO"
  run git-lore show flag-test --repo "$REPO"
  [ "$status" -eq 0 ]
  [[ "$output" == *"decisions.md"* ]]
}

#!/usr/bin/env bats
# Dispatcher and common helpers

setup() {
  EXT="$(cd "$(dirname "$BATS_TEST_FILENAME")/.." && pwd)"
  export PATH="$EXT/bin:$PATH"
  export GIT_LORE_EXT="$EXT"
}

@test "dispatcher routes help" {
  run git-lore help
  [ "$status" -eq 0 ]
  [[ "$output" == *"git lore <command>"* ]]
}

@test "dispatcher unknown command exits 2" {
  run git-lore nosuchcmd
  [ "$status" -eq 2 ]
}

@test "validate_work_id accepts valid ids" {
  source "$GIT_LORE_EXT/lib/git-lore/common.sh"
  validate_work_id "git-lore"
  validate_work_id "a"
  validate_work_id "oauth2"
}

@test "validate_work_id rejects invalid ids" {
  source "$GIT_LORE_EXT/lib/git-lore/common.sh"
  run validate_work_id "Git-Lore"
  [ "$status" -ne 0 ]
  run validate_work_id "-leading"
  [ "$status" -ne 0 ]
  run validate_work_id "has_underscore"
  [ "$status" -ne 0 ]
}

@test "validate_file_path rejects traversal" {
  source "$GIT_LORE_EXT/lib/git-lore/common.sh"
  run validate_file_path "../escape"
  [ "$status" -ne 0 ]
  run validate_file_path "/abs"
  [ "$status" -ne 0 ]
}

#!/usr/bin/env bats
# make install PREFIX handling

setup() {
  EXT="$(cd "$(dirname "$BATS_TEST_FILENAME")/.." && pwd)"
}

@test "make install defaults PREFIX to HOME/.local" {
  run make -C "$EXT" -n install
  [ "$status" -eq 0 ]
  [[ "$output" == *"install -d \"$HOME/.local/bin\""* ]]
}

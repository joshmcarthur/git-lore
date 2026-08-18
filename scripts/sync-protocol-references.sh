#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SOURCE="$ROOT/skills/protocol.md"

for skill in create-lore read-lore edit-lore sync-lore; do
  dest="$ROOT/skills/$skill/references/protocol.md"
  mkdir -p "$(dirname "$dest")"
  ln -sf ../../protocol.md "$dest"
  echo "linked $dest -> ../../protocol.md"
done

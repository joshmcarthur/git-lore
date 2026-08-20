#!/usr/bin/env bash
# Source: human-facing help for git-lore shell plugin
set -euo pipefail

cat <<'EOF'
git-lore — Git external command for Lore Works (refs/lore/*)

Usage:
  git lore <command> [flags] [args]
  git-lore <command> [flags] [args]

Commands:
  list       List local Lore Works
  show       Show a Work or file content
  create     Initialise a new Lore Work (plan.md from stdin or --from)
  edit       Update a lore file (stdin, --from, or --interactive)
  export     Export a Work tree to a directory
  delete     Delete a Lore Work ref
  branch     List or manage branch ↔ work associations
  remote     Fetch, push, sync, or compare lore refs on a remote
  serve      Serve the local Lore browser UI
  version    Print version
  help       Show this help

Common flags:
  --repo <path>   Git repository (default: current directory)

Examples:
  git lore list --repo .
  git lore show git-lore plan.md
  echo "# Plan" | git lore create my-work --associate-branch
  cat plan.md | git lore edit my-work --file plan.md --message "lore: update plan"
  git lore export my-work --dir /tmp/my-work
  git lore remote sync --remote origin
  git lore branch --set my-work
  git lore serve --open

The primary git-lore surface remains agent skills under skills/.
This CLI mirrors those workflows for humans and scripted agents.

Windows: git-ops commands require Git Bash or WSL; serve runs natively.
EOF

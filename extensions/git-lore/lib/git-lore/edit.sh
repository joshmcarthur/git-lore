#!/usr/bin/env bash
# Source: skills/edit-lore steps 5–6 — export, edit one file, commit with isolated index
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=common.sh
source "$SCRIPT_DIR/common.sh"

usage() {
  cat <<'EOF'
Usage: git lore edit <work-id> --file <path> [--from <file>] [--message <msg>] [--interactive] [--repo <path>]

Content from stdin (default), --from, or --interactive ($EDITOR).
EOF
}

repo=""
file_path=""
from=""
message=""
interactive=false
positional=()
while [[ $# -gt 0 ]]; do
  case "$1" in
    --repo)
      [[ $# -ge 2 ]] || usage_die "--repo requires a path"
      repo="$2"; shift 2 ;;
    --repo=*)
      repo="${1#--repo=}"; shift ;;
    --file)
      [[ $# -ge 2 ]] || usage_die "--file requires a path"
      file_path="$2"; shift 2 ;;
    --file=*)
      file_path="${1#--file=}"; shift ;;
    --from)
      [[ $# -ge 2 ]] || usage_die "--from requires a path"
      from="$2"; shift 2 ;;
    --from=*)
      from="${1#--from=}"; shift ;;
    --message)
      [[ $# -ge 2 ]] || usage_die "--message requires a value"
      message="$2"; shift 2 ;;
    --message=*)
      message="${1#--message=}"; shift ;;
    --interactive)
      interactive=true; shift ;;
    -h|--help)
      usage; exit 0 ;;
    --)
      shift; positional+=("$@"); break ;;
    -*)
      usage_die "unknown flag: $1" ;;
    *)
      positional+=("$1"); shift ;;
  esac
done

resolve_repo "$repo"

if [[ ${#positional[@]} -lt 1 ]]; then
  usage_die "work-id is required"
fi
if [[ -z "$file_path" ]]; then
  usage_die "--file is required"
fi
work_id="${positional[0]}"
validate_work_id "$work_id"
validate_file_path "$file_path"

if ! work_exists "$work_id"; then
  die "work \"$work_id\" not found"
fi

ref="refs/lore/${work_id}"
parent="$(git_lore rev-parse "$ref")"
git_dir_path="$(git_dir)"

dir="$(mktemp -d "${TMPDIR:-/tmp}/git-lore-edit.XXXXXX")"
index="$(mktemp "${TMPDIR:-/tmp}/git-lore-index.XXXXXX")"
cleanup() {
  rm -rf "$dir" "$index"
}
trap cleanup EXIT

# Export current lore tree
git_lore archive "$parent" | tar -x -C "$dir"

target="$dir/$file_path"
mkdir -p "$(dirname "$target")"

if [[ "$interactive" == true ]]; then
  editor="${EDITOR:-${VISUAL:-}}"
  if [[ -z "$editor" ]]; then
    die "set EDITOR or VISUAL for interactive edit"
  fi
  if [[ ! -f "$target" ]]; then
    : >"$target"
  fi
  # shellcheck disable=SC2086
  $editor "$target"
elif [[ -n "$from" ]]; then
  cat "$from" >"$target"
else
  if [[ -t 0 ]]; then
    die "no content: pipe to stdin or use --from"
  fi
  cat >"$target"
  if [[ ! -s "$target" ]]; then
    die "no content: pipe to stdin or use --from"
  fi
fi

if [[ -z "$message" ]]; then
  message="lore: update Work ${work_id}"
fi

# Isolated index — edit-lore skill step 6
export GIT_INDEX_FILE="$index"
git --git-dir="$git_dir_path" read-tree --empty
(
  cd "$dir"
  git --git-dir="$git_dir_path" add -A .
)
tree="$(git --git-dir="$git_dir_path" write-tree)"
unset GIT_INDEX_FILE

new="$(git_lore commit-tree "$tree" -p "$parent" -m "$message")"
git_lore update-ref "$ref" "$new"

printf 'Updated %s at %s (%s)\n' "$ref" "$new" "$message"

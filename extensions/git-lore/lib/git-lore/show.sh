#!/usr/bin/env bash
# Source: skills/read-lore — show Work metadata or file content
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=common.sh
source "$SCRIPT_DIR/common.sh"

usage() {
  cat <<'EOF'
Usage: git lore show <work-id> [file] [--repo <path>]
EOF
}

repo=""
positional=()
while [[ $# -gt 0 ]]; do
  case "$1" in
    --repo)
      [[ $# -ge 2 ]] || usage_die "--repo requires a path"
      repo="$2"; shift 2 ;;
    --repo=*)
      repo="${1#--repo=}"; shift ;;
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
work_id="${positional[0]}"
file_path="${positional[1]:-}"
validate_work_id "$work_id"

if [[ -n "$file_path" ]]; then
  validate_file_path "$file_path"
  git_lore show "refs/lore/${work_id}:${file_path}"
  exit 0
fi

ref="refs/lore/${work_id}"
if ! work_exists "$work_id"; then
  die "work \"$work_id\" not found"
fi

commit="$(git_lore rev-parse "$ref")"
subject="$(git_lore log -1 --format='%s' "$ref")"
printf 'Work: %s\n' "$work_id"
printf 'Ref:  %s\n' "$ref"
printf 'Head: %s\n' "$commit"
printf 'Subject: %s\n' "$subject"
printf 'Files:\n'
while IFS= read -r f; do
  [[ -n "$f" ]] && printf '  %s\n' "$f"
done < <(git_lore ls-tree -r --name-only "$ref")

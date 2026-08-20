#!/usr/bin/env bash
# Source: protocol — deleting the Lore ref removes the Work
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=common.sh
source "$SCRIPT_DIR/common.sh"

usage() {
  cat <<'EOF'
Usage: git lore delete <work-id> [--repo <path>]
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
validate_work_id "$work_id"

if ! work_exists "$work_id"; then
  die "work \"$work_id\" not found"
fi

git_lore update-ref -d "refs/lore/${work_id}"
printf 'Deleted refs/lore/%s\n' "$work_id"

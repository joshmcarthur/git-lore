#!/usr/bin/env bash
# Source: skills/edit-lore export step — git archive to directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=common.sh
source "$SCRIPT_DIR/common.sh"

usage() {
  cat <<'EOF'
Usage: git lore export <work-id> --dir <path> [--repo <path>]
EOF
}

repo=""
dir=""
positional=()
while [[ $# -gt 0 ]]; do
  case "$1" in
    --repo)
      [[ $# -ge 2 ]] || usage_die "--repo requires a path"
      repo="$2"; shift 2 ;;
    --repo=*)
      repo="${1#--repo=}"; shift ;;
    --dir)
      [[ $# -ge 2 ]] || usage_die "--dir requires a path"
      dir="$2"; shift 2 ;;
    --dir=*)
      dir="${1#--dir=}"; shift ;;
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
if [[ -z "$dir" ]]; then
  usage_die "--dir is required"
fi
work_id="${positional[0]}"
validate_work_id "$work_id"

if ! work_exists "$work_id"; then
  die "work \"$work_id\" not found"
fi

commit="$(git_lore rev-parse "refs/lore/${work_id}")"
mkdir -p "$dir"
git_lore archive "$commit" | tar -x -C "$dir"
printf 'Exported %s to %s\n' "$work_id" "$dir"

#!/usr/bin/env bash
# Source: skills/read-lore — list local refs/lore/*
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=common.sh
source "$SCRIPT_DIR/common.sh"

usage() {
  cat <<'EOF'
Usage: git lore list [--repo <path>]
EOF
}

repo=""
while [[ $# -gt 0 ]]; do
  case "$1" in
    --repo)
      [[ $# -ge 2 ]] || usage_die "--repo requires a path"
      repo="$2"; shift 2 ;;
    --repo=*)
      repo="${1#--repo=}"; shift ;;
    -h|--help)
      usage; exit 0 ;;
    -*)
      usage_die "unknown flag: $1" ;;
    *)
      usage_die "unexpected argument: $1" ;;
  esac
done

resolve_repo "$repo"

rows=()
while IFS= read -r short; do
  [[ -z "$short" ]] && continue
  id="${short#lore/}"
  if [[ "$id" == "$short" ]]; then
    continue
  fi
  if [[ "$id" == *-remote ]]; then
    continue
  fi
  commit="$(git_lore rev-parse "refs/lore/${id}")"
  subject="$(git_lore log -1 --format='%s' "refs/lore/${id}")"
  rows+=("$id"$'\t'"$(short_sha "$commit")"$'\t'"$subject")
done < <(git_lore for-each-ref --format='%(refname:short)' refs/lore 2>/dev/null || true)

if [[ ${#rows[@]} -eq 0 ]]; then
  printf 'No Lore Works found.\n'
  exit 0
fi

printf 'ID\tCOMMIT\tSUBJECT\n'
printf '%s\n' "${rows[@]}"

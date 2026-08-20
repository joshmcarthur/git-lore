#!/usr/bin/env bash
# Source: skills/create-lore + read-lore — branch.*.lore associations
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=common.sh
source "$SCRIPT_DIR/common.sh"

usage() {
  cat <<'EOF'
Usage:
  git lore branch [--repo <path>]
  git lore branch --set <work-id> [--branch <name>] [--repo <path>]
  git lore branch --unset [--branch <name>] [--repo <path>]
EOF
}

repo=""
set_work=""
unset_flag=false
branch_name=""
while [[ $# -gt 0 ]]; do
  case "$1" in
    --repo)
      [[ $# -ge 2 ]] || usage_die "--repo requires a path"
      repo="$2"; shift 2 ;;
    --repo=*)
      repo="${1#--repo=}"; shift ;;
    --set)
      [[ $# -ge 2 ]] || usage_die "--set requires a work-id"
      set_work="$2"; shift 2 ;;
    --set=*)
      set_work="${1#--set=}"; shift ;;
    --unset)
      unset_flag=true; shift ;;
    --branch)
      [[ $# -ge 2 ]] || usage_die "--branch requires a name"
      branch_name="$2"; shift 2 ;;
    --branch=*)
      branch_name="${1#--branch=}"; shift ;;
    -h|--help)
      usage; exit 0 ;;
    --)
      shift; break ;;
    -*)
      usage_die "unknown flag: $1" ;;
    *)
      usage_die "unexpected argument: $1" ;;
  esac
done

resolve_repo "$repo"

if [[ -z "$branch_name" ]]; then
  branch_name="$(git_lore branch --show-current)"
fi

if [[ -n "$set_work" ]]; then
  if [[ -z "$branch_name" ]]; then
    die "cannot determine branch: detached HEAD; use --branch"
  fi
  validate_work_id "$set_work"
  if ! work_exists "$set_work"; then
    die "work \"$set_work\" not found"
  fi
  git_lore config "branch.${branch_name}.lore" "$set_work"
  printf 'Associated branch %q with work %q\n' "$branch_name" "$set_work"
  exit 0
fi

if [[ "$unset_flag" == true ]]; then
  if [[ -z "$branch_name" ]]; then
    die "cannot determine branch: detached HEAD; use --branch"
  fi
  if ! git_lore config --unset "branch.${branch_name}.lore" 2>/dev/null; then
    # exit 5 = key not found — treat as success
    true
  fi
  printf 'Removed lore association from branch %q\n' "$branch_name"
  exit 0
fi

# List associations
out="$(git_lore config --get-regexp '^branch\..*\.lore$' 2>/dev/null || true)"
if [[ -z "$out" ]]; then
  printf 'No branch lore associations.\n'
  exit 0
fi

printf 'BRANCH\tWORK\n'
while IFS= read -r line; do
  [[ -z "$line" ]] && continue
  key="${line%% *}"
  work="${line#* }"
  # key = branch.<name>.lore
  b="${key#branch.}"
  b="${b%.lore}"
  printf '%s\t%s\n' "$b" "$work"
done <<<"$out"

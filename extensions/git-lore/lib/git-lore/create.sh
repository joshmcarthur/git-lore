#!/usr/bin/env bash
# Source: skills/create-lore steps 4–5 — orphan lore commit with plan.md
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=common.sh
source "$SCRIPT_DIR/common.sh"

usage() {
  cat <<'EOF'
Usage: git lore create <work-id> [--from <file>] [--message <msg>] [--associate-branch] [--repo <path>]

Reads plan.md content from stdin (or --from). Creates an orphan commit at refs/lore/<work-id>.
EOF
}

repo=""
from=""
message=""
associate=false
positional=()
while [[ $# -gt 0 ]]; do
  case "$1" in
    --repo)
      [[ $# -ge 2 ]] || usage_die "--repo requires a path"
      repo="$2"; shift 2 ;;
    --repo=*)
      repo="${1#--repo=}"; shift ;;
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
    --associate-branch)
      associate=true; shift ;;
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

if work_exists "$work_id"; then
  die "work \"$work_id\" already exists"
fi

# Read content preserving trailing newlines (bash $(cat) strips them)
content_file="$(mktemp "${TMPDIR:-/tmp}/git-lore-plan.XXXXXX")"
trap 'rm -f "$content_file"' EXIT
if [[ -n "$from" ]]; then
  cat "$from" >"$content_file"
else
  if [[ -t 0 ]]; then
    die "no content: pipe to stdin or use --from"
  fi
  cat >"$content_file"
  if [[ ! -s "$content_file" ]]; then
    die "no content: pipe to stdin or use --from"
  fi
fi

if [[ -z "$message" ]]; then
  message="lore: initialise Work ${work_id}"
fi

# create-lore skill: hash-object → mktree → commit-tree (orphan) → update-ref
plan_blob="$(git_lore hash-object -w --stdin <"$content_file")"
tree="$(printf '100644 blob %s\tplan.md\n' "$plan_blob" | git_lore mktree)"
commit="$(git_lore commit-tree "$tree" -m "$message")"
git_lore update-ref "refs/lore/${work_id}" "$commit"

if [[ "$associate" == true ]]; then
  branch="$(git_lore branch --show-current)"
  if [[ -z "$branch" ]]; then
    die "cannot associate branch: detached HEAD"
  fi
  git_lore config "branch.${branch}.lore" "$work_id"
  printf 'Created refs/lore/%s at %s (%s)\n' "$work_id" "$commit" "$message"
  printf 'Associated branch %q with work %q\n' "$branch" "$work_id"
else
  printf 'Created refs/lore/%s at %s (%s)\n' "$work_id" "$commit" "$message"
fi

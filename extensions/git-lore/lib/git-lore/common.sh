# Shared helpers for git-lore shell commands.
# Source: skills/protocol.md work-id rules; validate.go equivalent.

set -euo pipefail

GIT_LORE_REPO=""

die() {
  printf 'git-lore: %s\n' "$*" >&2
  exit 1
}

usage_die() {
  printf 'git-lore: %s\n\n' "$*" >&2
  if declare -F usage >/dev/null 2>&1; then
    usage >&2
  fi
  exit 2
}

# validate_work_id: lowercase letters, digits, hyphens (no leading/trailing hyphen)
validate_work_id() {
  local id="$1"
  if [[ -z "$id" ]]; then
    die "work id is empty"
  fi
  if [[ ! "$id" =~ ^[a-z0-9]+(-[a-z0-9]+)*$ ]]; then
    die "invalid work id \"$id\": must be lowercase letters, digits, and hyphens"
  fi
}

# validate_file_path: relative, no .., forward slashes only
validate_file_path() {
  local p="$1"
  if [[ -z "$p" ]]; then
    die "file path is empty"
  fi
  if [[ "$p" == *\\* ]]; then
    die "file path must use forward slashes: \"$p\""
  fi
  if [[ "$p" == /* ]]; then
    die "file path must be relative: \"$p\""
  fi
  if [[ "$p" == *..* ]]; then
    # Reject any path component that is ..
    local part
    IFS=/ read -ra parts <<<"$p"
    for part in "${parts[@]}"; do
      if [[ "$part" == ".." ]]; then
        die "file path must not contain ..: \"$p\""
      fi
    done
  fi
}

validate_remote_name() {
  local name="$1"
  if [[ -z "$name" ]]; then
    die "remote name is empty"
  fi
  if [[ "$name" == *" "* || "$name" == *$'\t'* || "$name" == *"/"* ]]; then
    die "invalid remote name: \"$name\""
  fi
}

# resolve_repo sets GIT_LORE_REPO from --repo or current directory
resolve_repo() {
  local dir="${1:-}"
  if [[ -n "$dir" ]]; then
    GIT_LORE_REPO="$(cd "$dir" && git rev-parse --show-toplevel)" || die "not a git repository: $dir"
  else
    GIT_LORE_REPO="$(git rev-parse --show-toplevel 2>/dev/null)" || die "not a git repository"
  fi
}

# Parse shared flags (--repo) from "$@"; sets GIT_LORE_REPO and leaves remaining in GIT_LORE_ARGS array.
# Call as: parse_common_flags "$@"; set -- "${GIT_LORE_ARGS[@]}"
parse_common_flags() {
  local repo=""
  GIT_LORE_ARGS=()
  while [[ $# -gt 0 ]]; do
    case "$1" in
      --repo)
        [[ $# -ge 2 ]] || usage_die "--repo requires a path"
        repo="$2"
        shift 2
        ;;
      --repo=*)
        repo="${1#--repo=}"
        shift
        ;;
      --)
        shift
        GIT_LORE_ARGS+=("$@")
        break
        ;;
      -*)
        # Leave unknown flags for the command to handle; collect as positional for now
        # Commands that need extra flags parse before/after calling this.
        GIT_LORE_ARGS+=("$1")
        shift
        ;;
      *)
        GIT_LORE_ARGS+=("$1")
        shift
        ;;
    esac
  done
  resolve_repo "$repo"
}

# git_lore: run git -C against GIT_LORE_REPO
git_lore() {
  git -C "$GIT_LORE_REPO" "$@"
}

work_exists() {
  local id="$1"
  git_lore rev-parse --verify "refs/lore/${id}" >/dev/null 2>&1
}

git_dir() {
  local d
  d="$(git_lore rev-parse --git-dir)"
  if [[ "$d" != /* ]]; then
    d="$GIT_LORE_REPO/$d"
  fi
  printf '%s\n' "$d"
}

# short_sha: first 8 chars or "-" if empty
short_sha() {
  local s="$1"
  if [[ -z "$s" ]]; then
    printf '%s\n' '-'
    return
  fi
  printf '%s\n' "${s:0:8}"
}

#!/usr/bin/env bash
# Source: skills/sync-lore — fetch, push, sync, status for refs/lore/*
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=common.sh
source "$SCRIPT_DIR/common.sh"

usage() {
  cat <<'EOF'
Usage:
  git lore remote fetch  [--remote <name>] [--work-id <id>] [--repo <path>]
  git lore remote push   [--remote <name>] [--work-id <id>] [--repo <path>]
  git lore remote sync   [--remote <name>] [--work-id <id>] [--repo <path>]
  git lore remote status [--remote <name>] [--work-id <id>] [--repo <path>]
EOF
}

ensure_refspec() {
  local remote="$1"
  local configured=false
  if git_lore config --get-all "remote.${remote}.fetch" 2>/dev/null | grep -Fq 'refs/lore/'; then
    configured=true
  fi
  if [[ "$configured" == false ]]; then
    git_lore config --add "remote.${remote}.fetch" 'refs/lore/*:refs/lore/*'
    printf 'true'
  else
    printf 'false'
  fi
}

classify_sync() {
  local local_sha="$1" remote_sha="$2" merge_base="$3"
  if [[ -z "$local_sha" && -z "$remote_sha" ]]; then
    printf 'synced'
  elif [[ -z "$local_sha" && -n "$remote_sha" ]]; then
    printf 'remote-only'
  elif [[ -n "$local_sha" && -z "$remote_sha" ]]; then
    printf 'local-only'
  elif [[ "$local_sha" == "$remote_sha" ]]; then
    printf 'synced'
  elif [[ "$merge_base" == "$remote_sha" ]]; then
    printf 'ahead'
  elif [[ "$merge_base" == "$local_sha" ]]; then
    printf 'behind'
  else
    printf 'diverged'
  fi
}

do_fetch() {
  local remote="$1" work_id="$2"
  local added
  added="$(ensure_refspec "$remote")"
  local refspec
  if [[ -n "$work_id" ]]; then
    refspec="refs/lore/${work_id}:refs/lore/${work_id}"
  else
    refspec='refs/lore/*:refs/lore/*'
  fi
  git_lore fetch "$remote" "$refspec"
  printf 'Fetched lore from %s' "$remote"
  if [[ -n "$work_id" ]]; then
    printf ' (work %s)' "$work_id"
  fi
  if [[ "$added" == true ]]; then
    printf ' [refspec added]'
  fi
  printf '\n'
}

do_push() {
  local remote="$1" work_id="$2"
  local refspec
  if [[ -n "$work_id" ]]; then
    refspec="refs/lore/${work_id}"
  else
    refspec='refs/lore/*'
  fi
  git_lore push "$remote" "$refspec"
  printf 'Pushed lore to %s' "$remote"
  if [[ -n "$work_id" ]]; then
    printf ' (work %s)' "$work_id"
  fi
  printf '\n'
}

do_status() {
  local remote="$1" filter_id="$2"
  local configured=false
  if git_lore config --get-all "remote.${remote}.fetch" 2>/dev/null | grep -Fq 'refs/lore/'; then
    configured=true
  fi
  printf 'Remote: %s (refspec configured: %s)\n' "$remote" "$configured"

  # Local refs
  declare -A local_shas=()
  local line short id sha
  while IFS= read -r line; do
    [[ -z "$line" ]] && continue
    short="${line%% *}"
    sha="${line##* }"
    id="${short#lore/}"
    [[ "$id" == "$short" ]] && continue
    [[ "$id" == *-remote ]] && continue
    local_shas["$id"]="$sha"
  done < <(git_lore for-each-ref --format='%(refname:short) %(objectname)' refs/lore 2>/dev/null || true)

  # Remote refs
  declare -A remote_shas=()
  while IFS= read -r line; do
    [[ -z "$line" ]] && continue
    sha="${line%%$'\t'*}"
    sha="${sha%% *}"
    ref="${line##*$'\t'}"
    ref="${ref##* }"
    id="${ref#refs/lore/}"
    [[ "$id" == "$ref" ]] && continue
    [[ "$id" == *-remote ]] && continue
    remote_shas["$id"]="$sha"
  done < <(git_lore ls-remote "$remote" 'refs/lore/*' 2>/dev/null || true)

  declare -A all_ids=()
  for id in "${!local_shas[@]}"; do all_ids["$id"]=1; done
  for id in "${!remote_shas[@]}"; do all_ids["$id"]=1; done

  if [[ ${#all_ids[@]} -eq 0 ]]; then
    printf 'No lore refs to compare.\n'
    return
  fi

  printf 'ID\tSTATE\tLOCAL\tREMOTE\n'
  # Stable sort by id
  local sorted
  sorted="$(printf '%s\n' "${!all_ids[@]}" | sort)"
  while IFS= read -r id; do
    [[ -z "$id" ]] && continue
    if [[ -n "$filter_id" && "$id" != "$filter_id" ]]; then
      continue
    fi
    local_sha="${local_shas[$id]:-}"
    remote_sha="${remote_shas[$id]:-}"
    merge_base=""
    if [[ -n "$local_sha" && -n "$remote_sha" && "$local_sha" != "$remote_sha" ]]; then
      merge_base="$(git_lore merge-base "$local_sha" "$remote_sha" 2>/dev/null || true)"
    fi
    state="$(classify_sync "$local_sha" "$remote_sha" "$merge_base")"
    printf '%s\t%s\t%s\t%s\n' "$id" "$state" "$(short_sha "$local_sha")" "$(short_sha "$remote_sha")"
  done <<<"$sorted"
}

if [[ $# -lt 1 ]]; then
  usage_die "remote subcommand required: fetch, push, sync, or status"
fi

sub="$1"
shift

repo=""
remote="origin"
work_id=""
while [[ $# -gt 0 ]]; do
  case "$1" in
    --repo)
      [[ $# -ge 2 ]] || usage_die "--repo requires a path"
      repo="$2"; shift 2 ;;
    --repo=*)
      repo="${1#--repo=}"; shift ;;
    --remote)
      [[ $# -ge 2 ]] || usage_die "--remote requires a name"
      remote="$2"; shift 2 ;;
    --remote=*)
      remote="${1#--remote=}"; shift ;;
    --work-id)
      [[ $# -ge 2 ]] || usage_die "--work-id requires an id"
      work_id="$2"; shift 2 ;;
    --work-id=*)
      work_id="${1#--work-id=}"; shift ;;
    -h|--help)
      usage; exit 0 ;;
    -*)
      usage_die "unknown flag: $1" ;;
    *)
      usage_die "unexpected argument: $1" ;;
  esac
done

resolve_repo "$repo"
validate_remote_name "$remote"
if [[ -n "$work_id" ]]; then
  validate_work_id "$work_id"
fi

# Guard: remote must exist
if ! git_lore remote get-url "$remote" >/dev/null 2>&1; then
  die "no such remote: $remote"
fi

case "$sub" in
  fetch)
    do_fetch "$remote" "$work_id"
    ;;
  push)
    do_push "$remote" "$work_id"
    ;;
  sync)
    do_fetch "$remote" "$work_id"
    do_push "$remote" "$work_id"
    printf 'Synced lore with %s' "$remote"
    if [[ -n "$work_id" ]]; then
      printf ' (work %s)' "$work_id"
    fi
    printf '\n'
    ;;
  status)
    do_status "$remote" "$work_id"
    ;;
  *)
    usage_die "unknown remote subcommand \"$sub\""
    ;;
esac

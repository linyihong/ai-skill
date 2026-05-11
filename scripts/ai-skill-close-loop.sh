#!/usr/bin/env bash
# Conservative Ai-skill close-loop helper.
#
# Default mode is dry-run: inspect status, detect active lock, and print commit groups.
# Use --commit to create grouped commits. Use --push with --commit to push after commits.
# Cursor bundle sync is opt-in via AI_SKILL_SYNC_CURSOR_BUNDLE=1.
set -euo pipefail

COMMIT=0
PUSH=0
LOCK_TTL_SECONDS="${AI_SKILL_LOCK_TTL_SECONDS:-1800}"
GROUP_DIR=""

usage() {
  cat <<'EOF'
Usage: scripts/ai-skill-close-loop.sh [--commit] [--push]

Default:
  Dry-run only. Prints grouped dirty files and refuses to modify when another
  active close-loop lock exists.

Options:
  --commit  Commit each recognized group separately.
  --push    Push after successful commits. Requires --commit.
  --help    Show this help.
EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --commit) COMMIT=1 ;;
    --push) PUSH=1 ;;
    --help) usage; exit 0 ;;
    *) echo "Unknown argument: $1" >&2; usage >&2; exit 2 ;;
  esac
  shift
done

if [[ "${PUSH}" -eq 1 && "${COMMIT}" -ne 1 ]]; then
  echo "ERROR: --push requires --commit." >&2
  exit 2
fi

REPO_ROOT="$(git rev-parse --show-toplevel)"
cd "${REPO_ROOT}"
LOCK_DIR="${REPO_ROOT}/.git/ai-skill-agent.lock"

die() {
  echo "ERROR: $*" >&2
  exit 1
}

is_pid_alive() {
  local pid="$1"
  [[ -n "${pid}" ]] && kill -0 "${pid}" 2>/dev/null
}

lock_age_seconds() {
  local file="$1"
  local now mtime
  now="$(date +%s)"
  if stat -f %m "${file}" >/dev/null 2>&1; then
    mtime="$(stat -f %m "${file}")"
  else
    mtime="$(stat -c %Y "${file}")"
  fi
  echo $((now - mtime))
}

check_lock() {
  if [[ ! -d "${LOCK_DIR}" ]]; then
    return 0
  fi

  local owner_pid="" owner=""
  [[ -f "${LOCK_DIR}/pid" ]] && owner_pid="$(<"${LOCK_DIR}/pid")"
  [[ -f "${LOCK_DIR}/owner" ]] && owner="$(<"${LOCK_DIR}/owner")"

  local age
  age="$(lock_age_seconds "${LOCK_DIR}")"
  if is_pid_alive "${owner_pid}" && [[ "${age}" -lt "${LOCK_TTL_SECONDS}" ]]; then
    echo "Active Ai-skill close-loop lock detected."
    echo "owner=${owner:-unknown}"
    echo "pid=${owner_pid}"
    echo "ageSeconds=${age}"
    echo "Action: not committing or pushing. Ask the active agent/user to finish, or clear stale lock only after verifying it is safe."
    exit 3
  fi

  echo "Removing stale Ai-skill close-loop lock (pid=${owner_pid:-unknown}, ageSeconds=${age})."
  rm -rf "${LOCK_DIR}"
}

acquire_lock() {
  check_lock
  mkdir "${LOCK_DIR}" || die "Cannot acquire lock: ${LOCK_DIR}"
  printf '%s\n' "$$" > "${LOCK_DIR}/pid"
  printf '%s\n' "${USER:-unknown}@$(hostname 2>/dev/null || echo unknown)" > "${LOCK_DIR}/owner"
  date -u +"%Y-%m-%dT%H:%M:%SZ" > "${LOCK_DIR}/startedAt"
  trap 'rm -rf "${LOCK_DIR}"' EXIT
}

ensure_no_git_operation_in_progress() {
  local git_dir
  git_dir="$(git rev-parse --git-dir)"
  [[ ! -e "${git_dir}/MERGE_HEAD" ]] || die "merge in progress"
  [[ ! -d "${git_dir}/rebase-merge" ]] || die "rebase in progress"
  [[ ! -d "${git_dir}/rebase-apply" ]] || die "rebase/apply in progress"
  [[ ! -e "${git_dir}/CHERRY_PICK_HEAD" ]] || die "cherry-pick in progress"
}

group_for_path() {
  local path="$1"
  case "${path}" in
    .cursor/rules/*|ai-tools/*) echo "tooling" ;;
    shared-rules/*|README.md) echo "shared" ;;
    scripts/*) echo "scripts" ;;
    skills/apk-analysis/*) echo "apk-analysis" ;;
    skills/app-development-guidance/*) echo "app-development-guidance" ;;
    skills/*)
      local rest skill
      rest="${path#skills/}"
      skill="${rest%%/*}"
      echo "skill-${skill}"
      ;;
    *) echo "unknown" ;;
  esac
}

commit_message_for_group() {
  case "$1" in
    shared) echo "docs(shared): close skill update loop" ;;
    scripts) echo "chore(scripts): update skill close-loop automation" ;;
    tooling) echo "docs(tools): update skill tool integration guidance" ;;
    apk-analysis) echo "docs(apk): close skill guidance updates" ;;
    app-development-guidance) echo "docs(app): close guidance updates" ;;
    skill-*) echo "docs(${1#skill-}): close skill updates" ;;
    *) echo "docs(ai): close skill updates" ;;
  esac
}

changed_paths() {
  git status --porcelain=v1 | while IFS= read -r line; do
    local path
    path="${line:3}"
    if [[ "${path}" == *" -> "* ]]; then
      path="${path##* -> }"
    fi
    printf '%s\n' "${path}"
  done
}

scan_diff_for_private_paths() {
  local diff="" file
  diff+="$(git diff -- "$@" || true)"
  for file in "$@"; do
    if ! git ls-files --error-unmatch "${file}" >/dev/null 2>&1 && [[ -f "${file}" ]]; then
      diff+=$'\n'
      diff+="$(git diff --no-index /dev/null "${file}" || true)"
    fi
  done
  if grep -E '/Users/[A-Za-z0-9._-]+/(Documents|Downloads|Desktop)/|Authorization:[[:space:]]*Bearer[[:space:]]+[A-Za-z0-9._~+/-]{16,}|x-api-key:[[:space:]]*[A-Za-z0-9._~+/-]{16,}' <<<"${diff}" >/dev/null; then
    echo "Sensitive-looking content detected in diff for: $*" >&2
    return 1
  fi
}

main() {
  acquire_lock
  ensure_no_git_operation_in_progress

  paths=()
  while IFS= read -r path; do
    [[ -n "${path}" ]] && paths+=("${path}")
  done < <(changed_paths)
  if [[ "${#paths[@]}" -eq 0 ]]; then
    echo "Ai-skill close-loop: working tree clean."
    return 0
  fi

  GROUP_DIR="$(mktemp -d "${TMPDIR:-/tmp}/ai-skill-close-loop.XXXXXX")"
  trap 'rm -rf "${LOCK_DIR}" "${GROUP_DIR}"' EXIT

  groups=()
  local path group group_file seen
  for path in "${paths[@]}"; do
    group="$(group_for_path "${path}")"
    if [[ "${group}" == "unknown" ]]; then
      echo "Unrecognized dirty path: ${path}" >&2
      echo "Action: not committing automatically. Add a grouping rule or handle manually." >&2
      exit 4
    fi
    group_file="${GROUP_DIR}/${group}"
    if [[ ! -f "${group_file}" ]]; then
      groups+=("${group}")
      : > "${group_file}"
    fi
    printf '%s\n' "${path}" >> "${group_file}"
  done

  echo "Ai-skill close-loop groups:"
  for group in "${groups[@]}"; do
    echo "## ${group}"
    sed '/^$/d;s/^/  - /' "${GROUP_DIR}/${group}"
  done

  if [[ "${COMMIT}" -ne 1 ]]; then
    echo "Dry-run only. Re-run with --commit to create grouped commits, or --commit --push to push."
    return 0
  fi

  for group in "${groups[@]}"; do
    group_file_array=()
    while IFS= read -r path; do
      [[ -n "${path}" ]] && group_file_array+=("${path}")
    done < "${GROUP_DIR}/${group}"
    scan_diff_for_private_paths "${group_file_array[@]}"
    git add -- "${group_file_array[@]}"
    git commit -m "$(commit_message_for_group "${group}")"
  done

  if [[ "${AI_SKILL_SYNC_CURSOR_BUNDLE:-0}" == "1" ]]; then
    ./scripts/sync-cursor-bundle.sh
  else
    echo "Skipping Cursor bundle sync (reference-only default). Set AI_SKILL_SYNC_CURSOR_BUNDLE=1 to sync local Cursor bundles."
  fi

  if [[ "${PUSH}" -eq 1 ]]; then
    git push origin "$(git branch --show-current)"
  fi

  git status --short --branch
}

main "$@"

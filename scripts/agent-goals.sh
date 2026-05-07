#!/usr/bin/env bash
# Tool-neutral project-local conversation goal ledger helper.
#
# Writes only under <PROJECT_ROOT>/.agent-goals/ by default. The ledger is
# temporary project state and should not be committed.
set -euo pipefail

PROJECT_ROOT="${PWD}"
GOAL_DIR_NAME=".agent-goals"
LOCK_TTL_SECONDS="${AGENT_GOALS_LOCK_TTL_SECONDS:-1800}"
COMMAND=""
AGENT_GOAL_ACTIVE_LOCK=""

usage() {
  cat <<'EOF'
Usage: scripts/agent-goals.sh [--project DIR] <command> [options]

Commands:
  init
      Create .agent-goals/ structure and exclude it from git via .git/info/exclude.

  status
      List active goal files and locks.

  start --id ID --title TITLE --source TEXT [--priority P1] [--next TEXT] [--criteria TEXT]
      Create or replace a goal file.

  update --id ID [--status STATUS] [--next TEXT] [--note TEXT] [--owner TEXT]
      Update status/next action and append progress.

  split --parent ID --id ID --title TITLE [--priority P2] [--next TEXT] [--criteria TEXT]
      Create a child goal linked to a parent.

  pause --id ID [--reason TEXT] [--superseded]
      Mark a goal paused or superseded.

  complete --id ID --validated [--note TEXT]
      Delete a goal only after validation is explicitly confirmed.
      Without --validated, marks the goal needs-validation and keeps it.

  cleanup
      Remove stale lock directories older than AGENT_GOALS_LOCK_TTL_SECONDS.

Options:
  --project DIR   Project root containing .agent-goals/. Defaults to cwd.
  --help          Show this help.
EOF
}

die() {
  echo "ERROR: $*" >&2
  exit 1
}

now_utc() {
  date -u +"%Y-%m-%dT%H:%M:%SZ"
}

slug_safe() {
  local value="$1"
  [[ "${value}" =~ ^[A-Za-z0-9._-]+$ ]] || die "Invalid id '${value}'. Use letters, numbers, dot, underscore, or dash."
}

ledger_root() {
  printf '%s/%s\n' "${PROJECT_ROOT}" "${GOAL_DIR_NAME}"
}

goals_dir() {
  printf '%s/goals\n' "$(ledger_root)"
}

locks_dir() {
  printf '%s/locks\n' "$(ledger_root)"
}

goal_path() {
  local id="$1"
  slug_safe "${id}"
  printf '%s/%s.md\n' "$(goals_dir)" "${id}"
}

ensure_dirs() {
  mkdir -p "$(goals_dir)" "$(locks_dir)"
  if [[ ! -f "$(ledger_root)/README.md" ]]; then
    cat > "$(ledger_root)/README.md" <<'EOF'
# Agent Goals

Temporary conversation goal ledger for agents. Do not commit this directory.
Delete goal files only after completion criteria and validation are satisfied.
EOF
  fi
}

ensure_git_exclude() {
  if git -C "${PROJECT_ROOT}" rev-parse --git-dir >/dev/null 2>&1; then
    local git_dir exclude_file
    git_dir="$(git -C "${PROJECT_ROOT}" rev-parse --git-dir)"
    exclude_file="${PROJECT_ROOT}/${git_dir}/info/exclude"
    mkdir -p "$(dirname "${exclude_file}")"
    touch "${exclude_file}"
    if ! grep -Fxq "${GOAL_DIR_NAME}/" "${exclude_file}"; then
      printf '\n%s/\n' "${GOAL_DIR_NAME}" >> "${exclude_file}"
    fi
  fi
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

is_pid_alive() {
  local pid="$1"
  [[ -n "${pid}" ]] && kill -0 "${pid}" 2>/dev/null
}

lock_path() {
  local id="$1"
  slug_safe "${id}"
  printf '%s/%s.lock\n' "$(locks_dir)" "${id}"
}

check_stale_or_active_lock() {
  local id="$1"
  local lock
  lock="$(lock_path "${id}")"
  [[ -d "${lock}" ]] || return 0

  local owner="" pid="" age
  [[ -f "${lock}/owner" ]] && owner="$(<"${lock}/owner")"
  [[ -f "${lock}/pid" ]] && pid="$(<"${lock}/pid")"
  age="$(lock_age_seconds "${lock}")"

  if is_pid_alive "${pid}" && [[ "${age}" -lt "${LOCK_TTL_SECONDS}" ]]; then
    echo "Active goal lock detected."
    echo "goal=${id}"
    echo "owner=${owner:-unknown}"
    echo "pid=${pid}"
    echo "ageSeconds=${age}"
    exit 3
  fi

  echo "Removing stale goal lock: ${id} (pid=${pid:-unknown}, ageSeconds=${age})"
  rm -rf "${lock}"
}

acquire_lock() {
  local id="$1"
  ensure_dirs
  check_stale_or_active_lock "${id}"
  local lock
  lock="$(lock_path "${id}")"
  mkdir "${lock}" || die "Cannot acquire goal lock for ${id}"
  AGENT_GOAL_ACTIVE_LOCK="${lock}"
  printf '%s\n' "$$" > "${lock}/pid"
  printf '%s\n' "${USER:-unknown}@$(hostname 2>/dev/null || echo unknown)" > "${lock}/owner"
  now_utc > "${lock}/startedAt"
  trap 'rm -rf "${AGENT_GOAL_ACTIVE_LOCK}"' EXIT
}

append_progress() {
  local file="$1"
  local note="$2"
  [[ -n "${note}" ]] || return 0
  printf -- '- %s: %s\n' "$(now_utc)" "${note}" >> "${file}"
}

set_frontmatter_value() {
  local file="$1" key="$2" value="$3"
  python3 - "$file" "$key" "$value" <<'PY'
from pathlib import Path
import sys

path = Path(sys.argv[1])
key = sys.argv[2]
value = sys.argv[3]
lines = path.read_text(encoding="utf-8").splitlines()
if not lines or lines[0] != "---":
    raise SystemExit("missing frontmatter")
end = None
for i in range(1, len(lines)):
    if lines[i] == "---":
        end = i
        break
if end is None:
    raise SystemExit("unterminated frontmatter")
prefix = f"{key}: "
for i in range(1, end):
    if lines[i].startswith(prefix):
        lines[i] = f"{key}: {value}"
        break
else:
    lines.insert(end, f"{key}: {value}")
path.write_text("\n".join(lines) + "\n", encoding="utf-8")
PY
}

cmd_init() {
  ensure_dirs
  ensure_git_exclude
  echo "Initialized $(ledger_root)"
}

cmd_status() {
  ensure_dirs
  echo "Goal ledger: $(ledger_root)"
  echo "Goals:"
  if compgen -G "$(goals_dir)/*.md" >/dev/null; then
    for file in "$(goals_dir)"/*.md; do
      printf '  - %s\n' "$(basename "${file}")"
    done
  else
    echo "  (none)"
  fi
  echo "Locks:"
  if compgen -G "$(locks_dir)/*.lock" >/dev/null; then
    for lock in "$(locks_dir)"/*.lock; do
      local owner="" pid=""
      [[ -f "${lock}/owner" ]] && owner="$(<"${lock}/owner")"
      [[ -f "${lock}/pid" ]] && pid="$(<"${lock}/pid")"
      printf '  - %s owner=%s pid=%s\n' "$(basename "${lock}")" "${owner:-unknown}" "${pid:-unknown}"
    done
  else
    echo "  (none)"
  fi
}

cmd_start() {
  local id="" title="" source="" priority="P1" next="" criteria=""
  while [[ $# -gt 0 ]]; do
    case "$1" in
      --id) id="$2"; shift 2 ;;
      --title) title="$2"; shift 2 ;;
      --source) source="$2"; shift 2 ;;
      --priority) priority="$2"; shift 2 ;;
      --next) next="$2"; shift 2 ;;
      --criteria) criteria="$2"; shift 2 ;;
      *) die "Unknown start option: $1" ;;
    esac
  done
  [[ -n "${id}" ]] || die "start requires --id"
  [[ -n "${title}" ]] || die "start requires --title"
  [[ -n "${source}" ]] || die "start requires --source"
  slug_safe "${id}"
  acquire_lock "${id}"
  local file
  file="$(goal_path "${id}")"
  cat > "${file}" <<EOF
---
id: ${id}
priority: ${priority}
status: active
owner: ${USER:-unknown}@$(hostname 2>/dev/null || echo unknown)
created: $(now_utc)
updated: $(now_utc)
project: ${PROJECT_ROOT}
---

# ${title}

## Source Request
${source}

## Scope
- In:
- Out:
- Affected paths/repos:

## Subgoals
- [ ] ${title}

## Dependencies
- none

## Progress
- $(now_utc): Goal created.

## Next Action
${next:-Define the next concrete action.}

## Completion Criteria
- [ ] ${criteria:-Define observable completion criteria.}

## Validation
- not yet validated

## Handoff Notes
- none
EOF
  echo "Started goal: ${file}"
}

cmd_update() {
  local id="" status="" next="" note="" owner=""
  while [[ $# -gt 0 ]]; do
    case "$1" in
      --id) id="$2"; shift 2 ;;
      --status) status="$2"; shift 2 ;;
      --next) next="$2"; shift 2 ;;
      --note) note="$2"; shift 2 ;;
      --owner) owner="$2"; shift 2 ;;
      *) die "Unknown update option: $1" ;;
    esac
  done
  [[ -n "${id}" ]] || die "update requires --id"
  acquire_lock "${id}"
  local file
  file="$(goal_path "${id}")"
  [[ -f "${file}" ]] || die "Goal not found: ${id}"
  [[ -z "${status}" ]] || set_frontmatter_value "${file}" "status" "${status}"
  [[ -z "${owner}" ]] || set_frontmatter_value "${file}" "owner" "${owner}"
  set_frontmatter_value "${file}" "updated" "$(now_utc)"
  if [[ -n "${next}" ]]; then
    python3 - "$file" "$next" <<'PY'
from pathlib import Path
import sys

path = Path(sys.argv[1])
next_action = sys.argv[2]
text = path.read_text(encoding="utf-8")
marker = "\n## Next Action\n"
next_marker = "\n## Completion Criteria\n"
start = text.index(marker) + len(marker)
end = text.index(next_marker, start)
text = text[:start] + next_action + "\n" + text[end:]
path.write_text(text, encoding="utf-8")
PY
  fi
  append_progress "${file}" "${note:-Goal updated.}"
  echo "Updated goal: ${file}"
}

cmd_split() {
  local parent="" id="" title="" priority="P2" next="" criteria=""
  while [[ $# -gt 0 ]]; do
    case "$1" in
      --parent) parent="$2"; shift 2 ;;
      --id) id="$2"; shift 2 ;;
      --title) title="$2"; shift 2 ;;
      --priority) priority="$2"; shift 2 ;;
      --next) next="$2"; shift 2 ;;
      --criteria) criteria="$2"; shift 2 ;;
      *) die "Unknown split option: $1" ;;
    esac
  done
  [[ -n "${parent}" ]] || die "split requires --parent"
  [[ -f "$(goal_path "${parent}")" ]] || die "Parent goal not found: ${parent}"
  cmd_start --id "${id}" --title "${title}" --source "Child goal of ${parent}" --priority "${priority}" --next "${next}" --criteria "${criteria:-Complete child goal ${title}.}"
  append_progress "$(goal_path "${parent}")" "Split child goal ${id}: ${title}"
}

cmd_pause() {
  local id="" reason="" status="paused"
  while [[ $# -gt 0 ]]; do
    case "$1" in
      --id) id="$2"; shift 2 ;;
      --reason) reason="$2"; shift 2 ;;
      --superseded) status="superseded"; shift ;;
      *) die "Unknown pause option: $1" ;;
    esac
  done
  [[ -n "${id}" ]] || die "pause requires --id"
  cmd_update --id "${id}" --status "${status}" --note "${reason:-Goal ${status}.}" --next "Resume only when this goal becomes the highest priority again."
}

cmd_complete() {
  local id="" note="" validated=0
  while [[ $# -gt 0 ]]; do
    case "$1" in
      --id) id="$2"; shift 2 ;;
      --note) note="$2"; shift 2 ;;
      --validated) validated=1; shift ;;
      *) die "Unknown complete option: $1" ;;
    esac
  done
  [[ -n "${id}" ]] || die "complete requires --id"
  if [[ "${validated}" -ne 1 ]]; then
    cmd_update --id "${id}" --status "needs-validation" --note "${note:-Completion requested without validation; keeping goal.}"
    echo "Goal kept because --validated was not provided."
    return 0
  fi
  acquire_lock "${id}"
  local file
  file="$(goal_path "${id}")"
  [[ -f "${file}" ]] || die "Goal not found: ${id}"
  rm -f "${file}"
  echo "Deleted completed goal: ${id}"
}

cmd_cleanup() {
  ensure_dirs
  local removed=0
  if compgen -G "$(locks_dir)/*.lock" >/dev/null; then
    for lock in "$(locks_dir)"/*.lock; do
      local pid="" age
      [[ -f "${lock}/pid" ]] && pid="$(<"${lock}/pid")"
      age="$(lock_age_seconds "${lock}")"
      if ! is_pid_alive "${pid}" || [[ "${age}" -ge "${LOCK_TTL_SECONDS}" ]]; then
        rm -rf "${lock}"
        removed=$((removed + 1))
      fi
    done
  fi
  echo "Removed stale locks: ${removed}"
}

parse_global() {
  while [[ $# -gt 0 ]]; do
    case "$1" in
      --project) PROJECT_ROOT="$(cd "$2" && pwd)"; shift 2 ;;
      --help) usage; exit 0 ;;
      init|status|start|update|split|pause|complete|cleanup) COMMAND="$1"; shift; break ;;
      *) die "Unknown argument: $1" ;;
    esac
  done
  [[ -n "${COMMAND}" ]] || die "Missing command. Use --help."

  case "${COMMAND}" in
    init) cmd_init "$@" ;;
    status) cmd_status "$@" ;;
    start) cmd_start "$@" ;;
    update) cmd_update "$@" ;;
    split) cmd_split "$@" ;;
    pause) cmd_pause "$@" ;;
    complete) cmd_complete "$@" ;;
    cleanup) cmd_cleanup "$@" ;;
    *) die "Unknown command: ${COMMAND}" ;;
  esac
}

parse_global "$@"

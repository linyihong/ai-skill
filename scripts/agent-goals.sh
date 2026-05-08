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
      Optional: [--plan REF] [--todo REF], repeatable.
      Optional: [--parallelization parallelizable|single-owner|non-parallelizable].

  update --id ID [--status STATUS] [--next TEXT] [--note TEXT] [--owner TEXT]
      Update status/next action and append progress.
      Optional: [--plan REF] [--todo REF], repeatable.
      Optional: [--missing TEXT] [--decision TEXT] [--strengthen TEXT].
      Optional: [--parallelization parallelizable|single-owner|non-parallelizable].

  split --parent ID --id ID --title TITLE [--priority P2] [--next TEXT] [--criteria TEXT]
      Create a child goal linked to a parent.
      Optional: [--plan REF] [--todo REF], repeatable.
      Optional: [--parallelization parallelizable|single-owner|non-parallelizable].

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

index_path() {
  printf '%s/README.md\n' "$(ledger_root)"
}

goal_path() {
  local id="$1"
  slug_safe "${id}"
  printf '%s/%s.md\n' "$(goals_dir)" "${id}"
}

refresh_index() {
  mkdir -p "$(goals_dir)" "$(locks_dir)"
  python3 - "$(ledger_root)" <<'PY'
from pathlib import Path
import sys

root = Path(sys.argv[1])
goals_dir = root / "goals"
rows = []

def frontmatter(lines):
    if not lines or lines[0] != "---":
        return {}
    end = None
    for i in range(1, len(lines)):
        if lines[i] == "---":
            end = i
            break
    if end is None:
        return {}
    data = {}
    for line in lines[1:end]:
        if ": " in line:
            key, value = line.split(": ", 1)
            data[key] = value
    return data

def section_text(text, start_marker, end_marker):
    if start_marker not in text:
        return ""
    start = text.index(start_marker) + len(start_marker)
    if end_marker in text[start:]:
        end = text.index(end_marker, start)
        return text[start:end].strip()
    return text[start:].strip()

def planning_links(text):
    section = section_text(text, "\n## Planning / Todo Links\n", "\n## Dependencies\n")
    refs = []
    for line in section.splitlines():
        if not line.startswith("|") or "---" in line or "Reference" in line:
            continue
        parts = [p.strip() for p in line.strip("|").split("|")]
        if len(parts) >= 2 and parts[1] and not parts[1].startswith("<"):
            refs.append(f"{parts[0]}: {parts[1]}")
    return "<br>".join(refs) if refs else "-"

def lock_state(goal_id):
    lock = root / "locks" / f"{goal_id}.lock"
    if not lock.exists():
        return "unlocked"
    owner = "unknown"
    pid = ""
    owner_file = lock / "owner"
    pid_file = lock / "pid"
    if owner_file.exists():
        owner = owner_file.read_text(encoding="utf-8").strip() or owner
    if pid_file.exists():
        pid = pid_file.read_text(encoding="utf-8").strip()
    if pid:
        return f"locked by {owner} pid={pid}"
    return f"locked by {owner}"

def open_work(text):
    section = section_text(text, "\n## Open Work / Decisions\n", "\n## Dependencies\n")
    items = []
    labels = {
        "- Missing work:": "missing",
        "- Decision needed:": "decision",
        "- Needs strengthening:": "strengthen",
    }
    for line in section.splitlines():
        for prefix, label in labels.items():
            if line.startswith(prefix):
                value = line[len(prefix):].strip()
                if value and not value.startswith("<"):
                    items.append(f"{label}: {value}")
    return "<br>".join(items) if items else "-"

for file in sorted(goals_dir.glob("*.md")):
    text = file.read_text(encoding="utf-8")
    lines = text.splitlines()
    meta = frontmatter(lines)
    title = file.stem
    for line in lines:
        if line.startswith("# "):
            title = line[2:].strip()
            break
    next_action = section_text(text, "\n## Next Action\n", "\n## Completion Criteria\n") or "-"
    next_action = " ".join(next_action.split())
    if len(next_action) > 100:
        next_action = next_action[:97] + "..."
    rows.append([
        meta.get("priority", ""),
        meta.get("status", ""),
        meta.get("parallelization", ""),
        meta.get("owner", ""),
        lock_state(meta.get("id", file.stem)),
        f"[{title}](goals/{file.name})",
        open_work(text),
        planning_links(text),
        next_action,
        meta.get("updated", ""),
    ])

content = [
    "# Agent Goals",
    "",
    "Temporary conversation goal ledger for humans and agents. Do not commit this directory.",
    "Use this table to see unfinished work, needed decisions, priority, and what needs strengthening.",
    "Delete goal files only after completion criteria and validation are satisfied.",
    "",
    "## Goal Table",
    "",
    "| Priority | Status | Mode | Owner | Lock | Goal | Open Work / Decisions | Planning / Todo Links | Next Action | Updated |",
    "| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |",
]
if rows:
    content.extend("| " + " | ".join(row) + " |" for row in rows)
else:
    content.append("| - | - | - | - | unlocked | No active goals | - | - | - | - |")
content.extend([
    "",
    "## Notes",
    "",
    "- Goal details live in `goals/`.",
    "- Locks live in `locks/` and are temporary.",
    "- This file is an index generated by the helper script.",
    "- Use `Open Work / Decisions` to choose what to do next after a long conversation.",
    "- If `Lock` shows another active owner for overlapping work, stop and ask before editing.",
])
(root / "README.md").write_text("\n".join(content) + "\n", encoding="utf-8")
PY
}

ensure_dirs() {
  mkdir -p "$(goals_dir)" "$(locks_dir)"
  if [[ ! -f "$(ledger_root)/README.md" ]]; then
    refresh_index
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

insert_planning_link() {
  local file="$1" type="$2" reference="$3" note="$4"
  [[ -n "${reference}" ]] || return 0
  python3 - "$file" "$type" "$reference" "$note" <<'PY'
from pathlib import Path
import sys

path = Path(sys.argv[1])
kind = sys.argv[2]
reference = sys.argv[3]
note = sys.argv[4]
text = path.read_text(encoding="utf-8")
row = f"| {kind} | {reference} | {note} |\n"
marker = "\n## Planning / Todo Links\n"
end_marker = "\n## Dependencies\n"
if marker not in text:
    insert_at = text.index(end_marker)
    section = (
        "\n## Planning / Todo Links\n"
        "| Type | Reference | Status / Note |\n"
        "| --- | --- | --- |\n"
        f"{row}"
    )
    text = text[:insert_at] + section + text[insert_at:]
else:
    start = text.index(marker) + len(marker)
    end = text.index(end_marker, start)
    section = text[start:end]
    if reference not in section:
        text = text[:end] + row + text[end:]
path.write_text(text, encoding="utf-8")
PY
}

replace_bullet_value() {
  local file="$1" section="$2" label="$3" value="$4"
  [[ -n "${value}" ]] || return 0
  python3 - "$file" "$section" "$label" "$value" <<'PY'
from pathlib import Path
import sys

path = Path(sys.argv[1])
section = sys.argv[2]
label = sys.argv[3]
value = sys.argv[4]
text = path.read_text(encoding="utf-8")
section_marker = f"\n## {section}\n"
next_marker = "\n## Dependencies\n"
if section_marker not in text:
    insert_at = text.index(next_marker)
    block = f"{section_marker}- Missing work:\n- Decision needed:\n- Needs strengthening:\n"
    text = text[:insert_at] + block + text[insert_at:]
start = text.index(section_marker) + len(section_marker)
end = text.index(next_marker, start)
lines = text[start:end].splitlines()
prefix = f"- {label}:"
for i, line in enumerate(lines):
    if line.startswith(prefix):
        lines[i] = f"{prefix} {value}"
        break
else:
    lines.append(f"{prefix} {value}")
text = text[:start] + "\n".join(lines).rstrip() + "\n" + text[end:]
path.write_text(text, encoding="utf-8")
PY
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
  refresh_index
  echo "Initialized $(ledger_root)"
}

cmd_status() {
  ensure_dirs
  refresh_index
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
  local id="" title="" source="" priority="P1" next="" criteria="" parallelization="parallelizable"
  local -a plan_links=()
  local -a todo_links=()
  while [[ $# -gt 0 ]]; do
    case "$1" in
      --id) id="$2"; shift 2 ;;
      --title) title="$2"; shift 2 ;;
      --source) source="$2"; shift 2 ;;
      --priority) priority="$2"; shift 2 ;;
      --next) next="$2"; shift 2 ;;
      --criteria) criteria="$2"; shift 2 ;;
      --parallelization) parallelization="$2"; shift 2 ;;
      --plan) plan_links+=("$2"); shift 2 ;;
      --todo) todo_links+=("$2"); shift 2 ;;
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
parallelization: ${parallelization}
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

## Planning / Todo Links
| Type | Reference | Status / Note |
| --- | --- | --- |

## Open Work / Decisions
- Missing work:
- Decision needed:
- Needs strengthening:

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
  local link
  for link in "${plan_links[@]+"${plan_links[@]}"}"; do
    insert_planning_link "${file}" "plan" "${link}" "linked at creation"
  done
  for link in "${todo_links[@]+"${todo_links[@]}"}"; do
    insert_planning_link "${file}" "todo" "${link}" "linked at creation"
  done
  refresh_index
  echo "Started goal: ${file}"
}

cmd_update() {
  local id="" status="" next="" note="" owner="" parallelization=""
  local missing="" decision="" strengthen=""
  local -a plan_links=()
  local -a todo_links=()
  while [[ $# -gt 0 ]]; do
    case "$1" in
      --id) id="$2"; shift 2 ;;
      --status) status="$2"; shift 2 ;;
      --next) next="$2"; shift 2 ;;
      --note) note="$2"; shift 2 ;;
      --owner) owner="$2"; shift 2 ;;
      --parallelization) parallelization="$2"; shift 2 ;;
      --plan) plan_links+=("$2"); shift 2 ;;
      --todo) todo_links+=("$2"); shift 2 ;;
      --missing) missing="$2"; shift 2 ;;
      --decision) decision="$2"; shift 2 ;;
      --strengthen) strengthen="$2"; shift 2 ;;
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
  [[ -z "${parallelization}" ]] || set_frontmatter_value "${file}" "parallelization" "${parallelization}"
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
  local link
  for link in "${plan_links[@]+"${plan_links[@]}"}"; do
    insert_planning_link "${file}" "plan" "${link}" "linked during update"
  done
  for link in "${todo_links[@]+"${todo_links[@]}"}"; do
    insert_planning_link "${file}" "todo" "${link}" "linked during update"
  done
  replace_bullet_value "${file}" "Open Work / Decisions" "Missing work" "${missing}"
  replace_bullet_value "${file}" "Open Work / Decisions" "Decision needed" "${decision}"
  replace_bullet_value "${file}" "Open Work / Decisions" "Needs strengthening" "${strengthen}"
  append_progress "${file}" "${note:-Goal updated.}"
  refresh_index
  echo "Updated goal: ${file}"
}

cmd_split() {
  local parent="" id="" title="" priority="P2" next="" criteria="" parallelization="parallelizable"
  local -a plan_links=()
  local -a todo_links=()
  while [[ $# -gt 0 ]]; do
    case "$1" in
      --parent) parent="$2"; shift 2 ;;
      --id) id="$2"; shift 2 ;;
      --title) title="$2"; shift 2 ;;
      --priority) priority="$2"; shift 2 ;;
      --next) next="$2"; shift 2 ;;
      --criteria) criteria="$2"; shift 2 ;;
      --parallelization) parallelization="$2"; shift 2 ;;
      --plan) plan_links+=("$2"); shift 2 ;;
      --todo) todo_links+=("$2"); shift 2 ;;
      *) die "Unknown split option: $1" ;;
    esac
  done
  [[ -n "${parent}" ]] || die "split requires --parent"
  [[ -f "$(goal_path "${parent}")" ]] || die "Parent goal not found: ${parent}"
  local -a start_args
  start_args=(--id "${id}" --title "${title}" --source "Child goal of ${parent}" --priority "${priority}" --next "${next}" --criteria "${criteria:-Complete child goal ${title}.}" --parallelization "${parallelization}")
  local link
  for link in "${plan_links[@]+"${plan_links[@]}"}"; do
    start_args+=(--plan "${link}")
  done
  for link in "${todo_links[@]+"${todo_links[@]}"}"; do
    start_args+=(--todo "${link}")
  done
  cmd_start "${start_args[@]}"
  append_progress "$(goal_path "${parent}")" "Split child goal ${id}: ${title}"
  refresh_index
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
  refresh_index
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
  refresh_index
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

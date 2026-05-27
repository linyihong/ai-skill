#!/usr/bin/env bash
# SessionStart hook: auto-load bootstrap context so Claude doesn't need to
# decide whether to run /bootstrap. The hook itself does the mechanical
# work (cat files, query SQLite) and injects everything as additionalContext
# via hookSpecificOutput. Claude only needs to echo the pre-computed
# Bootstrap Receipt in its first response.
#
# Per CLAUDE.md §「第一輪使用者互動」and runtime/core-bootstrap.yaml
# §per_session_obligations[obligation.bootstrap.receipt].
#
# Diagnostic logging: /tmp/ai-skill-sessionstart-hook.log

set -uo pipefail

LOG=/tmp/ai-skill-sessionstart-hook.log
TS=$(date +%Y-%m-%dT%H:%M:%S)

# Read hook input (JSON on stdin), but we don't strictly need fields from it
input_json="$(cat)"
echo "=== $TS SessionStart hook fired ===" >> "$LOG"
echo "input: $input_json" >> "$LOG"

PROJECT_DIR="${CLAUDE_PROJECT_DIR:-$(pwd)}"

# Query runtime.db for phase / obligation count / gate count
PHASE=$(sqlite3 "$PROJECT_DIR/runtime/runtime.db" "SELECT phase_id FROM phase_machine LIMIT 1;" 2>/dev/null || echo "unknown")
OBLIG_COUNT=$(sqlite3 "$PROJECT_DIR/runtime/runtime.db" "SELECT COUNT(*) FROM obligations;" 2>/dev/null || echo "?")
GATE_COUNT=$(sqlite3 "$PROJECT_DIR/runtime/runtime.db" "SELECT COUNT(*) FROM gates;" 2>/dev/null || echo "?")
# Per-turn obligations are defined in runtime/core-bootstrap.yaml
# §per_turn_obligations but NOT projected to the obligations table
# (that table holds phase-based execution obligations). Hardcode here
# to match the YAML source of truth.
ACTIVE_PER_TURN="obligation.cognitive.mode_report, obligation.finality.close_loop_check"

# Build the bootstrap context (3 required reads inline + receipt template)
read_safe() {
  local path="$1"
  if [[ -f "$path" ]]; then
    cat "$path"
  else
    echo "(missing: $path)"
  fi
}

CORE=$(read_safe "$PROJECT_DIR/CORE_BOOTSTRAP.md")
RULE_WEIGHT=$(read_safe "$PROJECT_DIR/enforcement/rule-weight.md")
DEPENDENCY=$(read_safe "$PROJECT_DIR/enforcement/dependency-reading.md")
GOAL_LEDGER=$(read_safe "$PROJECT_DIR/enforcement/conversation-goal-ledger.md")

CONTEXT=$(cat <<EOF
[ai-skill SessionStart] Bootstrap auto-loaded. The agent does NOT need to read these files again — they are already in context. Your first user-facing response MUST begin with this Bootstrap Receipt (verbatim), then proceed to answer the user:

Bootstrap: rules=✓ phase=${PHASE} obligations=${OBLIG_COUNT} gates=${GATE_COUNT}
Active per-turn obligations: ${ACTIVE_PER_TURN}

Final response MUST also end with a Cognitive Mode 報告 block (compact form is fine for trivial tasks). Per-turn enforcement: see runtime/core-bootstrap.yaml §per_turn_obligations.

--- CORE_BOOTSTRAP.md (companion) ---
${CORE}

--- enforcement/rule-weight.md ---
${RULE_WEIGHT}

--- enforcement/dependency-reading.md ---
${DEPENDENCY}

--- enforcement/conversation-goal-ledger.md ---
${GOAL_LEDGER}
EOF
)

# Emit hookSpecificOutput with additionalContext (JSON to stdout)
/opt/homebrew/bin/jq -n --arg ctx "$CONTEXT" '{
  hookSpecificOutput: {
    hookEventName: "SessionStart",
    additionalContext: $ctx
  }
}'

echo "phase=$PHASE obligations=$OBLIG_COUNT gates=$GATE_COUNT" >> "$LOG"
exit 0

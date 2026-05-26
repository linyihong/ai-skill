#!/usr/bin/env bash
# PostToolUse hook: after any tool execution, if Bootstrap Receipt not yet
# in transcript, inject a structured reminder into hookSpecificOutput so
# Claude sees it before deciding what to do next.
#
# This does NOT block (exit 0 always) — PostToolUse cannot reliably block
# without breaking tool result delivery. Instead it injects a visible
# reminder that Claude must emit Bootstrap Receipt before the next response.
#
# Diagnostic logging: /tmp/ai-skill-posttool-hook.log

set -uo pipefail

LOG=/tmp/ai-skill-posttool-hook.log
TS=$(date +%Y-%m-%dT%H:%M:%S)

input_json="$(cat)"
echo "=== $TS PostToolUse hook fired ===" >> "$LOG"

python3 - "$input_json" <<'PYEOF'
import json
import os
import sys
import hashlib

raw = sys.argv[1] if len(sys.argv) > 1 else ""
try:
    payload = json.loads(raw)
except Exception as e:
    print(f"ALLOW_BAD_INPUT: {e}", file=sys.stderr)
    sys.exit(0)

transcript_path = payload.get("transcript_path", "")
tool_name = payload.get("tool_name", "")

# Check cache first
if transcript_path:
    cache_key = hashlib.md5(transcript_path.encode()).hexdigest()[:12]
    cache_file = f"/tmp/ai-skill-bootstrap-{cache_key}.done"
    if os.path.exists(cache_file):
        print("CACHED_DONE", file=sys.stderr)
        sys.exit(0)

# Scan for Bootstrap Receipt
RECEIPT_MARKER = "Bootstrap:"
found = False
if transcript_path and os.path.isfile(transcript_path):
    try:
        with open(transcript_path, "r", encoding="utf-8") as fh:
            for line in fh:
                line = line.strip()
                if not line:
                    continue
                try:
                    entry = json.loads(line)
                except Exception:
                    continue
                role = entry.get("type") or entry.get("role")
                if role != "assistant":
                    continue
                chunks = []
                msg = entry.get("message")
                if isinstance(msg, dict):
                    c = msg.get("content")
                    if isinstance(c, str):
                        chunks.append(c)
                    elif isinstance(c, list):
                        for item in c:
                            if isinstance(item, dict) and isinstance(item.get("text"), str):
                                chunks.append(item["text"])
                            elif isinstance(item, str):
                                chunks.append(item)
                else:
                    c = entry.get("content")
                    if isinstance(c, str):
                        chunks.append(c)
                if RECEIPT_MARKER in "\n".join(chunks):
                    found = True
                    break
    except Exception as e:
        print(f"READ_ERROR: {e}", file=sys.stderr)
        sys.exit(0)

if found:
    if transcript_path:
        try:
            open(cache_file, "w").close()
        except Exception:
            pass
    print("RECEIPT_FOUND", file=sys.stderr)
    sys.exit(0)

# No receipt yet — inject reminder via hookSpecificOutput
reminder = (
    "[ai-skill PostToolUse] Bootstrap Receipt not yet emitted. "
    "Before writing your next response, you MUST:\n"
    "1. Read CORE_BOOTSTRAP.md\n"
    "2. Query runtime/runtime.db (phase / obligations / gates)\n"
    "3. Read enforcement/rule-weight.md, enforcement/dependency-reading.md, "
    "enforcement/conversation-goal-ledger.md\n"
    "4. Output Bootstrap Receipt as the first line of your response:\n"
    "   Bootstrap: rules=✓ phase=<phase-id> obligations=<n> gates=<n>\n"
    "   Active per-turn obligations: <obligation ids>"
)
output = {
    "hookSpecificOutput": {
        "hookEventName": "PostToolUse",
        "additionalContext": reminder
    }
}
print(json.dumps(output))
print("INJECTED_REMINDER", file=sys.stderr)
sys.exit(0)
PYEOF

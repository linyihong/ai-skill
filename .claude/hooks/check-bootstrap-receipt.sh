#!/usr/bin/env bash
# PreToolUse hook: block any non-Read tool call until Bootstrap Receipt is present
# in the transcript. Once confirmed, caches result in /tmp to avoid re-scanning.
#
# Per runtime/core-bootstrap.yaml §per_session_obligations[obligation.bootstrap.receipt]
# and gate.bootstrap.receipt_present.
#
# Diagnostic logging: /tmp/ai-skill-bootstrap-hook.log

set -uo pipefail

LOG=/tmp/ai-skill-bootstrap-hook.log
TS=$(date +%Y-%m-%dT%H:%M:%S)

input_json="$(cat)"
echo "=== $TS PreToolUse hook fired ===" >> "$LOG"

result=$(python3 - "$input_json" <<'PYEOF'
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

tool_name = payload.get("tool_name", "")
transcript_path = payload.get("transcript_path", "")

print(f"DIAG tool={tool_name!r} transcript={transcript_path!r}", file=sys.stderr)

# Read is the only always-safe passthrough before Bootstrap Receipt.
# Other tools require receipt, except the short-lived SessionStart flag path below.
ALWAYS_ALLOW = {"Read"}

if tool_name in ALWAYS_ALLOW:
    print(f"ALLOW_READ_TOOL: {tool_name}", file=sys.stderr)
    sys.exit(0)

if not transcript_path or not os.path.isfile(transcript_path):
    print(f"ALLOW_NO_TRANSCRIPT: {transcript_path!r}", file=sys.stderr)
    sys.exit(0)

# Cache flag: once receipt confirmed, skip re-scan for this transcript file
cache_key = hashlib.md5(transcript_path.encode()).hexdigest()[:12]
cache_file = f"/tmp/ai-skill-bootstrap-{cache_key}.done"

if os.path.exists(cache_file):
    print("ALLOW_CACHED", file=sys.stderr)
    sys.exit(0)

# Check SessionStart flag: if SessionStart fired < 120s ago for this project,
# treat bootstrap context as already injected and allow through.
# This prevents double-bootstrap when agent outputs text + calls tools in
# the same turn (tool fires before text is committed to transcript).
project_dir = os.environ.get("CLAUDE_PROJECT_DIR", "")
if project_dir:
    import time
    project_hash = hashlib.md5(project_dir.encode()).hexdigest()[:12]
    flag_file = f"/tmp/ai-skill-sessionstart-{project_hash}.flag"
    if os.path.exists(flag_file):
        try:
            flag_ts = int(open(flag_file).read().strip())
            if time.time() - flag_ts < 120:
                try:
                    open(cache_file, "w").close()
                except Exception:
                    pass
                print("ALLOW_SESSIONSTART_FLAG", file=sys.stderr)
                sys.exit(0)
        except Exception:
            pass

# Scan transcript for Bootstrap Receipt in any assistant message
RECEIPT_MARKER = "Bootstrap:"
found = False
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
            text = "\n".join(chunks)
            if RECEIPT_MARKER in text:
                found = True
                break
except Exception as e:
    print(f"ALLOW_READ_ERROR: {e}", file=sys.stderr)
    sys.exit(0)

if found:
    # Write cache flag so future calls skip scanning
    try:
        open(cache_file, "w").close()
    except Exception:
        pass
    print("ALLOW_RECEIPT_FOUND", file=sys.stderr)
    sys.exit(0)

# No receipt found — block this tool call
print("BLOCK_NO_RECEIPT", file=sys.stderr)
sys.exit(2)
PYEOF
)
exit_code=$?

echo "exit_code: $exit_code" >> "$LOG"

if [[ $exit_code -eq 0 ]]; then
    exit 0
fi

cat >&2 <<'REASON'
[ai-skill PreToolUse hook] Bootstrap Receipt missing.

Before calling any tool other than Read, you MUST:
1. Read CORE_BOOTSTRAP.md
2. Query runtime/runtime.db for phase / obligations / gates
3. Read the 3 required_reads: enforcement/rule-weight.md, enforcement/dependency-reading.md, enforcement/conversation-goal-ledger.md
4. Output the Bootstrap Receipt in your first user-facing response:
   Bootstrap: rules=✓ phase=<phase-id> obligations=<n> gates=<n>
   Active per-turn obligations: <obligation ids>

Only Read tool calls are allowed before the Receipt is emitted.
REASON
exit 2

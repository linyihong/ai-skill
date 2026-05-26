#!/usr/bin/env bash
# Stop hook: ensure last assistant response contains the per-turn
# `### Cognitive Mode 報告` block before Claude is allowed to stop.
#
# Per runtime/core-bootstrap.yaml §per_turn_obligations[obligation.cognitive.mode_report].
# Implements D-strategy from the archived bootstrap-contract-yaml-migration plan.
#
# Diagnostic logging: /tmp/ai-skill-stop-hook.log (append each invocation)

set -uo pipefail

LOG=/tmp/ai-skill-stop-hook.log
TS=$(date +%Y-%m-%dT%H:%M:%S)

input_json="$(cat)"
echo "=== $TS Stop hook fired ===" >> "$LOG"
echo "input_json: $input_json" >> "$LOG"

result=$(python3 - "$input_json" <<'PYEOF'
import json
import os
import sys

raw = sys.argv[1] if len(sys.argv) > 1 else ""
try:
    payload = json.loads(raw)
except Exception as e:
    print(f"ALLOW_BAD_INPUT: {e}", file=sys.stderr)
    sys.exit(0)

if payload.get("stop_hook_active") is True:
    print("ALLOW_LOOP_GUARD", file=sys.stderr)
    sys.exit(0)

transcript_path = payload.get("transcript_path", "")
if not transcript_path or not os.path.isfile(transcript_path):
    print(f"ALLOW_NO_TRANSCRIPT: path={transcript_path!r}", file=sys.stderr)
    sys.exit(0)

marker = "### Cognitive Mode 報告"
last_assistant_text = None
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
            if chunks:
                last_assistant_text = "\n".join(chunks)
except Exception as e:
    print(f"ALLOW_READ_ERROR: {e}", file=sys.stderr)
    sys.exit(0)

if last_assistant_text is None:
    print("ALLOW_NO_ASSISTANT_MSG", file=sys.stderr)
    sys.exit(0)

# Diagnostic: log last assistant text length and tail
print(f"DIAG last_msg_len={len(last_assistant_text)} tail={last_assistant_text[-200:]!r}", file=sys.stderr)

if marker in last_assistant_text:
    print("ALLOW_BLOCK_PRESENT", file=sys.stderr)
    sys.exit(0)

print("BLOCK_MISSING", file=sys.stderr)
sys.exit(2)
PYEOF
)
exit_code=$?

# Append diag to log
echo "result: $result" >> "$LOG"
echo "exit_code: $exit_code" >> "$LOG"

if [[ $exit_code -eq 0 ]]; then
  exit 0
fi

cat >&2 <<'REASON'
[ai-skill Stop hook] Missing obligation: your final response did not include the `### Cognitive Mode 報告` block.

Per runtime/core-bootstrap.yaml §per_turn_obligations[obligation.cognitive.mode_report], every final user-facing response MUST end with a Cognitive Mode block (compact 1-line for trivial all-default tasks: `Cognitive: <e>·<c>·<g>·<m> / V:<v> / Cost:<cost> / Sig:<signal>`; full 6-row markdown table otherwise).

Please append the block to your response now, then stop again. Canonical format spec: runtime/core-bootstrap.yaml. Query active obligations: `ai-skill runtime obligations`.
REASON
exit 2
